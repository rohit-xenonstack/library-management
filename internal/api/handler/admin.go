package handler

import (
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util/token"
	"log"
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

type RemoveBookRequest struct {
	ISBN string `json:"isbn" binding:"required"`
}

type ListIssueRequestResponse struct {
	Status  string                `json:"status"`
	Payload []model.RequestEvents `json:"payload"`
}

type RequestDetails struct {
	RequestID string `json:"request_id" binding:"required"`
	AdminID   string `json:"user_id" binding:"required"`
}

type UpdateBookRequest struct {
	ISBN      string `json:"isbn" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Authors   string `json:"authors" binding:"required"`
	Publisher string `json:"publisher" binding:"required"`
	Version   string `json:"version" binding:"required"`
}

type GetBookByISBNResponse struct {
	Status  string               `json:"status"`
	Payload *model.BookInventory `json:"payload"`
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

	responseMsg := "Book Added successfuly"
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
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

	responseMsg := "Book removed successfuly"
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
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

	log.Print(updateBookRequest)
	err := admin.AdminRepository.UpdateBook(ctx, updateBookRequest.ISBN, updateBookRequest.Title, updateBookRequest.Authors, updateBookRequest.Publisher, updateBookRequest.Version)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	responseMsg := "Book updated successfuly"
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
}

func (admin *AdminHandler) ListIssueRequests(ctx *gin.Context) {
	issueRequests := make([]repository.IssueRequestDetails, 0)

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "session not found",
		})
		return
	}
	log.Print(sessionPayload)
	userID := sessionPayload.(*token.Payload).UserID
	log.Print(userID)
	err := admin.AdminRepository.ListIssueRequests(ctx, &issueRequests, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"payload": issueRequests})
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

	responseMsg := "Book issue request approved"
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
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

	responseMsg := "Book issue Request rejected"
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": responseMsg,
	})
}

func (admin *AdminHandler) SearchBook(ctx *gin.Context) {
	var searchBookRequest SearchBookRequest

	if err := ctx.ShouldBindJSON(&searchBookRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	books := make([]model.BookInventory, 0)
	switch searchBookRequest.Field {
	case "title":
		err := admin.AdminRepository.SearchBookByTitle(ctx, searchBookRequest.SearchString, &books)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}
	case "authors":
		err := admin.AdminRepository.SearchBookByAuthor(ctx, searchBookRequest.SearchString, &books)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}
	case "publisher":
		err := admin.AdminRepository.SearchBookByPublisher(ctx, searchBookRequest.SearchString, &books)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}
	default:
		err := errors.New("wrong field specified for book search")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	response := SearchBookResponse{
		Status:  "success",
		Payload: books,
	}
	ctx.JSON(http.StatusCreated, response)
}

func (admin *AdminHandler) SearchBookByISBN(ctx *gin.Context) {
	var book model.BookInventory
	isbn := ctx.Param("isbn")

	err := admin.AdminRepository.SearchBookByISBN(ctx, isbn, &book)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	response := GetBookByISBNResponse{
		Status:  "success",
		Payload: &book,
	}
	ctx.JSON(http.StatusCreated, response)
}

func (admin *AdminHandler) GetBooks(ctx *gin.Context) {
	books := make([]model.BookInventory, 0)
	err := admin.AdminRepository.GetBooks(ctx, &books)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"payload": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Books fetched successfully",
		"payload": books,
	})
}
