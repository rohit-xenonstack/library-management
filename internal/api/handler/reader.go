package handler

import (
	"library-management/backend/internal/api/schema"
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

func (reader *ReaderHandler) RaiseIssueRequest(ctx *gin.Context) {
	var request schema.RaiseIssueRequest
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
	err := reader.ReaderRepository.RaiseIssueRequest(ctx, request.BookID, request.ReaderEmail, userID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "Raised Issue Request successfuly"
	ctx.JSON(http.StatusCreated, response)
}

func (reader *ReaderHandler) GetLatestAvailability(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	var latestDate string
	response := schema.GetLatestAvailabilityResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		Date: &latestDate,
	}

	err := reader.ReaderRepository.GetLatestBookAvailability(ctx, isbn, &latestDate)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "Fetched latest available date successfuly"
	response.Date = &latestDate
	ctx.JSON(http.StatusCreated, response)
}
