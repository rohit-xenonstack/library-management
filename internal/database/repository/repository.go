package repository

import "gorm.io/gorm"

type Repository struct {
	AuthRepository  *AuthRepository
	OwnerRepository *OwnerRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AuthRepository:  NewAuthRepository(db),
		OwnerRepository: NewOwnerRepository(db),
	}
}
