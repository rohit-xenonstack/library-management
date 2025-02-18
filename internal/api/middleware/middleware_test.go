package middleware

import (
	"fmt"
	"library-management/backend/internal/api"
	"library-management/backend/internal/config"
	"library-management/backend/internal/model"
	"library-management/backend/pkg/token"
	"library-management/backend/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var cfg = &config.SampleEnv
var db *gorm.DB

func addAuthorization(
	t *testing.T,
	request *http.Request,
	jwtoken token.Token,
	authorizationType string,
	username string,
	role string,
	duration time.Duration,
) {
	token, payload, err := jwtoken.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	userID := util.RandomUserID()
	role := model.OwnerRole

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, token token.Token)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, token token.Token) {
				addAuthorization(t, request, token, authorizationTypeBearer, userID, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, token token.Token) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, token token.Token) {
				addAuthorization(t, request, token, "unsupported", userID, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, token token.Token) {
				addAuthorization(t, request, token, "", userID, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, token token.Token) {
				addAuthorization(t, request, token, authorizationTypeBearer, userID, role, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			api, err := api.NewAPI(cfg, db)
			assert.NoError(t, err)

			authPath := "/auth"
			api.Router.GET(
				authPath,
				authMiddleware(api.Token),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			assert.NoError(t, err)

			tc.setupAuth(t, request, api.Token)
			api.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
