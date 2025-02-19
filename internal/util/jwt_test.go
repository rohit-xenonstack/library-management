package util

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWToken(t *testing.T) {
	jwtoken, err := NewJWTMaker(RandomString(32))
	assert.NoError(t, err)

	email := RandomEmail()
	role := ReaderRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := jwtoken.CreateToken(email, role, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)

	payload, err = jwtoken.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	assert.NotZero(t, payload.ID)
	assert.Equal(t, email, payload.Email)
	assert.Equal(t, role, payload.Role)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWToken(t *testing.T) {
	maker, err := NewJWTMaker(RandomString(32))
	assert.NoError(t, err)

	token, payload, err := maker.CreateToken(RandomEmail(), AdminRole, -time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrExpiredToken.Error())
	assert.Nil(t, payload)
}

func TestInvalidJWTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(RandomEmail(), AdminRole, time.Minute)
	assert.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	jwtoken, err := NewJWTMaker(RandomString(32))
	assert.NoError(t, err)

	payload, err = jwtoken.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidToken.Error())
	assert.Nil(t, payload)
}

func TestInvalidSecretKeySize(t *testing.T) {
	_, err := NewJWTMaker(RandomString(25))
	assert.Error(t, err)
}
