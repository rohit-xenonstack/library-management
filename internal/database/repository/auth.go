package repository

import (
	"library-management/backend/internal/api/model"

	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	Login(string) (*model.Users, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (auth *AuthRepository) Login(email string) (*model.Users, error) {
	result := auth.db.Where("email = ?", email).First(&model.Users{})

	if result.Error != nil {
		return nil, result.Error
	}

	var user *model.Users
	result.Scan(&user)
	return user, nil
}
