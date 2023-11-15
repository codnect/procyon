package oauth2

import "time"

type TokenType string

const (
	BearerToken TokenType = "Bearer"
)

type Token interface {
	Value() string
	IssuedAt() time.Time
	ExpiresAt() time.Time
}
