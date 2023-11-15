package http

type Status int

// HTTP status codes as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (
	StatusContinue           Status = 100 // RFC 9110, 15.2.1
	StatusSwitchingProtocols Status = 101 // RFC 9110, 15.2.2
	StatusProcessing         Status = 102 // RFC 2518, 10.1
	StatusEarlyHints         Status = 103 // RFC 8297

	StatusOK                   Status = 200 // RFC 9110, 15.3.1
	StatusCreated              Status = 201 // RFC 9110, 15.3.2
	StatusAccepted             Status = 202 // RFC 9110, 15.3.3
	StatusNonAuthoritativeInfo Status = 203 // RFC 9110, 15.3.4
	StatusNoContent            Status = 204 // RFC 9110, 15.3.5
	StatusResetContent         Status = 205 // RFC 9110, 15.3.6
	StatusPartialContent       Status = 206 // RFC 9110, 15.3.7
	StatusMultiStatus          Status = 207 // RFC 4918, 11.1
	StatusAlreadyReported      Status = 208 // RFC 5842, 7.1
	StatusIMUsed               Status = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   Status = 300 // RFC 9110, 15.4.1
	StatusMovedPermanently  Status = 301 // RFC 9110, 15.4.2
	StatusFound             Status = 302 // RFC 9110, 15.4.3
	StatusSeeOther          Status = 303 // RFC 9110, 15.4.4
	StatusNotModified       Status = 304 // RFC 9110, 15.4.5
	StatusUseProxy          Status = 305 // RFC 9110, 15.4.6
	_                       Status = 306 // RFC 9110, 15.4.7 (Unused)
	StatusTemporaryRedirect Status = 307 // RFC 9110, 15.4.8
	StatusPermanentRedirect Status = 308 // RFC 9110, 15.4.9

	StatusBadRequest                   Status = 400 // RFC 9110, 15.5.1
	StatusUnauthorized                 Status = 401 // RFC 9110, 15.5.2
	StatusPaymentRequired              Status = 402 // RFC 9110, 15.5.3
	StatusForbidden                    Status = 403 // RFC 9110, 15.5.4
	StatusNotFound                     Status = 404 // RFC 9110, 15.5.5
	StatusMethodNotAllowed             Status = 405 // RFC 9110, 15.5.6
	StatusNotAcceptable                Status = 406 // RFC 9110, 15.5.7
	StatusProxyAuthRequired            Status = 407 // RFC 9110, 15.5.8
	StatusRequestTimeout               Status = 408 // RFC 9110, 15.5.9
	StatusConflict                     Status = 409 // RFC 9110, 15.5.10
	StatusGone                         Status = 410 // RFC 9110, 15.5.11
	StatusLengthRequired               Status = 411 // RFC 9110, 15.5.12
	StatusPreconditionFailed           Status = 412 // RFC 9110, 15.5.13
	StatusRequestEntityTooLarge        Status = 413 // RFC 9110, 15.5.14
	StatusRequestURITooLong            Status = 414 // RFC 9110, 15.5.15
	StatusUnsupportedMediaType         Status = 415 // RFC 9110, 15.5.16
	StatusRequestedRangeNotSatisfiable Status = 416 // RFC 9110, 15.5.17
	StatusExpectationFailed            Status = 417 // RFC 9110, 15.5.18
	StatusTeapot                       Status = 418 // RFC 9110, 15.5.19 (Unused)
	StatusMisdirectedRequest           Status = 421 // RFC 9110, 15.5.20
	StatusUnprocessableEntity          Status = 422 // RFC 9110, 15.5.21
	StatusLocked                       Status = 423 // RFC 4918, 11.3
	StatusFailedDependency             Status = 424 // RFC 4918, 11.4
	StatusTooEarly                     Status = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              Status = 426 // RFC 9110, 15.5.22
	StatusPreconditionRequired         Status = 428 // RFC 6585, 3
	StatusTooManyRequests              Status = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  Status = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   Status = 451 // RFC 7725, 3

	StatusInternalServerError           Status = 500 // RFC 9110, 15.6.1
	StatusNotImplemented                Status = 501 // RFC 9110, 15.6.2
	StatusBadGateway                    Status = 502 // RFC 9110, 15.6.3
	StatusServiceUnavailable            Status = 503 // RFC 9110, 15.6.4
	StatusGatewayTimeout                Status = 504 // RFC 9110, 15.6.5
	StatusHTTPVersionNotSupported       Status = 505 // RFC 9110, 15.6.6
	StatusVariantAlsoNegotiates         Status = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           Status = 507 // RFC 4918, 11.5
	StatusLoopDetected                  Status = 508 // RFC 5842, 7.2
	StatusNotExtended                   Status = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired Status = 511 // RFC 6585, 6
)

func (code Status) String() string {
	switch code {
	case StatusContinue:
		return "Continue"
	case StatusSwitchingProtocols:
		return "Switching Protocols"
	case StatusProcessing:
		return "Processing"
	case StatusEarlyHints:
		return "Early Hints"
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusAccepted:
		return "Accepted"
	case StatusNonAuthoritativeInfo:
		return "Non-Authoritative Information"
	case StatusNoContent:
		return "No Content"
	case StatusResetContent:
		return "Reset Content"
	case StatusPartialContent:
		return "Partial Content"
	case StatusMultiStatus:
		return "Multi-Status"
	case StatusAlreadyReported:
		return "Already Reported"
	case StatusIMUsed:
		return "IM Used"
	case StatusMultipleChoices:
		return "Multiple Choices"
	case StatusMovedPermanently:
		return "Moved Permanently"
	case StatusFound:
		return "Found"
	case StatusSeeOther:
		return "See Other"
	case StatusNotModified:
		return "Not Modified"
	case StatusUseProxy:
		return "Use Proxy"
	case StatusTemporaryRedirect:
		return "Temporary Redirect"
	case StatusPermanentRedirect:
		return "Permanent Redirect"
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusPaymentRequired:
		return "Payment Required"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusNotAcceptable:
		return "Not Acceptable"
	case StatusProxyAuthRequired:
		return "Proxy Authentication Required"
	case StatusRequestTimeout:
		return "Request Timeout"
	case StatusConflict:
		return "Conflict"
	case StatusGone:
		return "Gone"
	case StatusLengthRequired:
		return "Length Required"
	case StatusPreconditionFailed:
		return "Precondition Failed"
	case StatusRequestEntityTooLarge:
		return "Request Entity Too Large"
	case StatusRequestURITooLong:
		return "Request URI Too Long"
	case StatusUnsupportedMediaType:
		return "Unsupported Media Type"
	case StatusRequestedRangeNotSatisfiable:
		return "Requested Range Not Satisfiable"
	case StatusExpectationFailed:
		return "Expectation Failed"
	case StatusTeapot:
		return "I'm a teapot"
	case StatusMisdirectedRequest:
		return "Misdirected Request"
	case StatusUnprocessableEntity:
		return "Unprocessable Entity"
	case StatusLocked:
		return "Locked"
	case StatusFailedDependency:
		return "Failed Dependency"
	case StatusTooEarly:
		return "Too Early"
	case StatusUpgradeRequired:
		return "Upgrade Required"
	case StatusPreconditionRequired:
		return "Precondition Required"
	case StatusTooManyRequests:
		return "Too Many Requests"
	case StatusRequestHeaderFieldsTooLarge:
		return "Request Header Fields Too Large"
	case StatusUnavailableForLegalReasons:
		return "Unavailable For Legal Reasons"
	case StatusInternalServerError:
		return "Internal Server Error"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	case StatusGatewayTimeout:
		return "Gateway Timeout"
	case StatusHTTPVersionNotSupported:
		return "HTTP Version Not Supported"
	case StatusVariantAlsoNegotiates:
		return "Variant Also Negotiates"
	case StatusInsufficientStorage:
		return "Insufficient Storage"
	case StatusLoopDetected:
		return "Loop Detected"
	case StatusNotExtended:
		return "Not Extended"
	case StatusNetworkAuthenticationRequired:
		return "Network Authentication Required"
	default:
		return ""
	}
}
