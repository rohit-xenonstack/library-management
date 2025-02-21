package handler

import (
	"fmt"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type RemoveBookRequest struct {
	ISBN string `json:"isbn" binding:"required"`
}

type ListIssueRequestResponse struct {
	Status  string                `json:"status"`
	Payload []model.RequestEvents `json:"payload"`
}

type RequestDetails struct {
	RequestID uuid.UUID `json:"request_id" binding:"required"`
	AdminID   uuid.UUID `json:"user_id" binding:"required"`
}

type UpdateBookRequest struct {
	ISBN      string  `json:"isbn" binding:"required"`
	Title     *string `json:"title,omitempty"`
	Authors   *string `json:"authors,omitempty"`
	Publisher *string `json:"publisher,omitempty"`
	Version   *string `json:"version,omitempty"`
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
		TotalCopies:     1,
		AvailableCopies: 1,
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

func (admin *AdminHandler) RemoveBook(ctx *gin.Context) {
	var removeBookRequest RemoveBookRequest

	if err := ctx.ShouldBindJSON(&removeBookRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	err := admin.AdminRepository.RemoveBook(ctx, removeBookRequest.ISBN)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("Book removed successfuly")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}

func (admin *AdminHandler) UpdateBook(ctx *gin.Context) {
	var updateBookRequest UpdateBookRequest

	if err := ctx.ShouldBindJSON(&updateBookRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	err := admin.AdminRepository.UpdateBook(ctx, updateBookRequest.ISBN, updateBookRequest.Title, updateBookRequest.Authors, updateBookRequest.Publisher, updateBookRequest.Version)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("Book updated successfuly")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}

func (admin *AdminHandler) ListIssueRequests(ctx *gin.Context) {
	issueRequests := make([]model.RequestEvents, 0)
	err := admin.AdminRepository.ListIssueRequests(ctx, &issueRequests)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, issueRequests)
	return
}

func (admin *AdminHandler) ApproveIssueRequest(ctx *gin.Context) {
	var requestDetails RequestDetails

	if err := ctx.ShouldBindJSON(&requestDetails); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	err := admin.AdminRepository.ApproveIssueRequest(ctx, requestDetails.RequestID, requestDetails.AdminID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("Book issue request approved")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}

func (admin *AdminHandler) RejectIssueRequest(ctx *gin.Context) {
	var requestDetails RequestDetails

	if err := ctx.ShouldBindJSON(&requestDetails); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	err := admin.AdminRepository.RejectIssueRequest(ctx, requestDetails.RequestID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := fmt.Sprint("Book issue Request rejected")
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
	return
}
