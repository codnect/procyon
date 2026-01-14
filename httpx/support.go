// Copyright 2025 Codnect
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

package httpx

import "net/http"

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
		return "reset Content"
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

// These constants represent common HTTP header fields.
const (
	// HeaderAccept specifies the media types that are acceptable for the response.
	HeaderAccept = "Accept"
	// HeaderCharset specifies the character sets that are acceptable.
	HeaderCharset = "Accept-Charset"
	// HeaderAcceptEncoding specifies the content codings that are acceptable in the response.
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderAcceptLanguage specifies the natural languages that are preferred in the response.
	HeaderAcceptLanguage = "Accept-Language"
	// HeaderAcceptPatch specifies the patch document formats that are acceptable.
	HeaderAcceptPatch = "Accept-Patch"
	// HeaderAcceptRanges allows the server to indicate its acceptance of range requests for a resource.
	HeaderAcceptRanges = "Accept-Ranges"

	// HeaderAccessControlAllowCredentials indicates whether the response to the defaultRequest can be exposed when the credentials flag is true.
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	// HeaderAccessControlAllowHeaders specifies the headers that are allowed in the actual defaultRequest.
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// HeaderAccessControlAllowMethods specifies the methods that are allowed when accessing the resource.
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	// HeaderAccessControlAllowOrigin specifies the origin that is allowed to access the resource.
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// HeaderAccessControlExposeHeaders specifies the headers that are exposed to the client.
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	// HeaderAccessControlMaxAge specifies the maximum amount of time that the results of a preflight defaultRequest can be cached.
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"
	// HeaderAccessControlRequestHeaders is used when issuing a preflight defaultRequest to let the server know what HTTP headers will be used in the actual defaultRequest.
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"
	// HeaderAccessControlRequestMethod is used when issuing a preflight defaultRequest to let the server know what HTTP method will be used in the actual defaultRequest.
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"

	// HeaderAge indicates the age of the response.
	HeaderAge = "Age"
	// HeaderAllow lists the set of methods supported by the resource.
	HeaderAllow = "Allow"
	// HeaderAuthorization contains the credentials to authenticate a user agent with a server.
	HeaderAuthorization = "Authorization"
	// HeaderCacheControl is used to specify directives for caching mechanisms.
	HeaderCacheControl = "Cache-Control"
	// HeaderConnection controls whether the network connection stays open after the current transaction finishes.
	HeaderConnection = "Connection"
	// HeaderContentDisposition is an extension header used in HTTP and MIME email to specify certain parameters related to the disposition of the message content.
	HeaderContentDisposition = "Content-Disposition"
	// HeaderContentEncoding is used to specify the content encodings applied to the entity-body.
	HeaderContentEncoding = "Content-Encoding"
	// HeaderContentLanguage describes the language(s) intended for the audience.
	HeaderContentLanguage = "Content-Language"
	// HeaderContentLength indicates the size of the entity-body in bytes.
	HeaderContentLength = "Content-Length"
	// HeaderContentType indicates the media type of the entity-body.
	HeaderContentType = "Content-Type"
	// HeaderCookie contains stored HTTP cookies previously sent by the server with the Set-Cookie header.
	HeaderCookie = "Cookie"
	// HeaderDate represents the date and time at which the message was originated.
	HeaderDate = "Date"
	// HeaderETag provides the current value of the entity tag for the requested variant.
	HeaderETag = "ETag"
	// HeaderExpect is used to indicate that particular server behaviors are required by the client.
	HeaderExpect = "Expect"
	// HeaderExpires gives the date/time after which the response is considered stale.
	HeaderExpires = "Expires"
	// HeaderFrom is an email address of the human user who controls the requesting user agent.
	HeaderFrom = "From"
	// HeaderHost specifies the domain name of the server and optionally the TCP port number.
	HeaderHost = "Host"

	// HeaderIfMatch is used to make a defaultRequest method conditional.
	HeaderIfMatch = "If-Match"
	// HeaderIfModifiedSince is used to make a GET or HEAD defaultRequest method conditional.
	HeaderIfModifiedSince = "If-Modified-Since"
	// HeaderIfNoneMatch is used to make a defaultRequest method conditional.
	HeaderIfNoneMatch = "If-None-Match"
	// HeaderIfRange is used to make a partial GET defaultRequest conditional.
	HeaderIfRange = "If-Range"
	// HeaderIfUnmodifiedSince is used to make a defaultRequest method conditional.
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	// HeaderLastModified indicates the date and time at which the server believes the variant was last modified.
	HeaderLastModified = "Last-Modified"
	// HeaderLink indicates that the response is part of a series of responses.
	HeaderLink = "Link"
	// HeaderLocation is used in redirection, or when a new resource has been created.
	HeaderLocation = "Location"
	// HeaderMaxForwards limits the number of times that the message can be forwarded through proxies or gateways.
	HeaderMaxForwards = "Max-Forwards"
	// HeaderOrigin indicates where a fetch originates from.
	HeaderOrigin = "Origin"
	// HeaderPragma allows backwards compatibility with HTTP/1.0 caches where the Cache-Control header is not yet present.
	HeaderPragma = "Pragma"
	// HeaderProxyAuthenticate must be included as part of a 407 Proxy Authentication Required response.
	HeaderProxyAuthenticate = "Proxy-Authenticate"
	// HeaderProxyAuthorization allows the client to identify itself (or its user) to a proxy which requires authentication.
	HeaderProxyAuthorization = "Proxy-Authorization"
	// HeaderRange is used in an HTTP defaultRequest to defaultRequest only part of a document.
	HeaderRange = "Range"
	// HeaderReferer allows the client to specify, for the server's benefit, the address of the document (or element within the document) from which the URI in the defaultRequest was obtained.
	HeaderReferer = "Referer"
	// HeaderRetryAfter indicates how long the user agent should wait before making a follow-up defaultRequest.
	HeaderRetryAfter = "Retry-After"
	// HeaderServer contains information about the software used by the origin server to handle the defaultRequest.
	HeaderServer = "Server"
	// HeaderSetCookie is sent by the server to the user agent with an HTTP response.
	HeaderSetCookie = "Set-Cookie"
	// HeaderSetCookie2 is the updated version of Set-Cookie header.
	HeaderSetCookie2 = "Set-Cookie2"
	// HeaderTE specifies the transfer encodings the user agent is willing to accept.
	HeaderTE = "TE"
	// HeaderTrailer allows the sender to include additional fields at the end of chunked messages.
	HeaderTrailer = "Trailer"
	// HeaderTransferEncoding indicates what (if any) type of transformation has been applied to the message body.
	HeaderTransferEncoding = "TransferEncoding"
	// HeaderUpgrade allows the client to specify what additional communication protocols it supports and would like to use if the server finds it appropriate to switch protocols.
	HeaderUpgrade = "Upgrade"
	// HeaderUserAgent contains information about the user agent originating the defaultRequest.
	HeaderUserAgent = "UserAgent"
	// HeaderVary determines how to match future defaultRequest headers to decide whether a cached response can be used rather than requesting a fresh one from the origin server.
	HeaderVary = "Vary"
	// HeaderVia is used by gateways and proxies to indicate the intermediate protocols and recipients between the user agent and the server on requests, and between the origin server and the client on responses.
	HeaderVia = "Via"
	// HeaderWarning is used to carry additional information about the status or transformation of a message which might not be reflected in the message.
	HeaderWarning = "Warning"
	// HeaderWWWAuthenticate must be included in 401 Unauthorized responses.
	HeaderWWWAuthenticate = "WWW-Authenticate"
)
