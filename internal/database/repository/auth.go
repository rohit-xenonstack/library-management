package repository

import (
	"context"
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"sync"

	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	Login(context.Context, string) (*model.Users, error)
	UserDetails(context.Context, string) (*model.Users, error)
	UserSignup(context.Context, model.Users) error
}

type AuthRepository struct {
	DB        *gorm.DB
	txManager *transaction.TxManager
	mu        sync.RWMutex
}

func NewAuthRepository(db *gorm.DB, txManager *transaction.TxManager) *AuthRepository {
	return &AuthRepository{
		DB:        db,
		txManager: txManager,
	}
}

func (auth *AuthRepository) Login(ctx context.Context, email string, user *model.Users) error {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		return tx.Set("gorm:query_option", "FOR SHARE").
			Where("email = ?", email).
			First(&user).Error
	})
}

func (auth *AuthRepository) UserDetails(ctx context.Context, userID string, user *model.Users) error {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		return tx.Set("gorm:query_option", "FOR SHARE").
			Where("ID = ?", userID).
			First(&user).Error
	})
}

func (auth *AuthRepository) UserSignup(ctx context.Context, user model.Users) error {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingLibrary model.Library
		result := tx.Set("gorm:query_option", "FOR SHARE").Where("ID = ?", user.LibID).First(&existingLibrary)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("library not found")
			}
			return result.Error
		}

		var existingUser model.Users
		result = tx.Set("gorm:query_option", "FOR SHARE").Where("email = ?", user.Email).First(&existingUser)

		if result.RowsAffected > 0 {
			return errors.New("user with supplied email already exists")
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) && result.Error != nil {
			return result.Error
		}

		return tx.Model(&model.Users{}).Create(&user).Error
	})
}
