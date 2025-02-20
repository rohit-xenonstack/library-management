package middleware_test

import (
	"fmt"
	"library-management/backend/internal/api"
	"library-management/backend/internal/api/handler"
	"library-management/backend/internal/api/middleware"
	"library-management/backend/internal/config"
	"library-management/backend/internal/util"
	"library-management/backend/internal/util/token"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	authorizationType string,
	userID uuid.UUID,
	role string,
	duration time.Duration,
) {
	jwtoken, _ := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	token, payload, err := jwtoken.CreateToken(userID, role, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(middleware.AuthorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	userID := util.RandomUUID()
	role := util.OwnerRole

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "unsupported", userID, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "", userID, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, role, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	cfg := &config.SampleEnv
	h := handler.NewHandler(nil, nil, nil)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			api := api.NewAPI(cfg, h)

			authPath := "/auth/login"
			api.Router.GET(
				authPath,
				middleware.JWTAuth(),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			assert.NoError(t, err)

			tc.setupAuth(t, request)
			api.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
