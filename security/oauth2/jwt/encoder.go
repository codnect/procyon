package jwt

type Encoder interface {
	Encode() (*Jwt, error)
}
