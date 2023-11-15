package jwt

type Decoder interface {
	Decode(token string) (*Jwt, error)
}
