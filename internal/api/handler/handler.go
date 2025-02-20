package handler

import (
	"library-management/backend/internal/database/repository"
)

type Handler struct {
	AuthHandler  *AuthHandler
	OwnerHandler *OwnerHandler
	AdminHandler *AdminHandler
}

func NewHandler(auth *repository.AuthRepository, owner *repository.OwnerRepository, admin *repository.AdminRepository) *Handler {
	return &Handler{
		AuthHandler:  NewAuthHandler(auth),
		OwnerHandler: NewOwnerHandler(owner),
		AdminHandler: NewAdminHandler(admin),
	}
}
