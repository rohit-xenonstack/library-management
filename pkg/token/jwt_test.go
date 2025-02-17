package token

import (
	"library-management/backend/internal/model"
	"library-management/backend/pkg/util"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWToken(t *testing.T) {
	jwtoken, err := NewJWToken(util.RandomString(32))
	assert.NoError(t, err)

	userID := util.RandomUserID()
	role := model.ReaderRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := jwtoken.CreateToken(userID, role, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)

	payload, err = jwtoken.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	assert.NotZero(t, payload.ID)
	assert.Equal(t, userID, payload.UserID)
	assert.Equal(t, role, payload.Role)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWToken(t *testing.T) {
	maker, err := NewJWToken(util.RandomString(32))
	assert.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomUserID(), model.AdminRole, -time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrExpiredToken.Error())
	assert.Nil(t, payload)
}

func TestInvalidJWTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomUserID(), model.AdminRole, time.Minute)
	assert.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	jwtoken, err := NewJWToken(util.RandomString(32))
	assert.NoError(t, err)

	payload, err = jwtoken.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidToken.Error())
	assert.Nil(t, payload)
}

func TestInvalidSecretKeySize(t *testing.T) {
	_, err := NewJWToken(util.RandomString(25))
	assert.Error(t, err)
}
