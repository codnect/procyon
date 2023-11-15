package jwt

import "time"

type Jwt struct {
	header map[string]any
	claims map[string]any
}

func (j *Jwt) Value() string {
	return ""
}

func (j *Jwt) IssuedAt() time.Time {
	return time.Time{}
}

func (j *Jwt) ExpiresAt() time.Time {
	return time.Time{}
}
