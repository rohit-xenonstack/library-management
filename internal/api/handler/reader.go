package handler

import (
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReaderHandler struct {
	ReaderRepository *repository.ReaderRepository
}

func NewReaderHandler(reader *repository.ReaderRepository) *ReaderHandler {
	return &ReaderHandler{
		ReaderRepository: reader,
	}
}

type SearchBookRequest struct {
	SearchString string `json:"search_string"`
	Field        string `json:"field"`
}

type SearchBookResponse struct {
	Status  string                `json:"status"`
	Payload []model.BookInventory `json:"payload"`
}

type RaiseIssueRequest struct {
	ReaderEmail string `json:"email" binding:"required"`
	BookID      string `json:"isbn" binding:"required"`
}

func (reader *ReaderHandler) SearchBook(ctx *gin.Context) {
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
		err := reader.ReaderRepository.SearchBookByTitle(ctx, searchBookRequest.SearchString, &books)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}
	case "authors":
		err := reader.ReaderRepository.SearchBookByAuthor(ctx, searchBookRequest.SearchString, &books)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}
	case "publisher":
		err := reader.ReaderRepository.SearchBookByPublisher(ctx, searchBookRequest.SearchString, &books)
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

func (reader *ReaderHandler) RaiseIssueRequest(ctx *gin.Context) {
	var raiseIssueRequest RaiseIssueRequest
	if err := ctx.ShouldBindJSON(&raiseIssueRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "invalid request parameters",
		})
		return
	}

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": "session not found",
		})
		return
	}

	userID := sessionPayload.(*token.Payload).UserID
	err := reader.ReaderRepository.RaiseIssueRequest(ctx, raiseIssueRequest.BookID, raiseIssueRequest.ReaderEmail, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": "Raised Issue Request successfuly",
	})
}

func (reader *ReaderHandler) GetLatestAvailability(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	var latestDate string
	err := reader.ReaderRepository.GetLatestBookAvailability(ctx, isbn, &latestDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"payload": latestDate,
	})
}

func (reader *ReaderHandler) GetLatestBooks(ctx *gin.Context) {
	books := make([]model.BookInventory, 0)
	err := reader.ReaderRepository.GetBooks(ctx, &books)
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
