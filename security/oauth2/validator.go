package oauth2

type TokenValidator interface {
	Validate(token Token) error
}
