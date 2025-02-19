package repository

import (
	"errors"
	"library-management/backend/internal/api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OwnerRepositoryInterface interface {
	CreateLibrary(*model.Library, string) error
	CreateNewOwner(*model.Users) error
	OnboardAdmin(*model.Users) error
}

type OwnerRepository struct {
	db *gorm.DB
}

func NewOwnerRepository(db *gorm.DB) *OwnerRepository {
	return &OwnerRepository{
		db: db,
	}
}

func (owner *OwnerRepository) CreateLibrary(library *model.Library, userID string) error {
	var userFields model.Users
	owner.db.Clauses(clause.Locking{Strength: "SHARE"}).Where("id = ?", userID).First(&model.Users{}).Scan(&userFields)
	if userFields.LibID != nil {
		return errors.New("only one library can be create per owner")
	}

	result := owner.db.Clauses(clause.Locking{Strength: "SHARE"}).Where("name = ?", library.Name).First(&model.Library{})
	if result.RowsAffected > 0 {
		return errors.New("library with supplied name already exists in database")
	}
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	tx := owner.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(library).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.Users{}).Where("id = ?", userID).Update("lib_id", library.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (owner *OwnerRepository) AddOwner(user *model.Users) error {
	result := owner.db.Where("email = ?", user.Email).First(&model.Users{})
	if result.RowsAffected > 0 {
		return errors.New("user with supplied email already exists in database")
	}
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	result = owner.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (owner *OwnerRepository) AddAdmin(user *model.Users, libID string) error {
	result := owner.db.Where("email = ?", user.Email).First(&model.Users{})
	if result.RowsAffected > 0 {
		return errors.New("user with supplied email already exists in database")
	}
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	result = owner.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
