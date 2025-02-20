package handler

import (
	"fmt"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminRepository *repository.AdminRepository
}

func NewAdminHandler(admin *repository.AdminRepository) *AdminHandler {
	return &AdminHandler{
		AdminRepository: admin,
	}
}

type AddBookRequest struct {
	AdminEmail string `json:"email" binding:"required"`
	ISBN       string `json:"isbn" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Authors    string `json:"authors" binding:"required"`
	Publisher  string `json:"publisher" binding:"required"`
	Version    string `json:"version" binding:"required"`
}

func (admin *AdminHandler) AddBook(ctx *gin.Context) {
	var addBookRequest AddBookRequest

	if err := ctx.ShouldBindJSON(&addBookRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	book := model.BookInventory{
		ISBN:            addBookRequest.ISBN,
		Title:           addBookRequest.Title,
		Authors:         addBookRequest.Authors,
		Publisher:       addBookRequest.Publisher,
		Version:         addBookRequest.Version,
		TotalCopies:     0,
		AvailableCopies: 0,
	}
	err := admin.AdminRepository.AddBook(ctx, &book, addBookRequest.AdminEmail)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("Book Added successfuly")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}
