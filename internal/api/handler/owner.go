package handler

import (
	"fmt"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OwnerHandler struct {
	OwnerRepository *repository.OwnerRepository
}

func NewOwnerHandler(owner *repository.OwnerRepository) *OwnerHandler {
	return &OwnerHandler{
		OwnerRepository: owner,
	}
}

type CreateLibraryRequest struct {
	OwnerID     uuid.UUID `json:"user_id" binding:"required"`
	LibraryName string    `json:"library_name" binding:"required"`
}

type CreateOwnerRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	ContactNumber string `json:"contact" binding:"required"`
}

type CreateAdminRequest struct {
	Name          string    `json:"name" binding:"required"`
	Email         string    `json:"email" binding:"required"`
	ContactNumber string    `json:"contact" binding:"required"`
	LibID         uuid.UUID `json:"library_id" binding:"required"`
}

func (owner *OwnerHandler) CreateLibrary(ctx *gin.Context) {
	var createLibraryRequest CreateLibraryRequest

	if err := ctx.ShouldBindJSON(&createLibraryRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}
	libID := util.RandomUUID()

	newLibrary := model.Library{
		ID:   libID,
		Name: createLibraryRequest.LibraryName,
	}

	err := owner.OwnerRepository.CreateLibrary(ctx, &newLibrary, createLibraryRequest.OwnerID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("Library created successfuly")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}

func (owner *OwnerHandler) CreateOwner(ctx *gin.Context) {
	var createOwnerRequest CreateOwnerRequest
	if err := ctx.ShouldBindJSON(&createOwnerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	newUser := model.Users{
		ID:            util.RandomUUID(),
		Name:          createOwnerRequest.Name,
		Email:         createOwnerRequest.Email,
		ContactNumber: createOwnerRequest.ContactNumber,
		Role:          util.OwnerRole,
		LibID:         nil,
	}

	err := owner.OwnerRepository.CreateUser(ctx, &newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("New user with role Owner created successfuly")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}

func (owner *OwnerHandler) CreateAdmin(ctx *gin.Context) {
	var createAdminRequest CreateAdminRequest
	if err := ctx.ShouldBindJSON(&createAdminRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	newUser := model.Users{
		ID:            util.RandomUUID(),
		Name:          createAdminRequest.Name,
		Email:         createAdminRequest.Email,
		ContactNumber: createAdminRequest.ContactNumber,
		Role:          util.AdminRole,
		LibID:         &createAdminRequest.LibID,
	}

	err := owner.OwnerRepository.OnboardAdmin(ctx, &newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("New admin onboarded successfuly")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}
