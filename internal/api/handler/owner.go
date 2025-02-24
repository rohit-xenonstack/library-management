package handler

import (
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
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
	LibraryName   string `json:"library_name" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	ContactNumber string `json:"contact" binding:"required"`
}

type CreateAdminRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	ContactNumber string `json:"contact" binding:"required"`
	LibID         string `json:"library_id" binding:"required"`
}

type GetAdminsRequest struct {
	LibID string `json:"library_id" binding:"required"`
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
	ownerID := util.RandomUUID()

	newLibrary := model.Library{
		ID:   libID,
		Name: createLibraryRequest.LibraryName,
	}

	newOwner := model.Users{
		ID:            ownerID,
		Name:          createLibraryRequest.Name,
		Email:         createLibraryRequest.Email,
		ContactNumber: createLibraryRequest.ContactNumber,
		Role:          util.OwnerRole,
		LibID:         &libID,
	}

	err := owner.OwnerRepository.CreateLibraryWithUser(ctx, &newLibrary, &newOwner)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := "Library created successfuly"
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
}

func (owner *OwnerHandler) GetLibraries(ctx *gin.Context) {
	libraries := make([]model.LibraryDetails, 0)
	err := owner.OwnerRepository.GetLibraries(ctx, &libraries)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"payload": "Failed to fetch libraries",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"libraries": libraries,
	})
}

func (owner *OwnerHandler) GetAdmins(ctx *gin.Context) {
	var getAdminsRequest GetAdminsRequest
	if err := ctx.ShouldBindJSON(&getAdminsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	admins := make([]model.Users, 0)
	err := owner.OwnerRepository.GetAdmins(ctx, &admins, getAdminsRequest.LibID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"payload": "Failed to fetch admins",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"admins": admins,
	})
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

	responseMsg := "New admin onboarded successfuly"
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
}
