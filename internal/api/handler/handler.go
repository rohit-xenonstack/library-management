package handler

import (
	"library-management/backend/internal/database/repository"
)

type Handler struct {
	AuthHandler   *AuthHandler
	OwnerHandler  *OwnerHandler
	AdminHandler  *AdminHandler
	ReaderHandler *ReaderHandler
	SharedHandler *SharedHandler
}

func NewHandler(auth *repository.AuthRepository, owner *repository.OwnerRepository, admin *repository.AdminRepository, reader *repository.ReaderRepository, shared *repository.SharedRepository) *Handler {
	return &Handler{
		AuthHandler:   NewAuthHandler(auth),
		OwnerHandler:  NewOwnerHandler(owner),
		AdminHandler:  NewAdminHandler(admin),
		ReaderHandler: NewReaderHandler(reader),
		SharedHandler: NewSharedHandler(shared),
	}
}
