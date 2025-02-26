package repository

import (
	"context"
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OwnerRepositoryInterface interface {
	CreateLibraryWithUser(*context.Context, *model.Library, *model.Users) error
	CreateUser(*context.Context, *model.Users) error
	OnboardAdmin(*context.Context, *model.Users) error
}

type OwnerRepository struct {
	db        *gorm.DB
	txManager *transaction.TxManager
	mu        sync.RWMutex
}

func NewOwnerRepository(db *gorm.DB, txManager *transaction.TxManager) *OwnerRepository {
	return &OwnerRepository{
		db:        db,
		txManager: txManager,
	}
}

func (owner *OwnerRepository) CreateLibraryWithUser(ctx context.Context, library *model.Library, user *model.Users) error {
	owner.mu.Lock()
	defer owner.mu.Unlock()

	return owner.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingUser model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("email = ?", user.Email).
			First(&existingUser)
		if result.RowsAffected > 0 && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user with supplied email already exists")
		}

		var existingLib model.Library
		result = tx.Set("gorm:query_option", "FOR UPDATE").
			Where("name = ?", library.Name).
			First(&existingLib)

		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return errors.New("library with supplied email already exists")
		}

		if err := tx.Create(library).Error; err != nil {
			return err
		}

		return tx.Model(&model.Users{}).Create(&user).Error
	})
}

func (owner *OwnerRepository) OnboardAdmin(ctx *gin.Context, user *model.Users) error {
	owner.mu.Lock()
	defer owner.mu.Unlock()

	return owner.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingUser model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("email = ?", user.Email).
			First(&existingUser)

		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}

		if result.RowsAffected > 0 {
			return errors.New("user with supplied email already exists")
		}

		var existingLibrary model.Library
		result = tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", user.LibID).
			First(&existingLibrary)

		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("no library found with given ID")
		}

		return tx.Create(user).Error
	})
}

func (owner *OwnerRepository) GetLibraries(ctx *gin.Context, libraryDetails *[]model.LibraryDetails) error {
	owner.mu.Lock()
	defer owner.mu.Unlock()

	return owner.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		query := `SELECT l.*, u.name as owner_name, u.email as owner_email, COALESCE(b.total_books, 0) as total_books
							FROM libraries l 
							LEFT JOIN users u ON l.id = u.lib_id
							LEFT JOIN ( 
								SELECT lib_id, COUNT(*) as total_books
								FROM book_inventories
								GROUP BY lib_id
							) b ON b.lib_id = l.id
							WHERE u.role = 'owner'
							`

		return tx.Raw(query).Scan(libraryDetails).Error
	})
}

func (owner *OwnerRepository) GetAdmins(ctx context.Context, admins *[]model.Users, libraryID string) error {
	owner.mu.Lock()
	defer owner.mu.Unlock()

	return owner.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		return tx.Model(&model.Users{}).Where("lib_id = ?", libraryID).Where("role = ?", "admin").Find(admins).Error
	})
}
