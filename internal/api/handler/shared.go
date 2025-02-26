package handler

import (
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/api/schema"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SharedHandler struct {
	SharedRepository *repository.SharedRepository
}

func NewSharedHandler(shared *repository.SharedRepository) *SharedHandler {
	return &SharedHandler{
		SharedRepository: shared,
	}
}

func (shared *SharedHandler) SearchBook(ctx *gin.Context) {
	var request schema.SearchBookRequest
	books := make([]model.BookInventory, 0)
	response := schema.SearchBookResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		Books: &books,
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		response.Message = "session not found in current context"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := sessionPayload.(*token.Payload).UserID

	switch request.Field {
	case "title":
		err := shared.SharedRepository.SearchBookByTitle(ctx, request.SearchString, &books, userID)
		if err != nil {
			response.Message = err.Error()
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	case "authors":
		err := shared.SharedRepository.SearchBookByAuthor(ctx, request.SearchString, &books, userID)
		if err != nil {
			response.Message = err.Error()
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	case "publisher":
		err := shared.SharedRepository.SearchBookByPublisher(ctx, request.SearchString, &books, userID)
		if err != nil {
			response.Message = err.Error()
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	default:
		response.Message = "wrong field specified for book search"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "book search successful"
	response.Books = &books
	ctx.JSON(http.StatusOK, response)
}

func (Shared *SharedHandler) SearchBookByISBN(ctx *gin.Context) {
	var book model.BookInventory
	response := schema.SearchBookByISBNResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		Book: &book,
	}

	isbn := ctx.Param("isbn")

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		response.Message = "session not found in current context"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := sessionPayload.(*token.Payload).UserID

	err := Shared.SharedRepository.SearchBookByISBN(ctx, isbn, &book, userID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "found book with supplied ISBN"
	response.Book = &book
	ctx.JSON(http.StatusCreated, response)
}

func (Shared *SharedHandler) GetBooks(ctx *gin.Context) {
	books := make([]model.BookInventory, 0)
	response := schema.SearchBookResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		Books: &books,
	}

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		response.Message = "session not found in current context"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := sessionPayload.(*token.Payload).UserID

	err := Shared.SharedRepository.GetBooks(ctx, &books, userID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "book search successful"
	response.Books = &books
	ctx.JSON(http.StatusOK, response)
}
