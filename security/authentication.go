package security

type Authentication interface {
	Authorities() []GrantedAuthority
	Details() any
	Principal() any
	IsAuthenticated() bool
	SetAuthenticated(isAuthenticated bool) error
}
