package handler

import (
	"fmt"
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
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (owner *OwnerHandler) CreateLibrary(ctx *gin.Context) {
	var createLibraryRequest CreateLibraryRequest
	if err := ctx.ShouldBindJSON(&createLibraryRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		return
	}
	libID := util.RandomUUID()
	newLibrary := model.Library{
		ID:   libID,
		Name: createLibraryRequest.Name,
	}
	fmt.Print(newLibrary)
}
