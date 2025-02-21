package repository

import (
	"context"
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	Login(context.Context, string) (*model.Users, error)
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

func (auth *AuthRepository) Login(ctx context.Context, email string) (*model.Users, error) {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	var user model.Users
	err := auth.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		return tx.Set("gorm:query_option", "FOR SHARE").
			Where("email = ?", email).
			First(&user).Error
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (auth *AuthRepository) UserDetails(ctx context.Context, userID uuid.UUID) (*model.Users, error) {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	var user model.Users
	err := auth.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		return tx.Set("gorm:query_option", "FOR SHARE").
			Where("ID = ?", userID).
			First(&user).Error
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (auth *AuthRepository) UserSignup(ctx context.Context, user model.Users) error {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		var existingUser model.Users
		result := tx.Set("gorm:query_option", "FOR SHARE").Where("email = ?", user.Email).First(&existingUser)

		if result.RowsAffected > 0 {
			return errors.New("User with supplied email already exists")
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) && result.Error != nil {
			return result.Error
		}

		return tx.Model(&model.Users{}).Create(&user).Error
	})
}
