package handler

import (
	"errors"
	"fmt"
	"library-management/backend/internal/api/middleware"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util"
	"library-management/backend/internal/util/token"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
}

type UserSignupRequest struct {
	Name    string    `json:"name" binding:"required"`
	Email   string    `json:"email" binding:"required"`
	Contact string    `json:"contact" binding:"required"`
	LibID   uuid.UUID `json:"library_id" binding:"required"`
}

type UserSignupResponse struct {
	Status  string `json:"status" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}

type LoginPayload struct {
	AccessToken string `json:"access_token" binding:"required"`
	Role        string `json:"role" binding:"required"`
}

type LoginResponse struct {
	Status  string       `json:"status" binding:"required"`
	Payload LoginPayload `json:"payload" binding:"required"`
}

type AuthHandler struct {
	AuthRepository *repository.AuthRepository
}

func NewAuthHandler(auth *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{
		AuthRepository: auth,
	}
}

func (auth *AuthHandler) Login(ctx *gin.Context) {
	var loginRequest LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	user, err := auth.AuthRepository.Login(ctx, loginRequest.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found with given email"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtoken, err := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))

	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, _, err := jwtoken.CreateToken(user.ID, user.Role, duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := LoginResponse{
		Status: "success",
		Payload: LoginPayload{
			AccessToken: accessToken,
			Role:        user.Role,
		},
	}

	ctx.JSON(http.StatusOK, response)
	return
}

func (auth *AuthHandler) UserDetails(ctx *gin.Context) {
	session, exists := ctx.Get(middleware.AuthorizationPayloadKey)
	sessionPayload := session.(*token.Payload)
	if !exists {
		err := fmt.Errorf("Session not found in current context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := auth.AuthRepository.UserDetails(ctx, sessionPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (auth *AuthHandler) UserSignup(ctx *gin.Context) {
	var userSignupRequest UserSignupRequest

	if err := ctx.ShouldBindJSON(&userSignupRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	newUser := model.Users{
		ID:            util.RandomUUID(),
		Name:          userSignupRequest.Name,
		Email:         userSignupRequest.Email,
		ContactNumber: userSignupRequest.Contact,
		Role:          util.ReaderRole,
		LibID:         &userSignupRequest.LibID,
	}

	err := auth.AuthRepository.UserSignup(ctx, newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := UserSignupResponse{
		Status:  "success",
		Payload: "User Signup Successful",
	}
	ctx.JSON(http.StatusOK, response)
	return
}
