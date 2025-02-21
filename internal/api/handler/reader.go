package handler

import (
	"errors"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
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
		err := reader.ReaderRepository.SearchBookByTitle(ctx, searchBookRequest.SearchString, &books)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}
	case "publisher":
		err := reader.ReaderRepository.SearchBookByTitle(ctx, searchBookRequest.SearchString, &books)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}
	default:
		err := errors.New("Wrong field specified for book search")
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
	return
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

	err := reader.ReaderRepository.RaiseIssueRequest(ctx, raiseIssueRequest.BookID, raiseIssueRequest.ReaderEmail)
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
	return
}
