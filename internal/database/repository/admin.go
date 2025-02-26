package repository

import (
	"context"
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"library-management/backend/internal/util"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type AdminRepositoryInterface interface {
	AddBook(*context.Context, *model.BookInventory, string) error
	RemoveBook(*context.Context, string) error
	UpdateBook(*context.Context, *model.BookInventory) error
	ListIssueRequests(*context.Context) (*[]model.BookInventory, error)
	ApproveIssueRequest(*context.Context, string) error
	RejectIssueRequest(*context.Context, string) error
}

type AdminRepository struct {
	db        *gorm.DB
	txManager *transaction.TxManager
	mu        sync.RWMutex
}

func NewAdminRepository(db *gorm.DB, txManager *transaction.TxManager) *AdminRepository {
	return &AdminRepository{
		db:        db,
		txManager: txManager,
	}
}

func (admin *AdminRepository) AddBook(ctx context.Context, book *model.BookInventory, email string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("email = ?", email).First(&user)
		if result.Error != nil {
			return result.Error
		}

		if user.Role != util.AdminRole {
			return errors.New("access denied, provide a valid Admin email")
		}

		book.LibID = user.LibID
		var existingBook model.BookInventory
		result = tx.Set("gorm:query_option", "FOR UPDATE").Where("isbn = ?", book.ISBN).First(&existingBook)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}

		log.Print(*existingBook.LibID, *user.LibID)
		if result.RowsAffected > 0 {
			if *existingBook.LibID != *user.LibID {
				return errors.New("book with same ISBN already exists in another library")
			}
			return tx.Model(&model.BookInventory{}).Where("isbn = ?", existingBook.ISBN).Update("total_copies", existingBook.TotalCopies+1).Update("available_copies", existingBook.AvailableCopies+1).Error
		}

		return tx.Create(book).Error
	})
}

func (admin *AdminRepository) RemoveBook(ctx context.Context, isbn string, userID string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user)
		if result.Error != nil {
			return result.Error
		}

		var existingBook model.BookInventory
		result = tx.Set("gorm:query_option", "FOR UPDATE").Where("isbn = ?", isbn).Where("lib_id = ?", user.LibID).First(&existingBook)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("book with supplied ISBN not found in database")
			}
			return result.Error
		}

		if (existingBook.AvailableCopies == 1) && (existingBook.TotalCopies == 1) {
			return tx.Model(&model.BookInventory{}).
				Where("isbn = ?", existingBook.ISBN).
				Delete(&existingBook).Error
		}

		if existingBook.AvailableCopies > 0 {
			return tx.Model(&model.BookInventory{}).
				Where("isbn = ?", existingBook.ISBN).
				Update("total_copies", existingBook.TotalCopies-1).
				Update("available_copies", existingBook.AvailableCopies-1).Error
		}

		return errors.New("cannot remove issued books")
	})
}

func (admin *AdminRepository) UpdateBook(ctx context.Context, isbn string, title, authors, publisher, version, userID string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user)
		if result.Error != nil {
			return result.Error
		}

		var existingBook model.BookInventory
		result = tx.Set("gorm:query_option", "FOR UPDATE").Where("isbn = ?", isbn).Where("lib_id = ?", user.LibID).First(&existingBook)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("book with supplied ISBN not found in database")
			}
			return result.Error
		}

		query := `update book_inventories set title = ?, authors = ?, publisher = ?, version = ? where isbn = ?`
		return tx.Exec(query, title, authors, publisher, version, isbn).Error
	})
}

func (admin *AdminRepository) ListIssueRequests(ctx context.Context, requestDetails *[]model.IssueRequestDetails, adminID string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var admin model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", adminID).First(&admin)
		if result.Error != nil {
			return result.Error
		}
		lib_id := admin.LibID
		log.Print(lib_id)
		query := `SELECT r.*, b.title as book_title, b.available_copies FROM request_events r, book_inventories b
              WHERE r.book_id = b.isbn AND r.approver_id IS NULL AND b.lib_id = '` + *lib_id + `'`
		return tx.Set("gorm:query_option", "FOR SHARE").
			Raw(query).
			Scan(requestDetails).
			Error
	})
}

func (admin *AdminRepository) ApproveIssueRequest(ctx context.Context, requestID string, approverID string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingIssueRequest model.RequestEvents
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("req_id = ?", requestID).First(&existingIssueRequest).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("invalid Issue Request ID")
			}
			return err
		}

		var bookInventory model.BookInventory
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("isbn = ?", existingIssueRequest.BookID).First(&bookInventory).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("invalid ISBN in issue request")
			}
			return err
		}

		if bookInventory.AvailableCopies < 1 {
			return errors.New("no available copies in inventory")
		}

		if err := tx.Model(&model.BookInventory{}).Where("isbn = ?", bookInventory.ISBN).Update("available_copies", bookInventory.AvailableCopies-1).Error; err != nil {
			return err
		}

		approvalDate := time.Now()
		expectedReturnDate := approvalDate.Add(time.Hour * 24 * 7)
		if err := tx.Model(&model.RequestEvents{}).Where("req_id = ?", requestID).Update("approver_id", approverID).Update("approval_date", approvalDate.Format(time.RFC3339)).Error; err != nil {
			return err
		}

		issueRegister := model.IssueRegistry{
			IssueID:            util.RandomUUID(),
			BookID:             existingIssueRequest.BookID,
			ReaderID:           existingIssueRequest.ReaderID,
			IssueApproverID:    approverID,
			IssueStatus:        "open",
			IssueDate:          approvalDate.Format(time.RFC3339),
			ExpectedReturnDate: expectedReturnDate.Format(time.RFC3339),
			ReturnDate:         nil,
			ReturnApproverID:   nil,
		}
		return tx.Model(&model.IssueRegistry{}).Create(issueRegister).Error
	})
}

func (admin *AdminRepository) RejectIssueRequest(ctx context.Context, requestID string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingIssueRequest model.RequestEvents
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("req_id = ?", requestID).First(&existingIssueRequest)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("invalid Issue Request ID")
			}
			return result.Error
		}

		return tx.Model(&model.RequestEvents{}).Where("req_id = ?", requestID).Delete(&existingIssueRequest).Error
	})
}
