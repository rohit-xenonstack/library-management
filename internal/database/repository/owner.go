package repository

import (
	"context"
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OwnerRepositoryInterface interface {
	CreateLibrary(*context.Context, *model.Library, string) error
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

func (owner *OwnerRepository) CreateLibrary(ctx context.Context, library *model.Library, userID uuid.UUID) error {
	owner.mu.Lock()
	defer owner.mu.Unlock()

	return owner.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var userFields model.Users
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", userID).
			First(&userFields).Error; err != nil {
			return err
		}

		if userFields.LibID != nil {
			return errors.New("only one library can be created per owner")
		}

		var existingLib model.Library
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("name = ?", library.Name).
			First(&existingLib)

		if result.RowsAffected > 0 {
			return errors.New("library with supplied email already exists")
		}

		if err := tx.Create(library).Error; err != nil {
			return err
		}

		return tx.Model(&model.Users{}).
			Where("id = ?", userID).
			Update("lib_id", library.ID).Error
	})
}

func (owner *OwnerRepository) CreateUser(ctx *gin.Context, user *model.Users) error {
	owner.mu.Lock()
	defer owner.mu.Unlock()

	return owner.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingUser model.Users
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("email = ?", user.Email).
			First(&existingUser)

		if result.RowsAffected > 0 {
			return errors.New("user with supplied email already exists")
		}

		return tx.Create(user).Error
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
		if result.RowsAffected > 0 {
			return errors.New("user with supplied email already exists")
		}

		var existingLibrary model.Library
		result = tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", user.LibID).
			First(&existingLibrary)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("No library found with given ID")
		}

		return tx.Create(user).Error
	})
}
