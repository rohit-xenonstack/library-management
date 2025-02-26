package handler

import (
	"errors"
	"library-management/backend/internal/api/middleware"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/api/schema"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util"
	"library-management/backend/internal/util/token"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	AuthRepository *repository.AuthRepository
}

func NewAuthHandler(auth *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{
		AuthRepository: auth,
	}
}

func (auth *AuthHandler) Login(ctx *gin.Context) {
	var loginRequest schema.LoginRequest
	loginResponse := schema.LoginResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		AccessToken: nil,
		User:        nil,
	}

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		loginResponse.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, loginResponse)
		return
	}

	var user model.Users
	err := auth.AuthRepository.Login(ctx, loginRequest.Email, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			loginResponse.Message = "user not found with given email"
			ctx.JSON(http.StatusNotFound, loginResponse)
			return
		}
		loginResponse.Message = "internal server error"
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	jwtoken, err := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		loginResponse.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		loginResponse.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	accessToken, _, err := jwtoken.CreateToken(user.ID, user.Role, duration)
	if err != nil {
		loginResponse.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, loginResponse)
		return
	}

	loginResponse.Status = "success"
	loginResponse.Message = "login successful"
	loginResponse.AccessToken = &accessToken
	loginResponse.User = &user
	ctx.JSON(http.StatusOK, loginResponse)
}

func (auth *AuthHandler) UserDetails(ctx *gin.Context) {
	var user model.Users
	response := schema.UserDetailsResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		User: &user,
	}

	session, exists := ctx.Get(middleware.AuthorizationPayloadKey)
	if !exists {
		response.Message = "session not found in current context"
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	sessionPayload := session.(*token.Payload)
	err := auth.AuthRepository.UserDetails(ctx, sessionPayload.UserID, &user)
	if err != nil {
		response.Message = "internal server error"
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = "success"
	response.Message = "user details fetched successfully"
	response.User = &user
	ctx.JSON(http.StatusOK, response)
}

func (auth *AuthHandler) ReaderSignup(ctx *gin.Context) {
	var request schema.ReaderSignupRequest
	response := schema.ReaderSignupResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Message = "invalid request parameters"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	newUser := model.Users{
		ID:            util.RandomUUID(),
		Name:          request.Name,
		Email:         request.Email,
		ContactNumber: request.Contact,
		Role:          util.ReaderRole,
		LibID:         &request.LibID,
	}

	err := auth.AuthRepository.UserSignup(ctx, newUser)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.Message = "user signed up successfully"
	ctx.JSON(http.StatusOK, response)
}

func (auth *AuthHandler) RefreshAccessToken(ctx *gin.Context) {
	response := schema.RefreshAccessTokenResponse{
		RequiredResponseFields: schema.RequiredResponseFields{
			Status:  "error",
			Message: "",
		},
		AccessToken: nil,
	}

	authorizationHeader := ctx.GetHeader("Authorization")
	if len(authorizationHeader) == 0 {
		response.Message = "authorization header not provided"
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		response.Message = "invalid authorization header format"
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != "bearer" {
		response.Message = `unsupported authorization type, only "bearer" supported`
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	accessToken := fields[1]

	tokenMaker, err := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	payload, err := tokenMaker.VerifyToken(accessToken)
	if err != nil && !errors.Is(err, token.ErrExpiredToken) {
		response.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	newAccessToken, _, err := tokenMaker.CreateToken(payload.UserID, payload.Role, duration)
	if err != nil {
		response.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = "success"
	response.Message = "token refreshed successfully"
	response.AccessToken = &newAccessToken
	ctx.JSON(http.StatusOK, response)
}
