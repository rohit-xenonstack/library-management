package token

import "time"

type Token interface {
	CreateToken(userID string, role string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
