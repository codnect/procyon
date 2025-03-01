package http

// Method represents an HTTP method.
type Method string

const (
	// MethodGet represents the GET HTTP method.
	MethodGet Method = "GET"
	// MethodHead represents the HEAD HTTP method.
	MethodHead Method = "HEAD"
	// MethodPost represents the POST HTTP method.
	MethodPost Method = "POST"
	// MethodPut represents the PUT HTTP method.
	MethodPut Method = "PUT"
	// MethodPatch represents the PATCH HTTP method.
	MethodPatch Method = "PATCH" // RFC 5789
	// MethodDelete represents the DELETE HTTP method.
	MethodDelete Method = "DELETE"
	// MethodOptions represents the OPTIONS HTTP method.
	MethodOptions Method = "OPTIONS"
	// MethodTrace represents the TRACE HTTP method.
	MethodTrace Method = "TRACE"
)

func httpMethodToInt(method Method) int {
	switch method {
	case MethodGet:
		return 0
	case MethodHead:
		return 1
	case MethodPost:
		return 2
	case MethodPut:
		return 3
	case MethodPatch:
		return 4
	case MethodDelete:
		return 5
	case MethodOptions:
		return 6
	case MethodTrace:
		return 7
	default:
		return -1
	}
}
