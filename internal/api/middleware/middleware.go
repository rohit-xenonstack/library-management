package middleware

import (
	"errors"
	"fmt"
	"library-management/backend/internal/database/repository"
	"library-management/backend/internal/util/token"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBasic  = "basic"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "session_payload"
)

type AuthMiddleware struct {
	AuthRepository *repository.AuthRepository
}

func NewAuthMiddleware(auth *repository.AuthRepository) *AuthMiddleware {
	return &AuthMiddleware{
		AuthRepository: auth,
	}
}

// type authHeader struct {
// 	BasicToken string `header:"Authorization"`
// }

func JWTAuth() gin.HandlerFunc {
	tokenMaker, _ := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

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
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func RequirePrivilege(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, ok := ctx.Get(AuthorizationPayloadKey)
		if !ok {
			err := fmt.Errorf("session not found in context")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}

		contextPayload := payload.(*token.Payload)
		if contextPayload.Role != requiredRole {
			err := fmt.Errorf("access denied. %s role required", requiredRole)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"payload": err.Error(),
			})
			return
		}

		ctx.Next()
	}
}
