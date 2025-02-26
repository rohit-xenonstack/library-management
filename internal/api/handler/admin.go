package handler

import (
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/api/schema"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util/token"
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

func (admin *AdminHandler) AddBook(ctx *gin.Context) {
	var request schema.AddBookRequest
	response := schema.RequiredResponseFields{
		Status:  "error",
		Message: "",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	book := model.BookInventory{
		ISBN:            request.ISBN,
		Title:           request.Title,
		Authors:         request.Authors,
		Publisher:       request.Publisher,
		Version:         request.Version,
		TotalCopies:     1,
		AvailableCopies: 1,
	}

	err := admin.AdminRepository.AddBook(ctx, &book, request.AdminEmail)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "Book Added successfuly"
	ctx.JSON(http.StatusCreated, response)
}

func (admin *AdminHandler) RemoveBook(ctx *gin.Context) {
	var request schema.RemoveBookRequest
	response := schema.RequiredResponseFields{
		Status:  "error",
		Message: "",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		response.Message = "session not found in context"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := sessionPayload.(*token.Payload).UserID

	err := admin.AdminRepository.RemoveBook(ctx, request.ISBN, userID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "Book removed successfuly"
	ctx.JSON(http.StatusCreated, response)
}

func (admin *AdminHandler) UpdateBook(ctx *gin.Context) {
	var request schema.UpdateBookRequest
	response := schema.RequiredResponseFields{
		Status:  "error",
		Message: "",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		response.Message = "session not found in context"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := sessionPayload.(*token.Payload).UserID

	err := admin.AdminRepository.UpdateBook(ctx, request.ISBN, request.Title, request.Authors, request.Publisher, request.Version, userID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "Book updated successfuly"
	ctx.JSON(http.StatusCreated, response)
}

func (admin *AdminHandler) ListIssueRequests(ctx *gin.Context) {
	issueRequests := make([]model.IssueRequestDetails, 0)
	response := schema.ListIssueRequestResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		Requests: issueRequests,
	}

	sessionPayload, ok := ctx.Get("session_payload")
	if !ok {
		response.Message = "session not found in current context"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userID := sessionPayload.(*token.Payload).UserID

	err := admin.AdminRepository.ListIssueRequests(ctx, &issueRequests, userID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = "success"
	response.Message = "retrieved issue requests successfuly"
	response.Requests = issueRequests
	ctx.JSON(http.StatusOK, response)
}

func (admin *AdminHandler) ApproveIssueRequest(ctx *gin.Context) {
	var request schema.RequestDetails
	response := schema.RequiredResponseFields{
		Status:  "error",
		Message: "",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := admin.AdminRepository.ApproveIssueRequest(ctx, request.RequestID, request.AdminID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "book issue request approved"
	ctx.JSON(http.StatusCreated, response)
}

func (admin *AdminHandler) RejectIssueRequest(ctx *gin.Context) {
	var request schema.RequestDetails
	response := schema.RequiredResponseFields{
		Status:  "error",
		Message: "",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := admin.AdminRepository.RejectIssueRequest(ctx, request.RequestID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "book issue request rejected"
	ctx.JSON(http.StatusCreated, response)
}
