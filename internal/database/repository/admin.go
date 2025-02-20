package repository

import (
	"context"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"sync"

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

func (admin *AdminRepository) AddBook(ctx context.Context, book *model.BookInventory, adminEmail string) error {
	admin.mu.Lock()
	defer admin.mu.Unlock()

	return admin.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var admin model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("email = ?", adminEmail).
			First(&admin)
		if result.Error != nil {
			return result.Error
		}

		book.LibID = admin.LibID

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
