package oauth2

type AuthorizationGrantType string

const (
	GrantTypeRefreshToken      AuthorizationGrantType = "refresh_token"
	GrantTypeClientCredentials AuthorizationGrantType = "client_credentials"
	GrantTypePassword          AuthorizationGrantType = "password"
)
