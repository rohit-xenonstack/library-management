package token

import (
	"errors"
	"library-management/backend/internal/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	IssuedAt time.Time `json:"issued_at"`
	Expires  time.Time `json:"expires"`
}

func NewPayload(userID string, role string, duration time.Duration) (*Payload, error) {
	tokenID := util.RandomUUID()

	payload := &Payload{
		ID:       tokenID,
		UserID:   userID,
		Role:     role,
		IssuedAt: time.Now(),
		Expires:  time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.Expires) {
		return ErrExpiredToken
	}
	return nil
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.Expires,
	}, nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetIssuer() (string, error) {
	return "", nil
}

func (payload *Payload) GetSubject() (string, error) {
	return "", nil
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
