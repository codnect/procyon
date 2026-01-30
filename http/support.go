// Copyright 2026 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import "net/http"

// Cookie represents an HTTP cookie.
type Cookie = http.Cookie

// SameSite represents the SameSite attribute of the cookie.
type SameSite = http.SameSite

const (
	// SameSiteDefaultMode represents the default mode of the SameSite attribute.
	SameSiteDefaultMode SameSite = iota + 1
	// SameSiteLaxMode represents the lax mode of the SameSite attribute.
	SameSiteLaxMode
	// SameSiteStrictMode represents the strict mode of the SameSite attribute.
	SameSiteStrictMode
	// SameSiteNoneMode represents the none mode of the SameSite attribute.
	SameSiteNoneMode
)

// A Header represents the key-value pairs in an HTTP header.
type Header = http.Header

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

// Status represents an HTTP status code.
type Status int

// These constants represent common HTTP status codes.
const (
	StatusContinue           Status = 100
	StatusSwitchingProtocols Status = 101
	StatusProcessing         Status = 102
	StatusEarlyHints         Status = 103

	StatusOK                   Status = 200
	StatusCreated              Status = 201
	StatusAccepted             Status = 202
	StatusNonAuthoritativeInfo Status = 203
	StatusNoContent            Status = 204
	StatusResetContent         Status = 205
	StatusPartialContent       Status = 206
	StatusMultiStatus          Status = 207
	StatusAlreadyReported      Status = 208
	StatusIMUsed               Status = 226

	StatusMultipleChoices   Status = 300
	StatusMovedPermanently  Status = 301
	StatusFound             Status = 302
	StatusSeeOther          Status = 303
	StatusNotModified       Status = 304
	StatusUseProxy          Status = 305
	_                       Status = 306
	StatusTemporaryRedirect Status = 307
	StatusPermanentRedirect Status = 308

	StatusBadRequest                   Status = 400
	StatusUnauthorized                 Status = 401
	StatusPaymentRequired              Status = 402
	StatusForbidden                    Status = 403
	StatusNotFound                     Status = 404
	StatusMethodNotAllowed             Status = 405
	StatusNotAcceptable                Status = 406
	StatusProxyAuthRequired            Status = 407
	StatusRequestTimeout               Status = 408
	StatusConflict                     Status = 409
	StatusGone                         Status = 410
	StatusLengthRequired               Status = 411
	StatusPreconditionFailed           Status = 412
	StatusRequestEntityTooLarge        Status = 413
	StatusRequestURITooLong            Status = 414
	StatusUnsupportedMediaType         Status = 415
	StatusRequestedRangeNotSatisfiable Status = 416
	StatusExpectationFailed            Status = 417
	StatusTeapot                       Status = 418
	StatusMisdirectedRequest           Status = 421
	StatusUnprocessableEntity          Status = 422
	StatusLocked                       Status = 423
	StatusFailedDependency             Status = 424
	StatusTooEarly                     Status = 425
	StatusUpgradeRequired              Status = 426
	StatusPreconditionRequired         Status = 428
	StatusTooManyRequests              Status = 429
	StatusRequestHeaderFieldsTooLarge  Status = 431
	StatusUnavailableForLegalReasons   Status = 451

	StatusInternalServerError           Status = 500
	StatusNotImplemented                Status = 501
	StatusBadGateway                    Status = 502
	StatusServiceUnavailable            Status = 503
	StatusGatewayTimeout                Status = 504
	StatusHTTPVersionNotSupported       Status = 505
	StatusVariantAlsoNegotiates         Status = 506
	StatusInsufficientStorage           Status = 507
	StatusLoopDetected                  Status = 508
	StatusNotExtended                   Status = 510
	StatusNetworkAuthenticationRequired Status = 511
)

// String method returns the text associated with the HTTP status code.
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
		return "Request-URI Too Long"
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
		return "Unknown Status"
	}
}
