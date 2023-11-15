package oauth2

import "time"

type RefreshToken struct {
}

func NewRefreshToken(tokenValue string, issuedAt time.Time, expiresAt time.Time) *AccessToken {
	return &AccessToken{}
}

func (t *RefreshToken) Value() string {
	return ""
}

func (t *RefreshToken) IssuedAt() time.Time {
	return time.Time{}
}

func (t *RefreshToken) ExpiresAt() time.Time {
	return time.Time{}
}
