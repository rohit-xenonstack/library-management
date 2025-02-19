package handler

import (
	"errors"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util/token"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// type LoginResponse struct {
// 	AccessToken    string `json:"access_token" binding:"required"`
// 	RefreshToken   string `json:"refresh_token" binding:"required"`
// 	AccessPayload  util.Payload
// 	RefreshPayload util.Payload
// }

type AccessToken struct {
	AccessToken string `json:"access_token" binding:"required"`
}
type LoginResponse struct {
	Status  string      `json:"status" binding:"required"`
	Payload AccessToken `json:"payload"`
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
		Payload: AccessToken{
			AccessToken: accessToken,
		},
	}

	ctx.IndentedJSON(http.StatusOK, response)
	return
}
