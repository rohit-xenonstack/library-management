package repository

import (
	"database/sql"
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"library-management/backend/internal/util"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReaderRepository struct {
	db        *gorm.DB
	txManager *transaction.TxManager
	mu        sync.RWMutex
}

func NewReaderRepository(db *gorm.DB, txManager *transaction.TxManager) *ReaderRepository {
	return &ReaderRepository{
		db:        db,
		txManager: txManager,
	}
}

func (reader *ReaderRepository) SearchBookByTitle(ctx *gin.Context, title string, books *[]model.BookInventory) error {
	reader.mu.Lock()
	defer reader.mu.Unlock()

	return reader.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		query := `select * from book_inventories where title like '%` + title + `%'`
		return tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.BookInventory{}).Raw(query).Scan(&books).Error
	})
}

func (reader *ReaderRepository) SearchBookByAuthor(ctx *gin.Context, author string, books *[]model.BookInventory) error {
	reader.mu.Lock()
	defer reader.mu.Unlock()

	return reader.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		query := `select * from book_inventories where authors like '%` + author + `%'`
		return tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.BookInventory{}).Raw(query).Scan(&books).Error
	})
}

func (reader *ReaderRepository) SearchBookByPublisher(ctx *gin.Context, publisher string, books *[]model.BookInventory) error {
	reader.mu.Lock()
	defer reader.mu.Unlock()

	return reader.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		query := `select * from book_inventories where publisher like '%` + publisher + `%'`
		return tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.BookInventory{}).Raw(query).Scan(&books).Error
	})
}

func (reader *ReaderRepository) RaiseIssueRequest(ctx *gin.Context, isbn string, email string) error {
	reader.mu.Lock()
	defer reader.mu.Unlock()

	return reader.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("email = ?", email).First(&user)
		if result.Error != nil {
			return result.Error
		}

		if user.Role != util.ReaderRole {
			return errors.New("Access Denied. Provide a valid Reader email")
		}
		readerID := user.ID

		var existingBook model.BookInventory
		result = tx.Set("gorm:query_option", "FOR UPDATE").Where("isbn = ?", isbn).First(&existingBook)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Book with supplied ISBN not found in database")
		}

		if existingBook.AvailableCopies < 1 {
			return errors.New("No books available to issue")
		}

		issueRequest := model.RequestEvents{
			ReqID:        util.RandomUUID(),
			BookID:       isbn,
			ReaderID:     readerID,
			RequestDate:  time.Now(),
			ApprovalDate: sql.NullTime{},
			ApproverID:   nil,
			RequestType:  "issue",
		}
		return tx.Model(&model.RequestEvents{}).Create(&issueRequest).Error
	})
}
