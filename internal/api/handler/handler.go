package handler

import (
	"library-management/backend/internal/database/repository"
)

type Handler struct {
	AuthHandler   *AuthHandler
	OwnerHandler  *OwnerHandler
	AdminHandler  *AdminHandler
	ReaderHandler *ReaderHandler
}

func NewHandler(auth *repository.AuthRepository, owner *repository.OwnerRepository, admin *repository.AdminRepository, reader *repository.ReaderRepository) *Handler {
	return &Handler{
		AuthHandler:   NewAuthHandler(auth),
		OwnerHandler:  NewOwnerHandler(owner),
		AdminHandler:  NewAdminHandler(admin),
		ReaderHandler: NewReaderHandler(reader),
	}
}
