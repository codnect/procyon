package security

type GrantedAuthority interface {
	Authority() string
}
