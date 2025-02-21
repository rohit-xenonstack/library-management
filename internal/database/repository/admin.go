package repository

import (
	"context"
	"database/sql"
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"library-management/backend/internal/util"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepositoryInterface interface {
	AddBook(*context.Context, *model.BookInventory, string) error
	RemoveBook(*context.Context, string) error
	UpdateBook(*context.Context, *model.BookInventory) error
	ListIssueRequest(*context.Context) (*[]model.BookInventory, error)
	ApproveIssueRequest(*context.Context, uuid.UUID) error
	DenyIssueRequest(*context.Context, uuid.UUID) error
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

type UpdateFields struct {
	Title     *string
	Authors   *string
	Publisher *string
	Version   *string
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
			return errors.New("Access Denied. Provide a valid Admin email")
		}

		book.LibID = user.LibID
		var existingBook model.BookInventory
		result = tx.Set("gorm:query_option", "FOR UPDATE").
			Where("isbn = ?", book.ISBN).
			First(&existingBook)

		if result.RowsAffected > 0 {
			return tx.Model(&model.BookInventory{}).
				Where("isbn = ?", existingBook.ISBN).
				Update("total_copies", existingBook.TotalCopies+1).
				Update("available_copies", existingBook.AvailableCopies+1).Error
		}

		return tx.Create(book).Error
	})
}

func (admin *AdminRepository) RemoveBook(ctx context.Context, isbn string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingBook model.BookInventory
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("isbn = ?", isbn).First(&existingBook)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Book with supplied ISBN not found in database")
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

		return errors.New("Cannot Remove issued books")
	})
}

func (admin *AdminRepository) UpdateBook(ctx context.Context, isbn string, title, authors, publisher, version *string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingBook model.BookInventory
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("isbn = ?", isbn).First(&existingBook)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Book with supplied ISBN not found in database")
		}

		updateFields := UpdateFields{
			Title:     title,
			Authors:   authors,
			Publisher: publisher,
			Version:   version,
		}
		if updateFields.Title == nil {
			updateFields.Title = &existingBook.Title
		}
		if updateFields.Authors == nil {
			updateFields.Authors = &existingBook.Authors
		}
		if updateFields.Publisher == nil {
			updateFields.Publisher = &existingBook.Publisher
		}
		if updateFields.Version == nil {
			updateFields.Version = &existingBook.Version
		}

		return tx.Model(&model.BookInventory{}).
			Where("isbn = ?", existingBook.ISBN).Select("title", "authors", "publishers", "version").Updates(updateFields).Error
	})
}

func (admin *AdminRepository) ListIssueRequests(ctx context.Context, requestEvents *[]model.RequestEvents) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		result := tx.Set("gorm:query_option", "FOR SHARE").Model(&model.RequestEvents{}).Where("approver_id IS NULL").Find(&requestEvents)
		return result.Error
	})
}

func (admin *AdminRepository) ApproveIssueRequest(ctx context.Context, requestID uuid.UUID, approverID uuid.UUID) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingIssueRequest model.RequestEvents
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("request_id = ?", requestID).First(&existingIssueRequest).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("Invalid Issue Request ID")
			}
			return err
		}

		approvalDate := time.Now()
		expectedReturnDate := approvalDate.Add(time.Hour * 24 * 7)
		if err := tx.Model(&model.RequestEvents{}).Where("request_id = ?", requestID).Update("approver_id", approverID).Update("approval_date", approvalDate).Error; err != nil {
			return err
		}

		issueRegister := model.IssueRegistry{
			IssueID:            util.RandomUUID(),
			BookID:             existingIssueRequest.BookID,
			ReaderID:           &existingIssueRequest.ReaderID,
			IssueApproverID:    &approverID,
			IssueStatus:        "open",
			IssueDate:          approvalDate,
			ExpectedReturnDate: expectedReturnDate,
			ReturnDate:         sql.NullTime{},
			ReturnApproverID:   nil,
		}
		return tx.Model(&model.IssueRegistry{}).Create(issueRegister).Error
	})
}

func (admin *AdminRepository) RejectIssueRequest(ctx context.Context, requestID uuid.UUID) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingIssueRequest model.RequestEvents
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("request_id = ?", requestID).First(&existingIssueRequest)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Invalid Issue Request ID")
		}

		return tx.Model(&model.RequestEvents{}).Where("request_id = ?", requestID).Delete(&existingIssueRequest).Error
	})
}
