package oauth2

import "time"

type AccessToken struct {
}

func NewAccessToken(tokenType TokenType, tokenValue string, issuedAt time.Time, expiresAt time.Time, scopes []string) *AccessToken {
	return &AccessToken{}
}

func (t *AccessToken) Value() string {
	return ""
}

func (t *AccessToken) IssuedAt() time.Time {
	return time.Time{}
}

func (t *AccessToken) ExpiresAt() time.Time {
	return time.Time{}
}

func (t *AccessToken) TokenType() TokenType {
	return ""
}

func (t *AccessToken) Scopes() []string {
	return nil
}
