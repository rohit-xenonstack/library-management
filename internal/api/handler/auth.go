package handler

import (
	"errors"
	"fmt"
	"library-management/backend/internal/api/middleware"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util"
	"library-management/backend/internal/util/token"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
}

// Update the UserSignupRequest struct to implement custom unmarshalling
type UserSignupRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Contact string `json:"contact" binding:"required"`
	LibID   string `json:"library_id" binding:"required"`
}

type UserSignupResponse struct {
	Status  string `json:"status" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}

type LoginPayload struct {
	Message     string       `json:"message" binding:"required"`
	AccessToken *string      `json:"access_token,omitempty"`
	User        *model.Users `json:"user,omitempty"`
}

type LoginResponse struct {
	Status  string       `json:"status" binding:"required"`
	Payload LoginPayload `json:"payload" binding:"required"`
}

type RefreshPayload struct {
	AccessToken string      `json:"access_token,omitempty"`
	Message     string      `json:"message,omitempty"`
	User        model.Users `json:"user,omitempty"`
}

type RefreshTokenResponse struct {
	Status  string         `json:"status"`
	Payload RefreshPayload `json:"payload"`
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
	var loginResponse LoginResponse

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		loginResponse.Status = "error"
		loginResponse.Payload.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, loginResponse)
		return
	}

	user, err := auth.AuthRepository.Login(ctx, loginRequest.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			loginResponse.Status = "error"
			loginResponse.Payload.Message = "user not found with given email"
			ctx.JSON(http.StatusNotFound, loginResponse)
			return
		}
		loginResponse.Status = "error"
		loginResponse.Payload.Message = "internal server error"
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	jwtoken, err := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		loginResponse.Status = "error"
		loginResponse.Payload.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		loginResponse.Status = "error"
		loginResponse.Payload.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	accessToken, _, err := jwtoken.CreateToken(user.ID, user.Role, duration)
	if err != nil {
		loginResponse.Status = "error"
		loginResponse.Payload.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	response := LoginResponse{
		Status: "success",
		Payload: LoginPayload{
			AccessToken: &accessToken,
			User:        user,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func (auth *AuthHandler) UserDetails(ctx *gin.Context) {
	session, exists := ctx.Get(middleware.AuthorizationPayloadKey)
	sessionPayload := session.(*token.Payload)
	if !exists {
		err := fmt.Errorf("session not found in current context")
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
}

func (auth *AuthHandler) RefreshToken(ctx *gin.Context) {
	var response RefreshTokenResponse

	authorizationHeader := ctx.GetHeader("Authorization")
	log.Print(authorizationHeader)
	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header is missing")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		err := errors.New("authorization header is invalid")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != "bearer" {
		err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"payload": err.Error(),
		})
		return
	}

	accessToken := fields[1]
	log.Print(accessToken)

	tokenMaker, err := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		response.Status = "error"
		response.Payload.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	payload, err := tokenMaker.VerifyToken(accessToken)
	if err != nil && !errors.Is(err, token.ErrExpiredToken) {
		response.Status = "error"
		response.Payload.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		response.Status = "error"
		response.Payload.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	log.Print("payload", payload)
	newAccessToken, _, err := tokenMaker.CreateToken(payload.UserID, payload.Role, duration)
	if err != nil {
		response.Status = "error"
		response.Payload.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = "success"
	response.Payload.AccessToken = newAccessToken
	ctx.JSON(http.StatusOK, response)
}
