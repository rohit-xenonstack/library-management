package repository

import (
	"context"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"sync"

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
