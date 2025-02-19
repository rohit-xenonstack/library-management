package repository

import (
	"library-management/backend/internal/database/transaction"

	"gorm.io/gorm"
)

type Repository struct {
	AuthRepository  *AuthRepository
	OwnerRepository *OwnerRepository
	txManager       *transaction.TxManager
}

func NewRepository(db *gorm.DB) *Repository {
	txManager := transaction.NewTxManager(db)
	return &Repository{
		AuthRepository:  NewAuthRepository(db, txManager),
		OwnerRepository: NewOwnerRepository(db, txManager),
		txManager:       txManager,
	}
}
