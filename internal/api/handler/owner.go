package handler

import (
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/api/schema"
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

func (owner *OwnerHandler) CreateLibraryWithOwner(ctx *gin.Context) {
	var request schema.CreateLibraryRequest
	response := schema.RequiredResponseFields{
		Status:  "error",
		Message: "",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	libID := util.RandomUUID()
	ownerID := util.RandomUUID()

	newLibrary := model.Library{
		ID:   libID,
		Name: request.LibraryName,
	}

	newOwner := model.Users{
		ID:            ownerID,
		Name:          request.Name,
		Email:         request.Email,
		ContactNumber: request.ContactNumber,
		Role:          util.OwnerRole,
		LibID:         &libID,
	}

	err := owner.OwnerRepository.CreateLibraryWithUser(ctx, &newLibrary, &newOwner)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "Library created successfuly"
	ctx.JSON(http.StatusCreated, response)
}

func (owner *OwnerHandler) GetLibraries(ctx *gin.Context) {
	libraries := make([]model.LibraryDetails, 0)
	response := schema.GetLibrariesResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		Libraries: &libraries,
	}

	err := owner.OwnerRepository.GetLibraries(ctx, &libraries)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = "success"
	response.Message = "fetched libraries successfuly"
	response.Libraries = &libraries

	ctx.JSON(http.StatusOK, response)
}

func (owner *OwnerHandler) GetAdmins(ctx *gin.Context) {
	admins := make([]model.Users, 0)
	var request schema.GetAdminsRequest
	response := schema.GetAdminsResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		Admins: &admins,
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := owner.OwnerRepository.GetAdmins(ctx, &admins, request.LibID)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = "success"
	response.Message = "fetched admins for current library successfuly"
	response.Admins = &admins

	ctx.JSON(http.StatusOK, response)
}

func (owner *OwnerHandler) CreateAdmin(ctx *gin.Context) {
	var request schema.CreateAdminRequest
	response := schema.RequiredResponseFields{
		Status:  "error",
		Message: "",
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	newUser := model.Users{
		ID:            util.RandomUUID(),
		Name:          request.Name,
		Email:         request.Email,
		ContactNumber: request.ContactNumber,
		Role:          util.AdminRole,
		LibID:         &request.LibID,
	}

	err := owner.OwnerRepository.OnboardAdmin(ctx, &newUser)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "new admin onboarded successfuly"
	ctx.JSON(http.StatusCreated, response)
}
