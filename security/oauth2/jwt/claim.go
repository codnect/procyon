package jwt

import (
	"net/url"
	"time"
)

type Claims interface {
	Id() string
	Issuer() url.URL
	Subject() string
	Audience() []string
	ExpiresAt() time.Time
	NotBefore() time.Time
	IssuedAt() time.Time
}
