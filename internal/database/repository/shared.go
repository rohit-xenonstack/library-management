package repository

import (
	"context"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SharedRepositoryInterface interface {
	AddBook(*context.Context, *model.BookInventory, string) error
	RemoveBook(*context.Context, string) error
	UpdateBook(*context.Context, *model.BookInventory) error
	ListIssueRequest(*context.Context) (*[]model.BookInventory, error)
	ApproveIssueRequest(*context.Context, string) error
	DenyIssueRequest(*context.Context, string) error
}

type SharedRepository struct {
	db        *gorm.DB
	txManager *transaction.TxManager
	mu        sync.RWMutex
}

func NewSharedRepository(db *gorm.DB, txManager *transaction.TxManager) *SharedRepository {
	return &SharedRepository{
		db:        db,
		txManager: txManager,
	}
}

func (shared *SharedRepository) SearchBookByTitle(ctx *gin.Context, title string, books *[]model.BookInventory, userID string) error {
	shared.mu.Lock()
	defer shared.mu.Unlock()

	return shared.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user)
		if result.Error != nil {
			return result.Error
		}

		query := `select * from book_inventories where lower(title) like lower('%` + title + `%') and lib_id = '` + *user.LibID + `'`
		return tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.BookInventory{}).Raw(query).Scan(&books).Error
	})
}

func (shared *SharedRepository) SearchBookByAuthor(ctx *gin.Context, author string, books *[]model.BookInventory, userID string) error {
	shared.mu.Lock()
	defer shared.mu.Unlock()

	return shared.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user)
		if result.Error != nil {
			return result.Error
		}

		query := `select * from book_inventories where lower(title) like lower('%` + author + `%') and lib_id = '` + *user.LibID + `'`
		return tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.BookInventory{}).Raw(query).Scan(&books).Error
	})
}

func (shared *SharedRepository) SearchBookByPublisher(ctx *gin.Context, publisher string, books *[]model.BookInventory, userID string) error {
	shared.mu.Lock()
	defer shared.mu.Unlock()

	return shared.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user)
		if result.Error != nil {
			return result.Error
		}

		query := `select * from book_inventories where lower(title) like lower('%` + publisher + `%') and lib_id = '` + *user.LibID + `'`
		return tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.BookInventory{}).Raw(query).Scan(&books).Error
	})
}

func (shared *SharedRepository) SearchBookByISBN(ctx *gin.Context, isbn string, book *model.BookInventory, userID string) error {
	shared.mu.Lock()
	defer shared.mu.Unlock()

	return shared.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user)
		if result.Error != nil {
			return result.Error
		}

		return tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.BookInventory{}).Where("isbn = ?", isbn).Where("lib_id = ?", user.LibID).First(&book).Error
	})
}

func (shared *SharedRepository) GetBooks(ctx *gin.Context, books *[]model.BookInventory, userID string) error {
	shared.mu.Lock()
	defer shared.mu.Unlock()

	return shared.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var user model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user)
		if result.Error != nil {
			return result.Error
		}
		return tx.Model(&model.BookInventory{}).Where("lib_id = ?", user.LibID).Find(&books).Error
	})
}
