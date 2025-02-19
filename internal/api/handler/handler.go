package handler

import (
	"library-management/backend/internal/database/repository"
)

type Handler struct {
	AuthHandler  *AuthHandler
	OwnerHandler *OwnerHandler
}

func NewHandler(auth *repository.AuthRepository, owner *repository.OwnerRepository) *Handler {
	return &Handler{
		AuthHandler:  NewAuthHandler(auth),
		OwnerHandler: NewOwnerHandler(owner),
	}
}
