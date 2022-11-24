package web

type HttpStatus int

const (
	StatusContinue           HttpStatus = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols HttpStatus = 101 // RFC 7231, 6.2.2
	StatusProcessing         HttpStatus = 102 // RFC 2518, 10.1
	StatusEarlyHints         HttpStatus = 103 // RFC 8297

	StatusOK                   HttpStatus = 200 // RFC 7231, 6.3.1
	StatusCreated              HttpStatus = 201 // RFC 7231, 6.3.2
	StatusAccepted             HttpStatus = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo HttpStatus = 203 // RFC 7231, 6.3.4
	StatusNoContent            HttpStatus = 204 // RFC 7231, 6.3.5
	StatusResetContent         HttpStatus = 205 // RFC 7231, 6.3.6
	StatusPartialContent       HttpStatus = 206 // RFC 7233, 4.1
	StatusMultiStatus          HttpStatus = 207 // RFC 4918, 11.1
	StatusAlreadyReported      HttpStatus = 208 // RFC 5842, 7.1
	StatusIMUsed               HttpStatus = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   HttpStatus = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  HttpStatus = 301 // RFC 7231, 6.4.2
	StatusFound             HttpStatus = 302 // RFC 7231, 6.4.3
	StatusSeeOther          HttpStatus = 303 // RFC 7231, 6.4.4
	StatusNotModified       HttpStatus = 304 // RFC 7232, 4.1
	StatusUseProxy          HttpStatus = 305 // RFC 7231, 6.4.5
	_                       HttpStatus = 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect HttpStatus = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect HttpStatus = 308 // RFC 7538, 3

	StatusBadRequest                   HttpStatus = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 HttpStatus = 401 // RFC 7235, 3.1
	StatusPaymentRequired              HttpStatus = 402 // RFC 7231, 6.5.2
	StatusForbidden                    HttpStatus = 403 // RFC 7231, 6.5.3
	StatusNotFound                     HttpStatus = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             HttpStatus = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                HttpStatus = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            HttpStatus = 407 // RFC 7235, 3.2
	StatusRequestTimeout               HttpStatus = 408 // RFC 7231, 6.5.7
	StatusConflict                     HttpStatus = 409 // RFC 7231, 6.5.8
	StatusGone                         HttpStatus = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               HttpStatus = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           HttpStatus = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        HttpStatus = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            HttpStatus = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         HttpStatus = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable HttpStatus = 416 // RFC 7233, 4.4
	StatusExpectationFailed            HttpStatus = 417 // RFC 7231, 6.5.14
	StatusTeapot                       HttpStatus = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           HttpStatus = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          HttpStatus = 422 // RFC 4918, 11.2
	StatusLocked                       HttpStatus = 423 // RFC 4918, 11.3
	StatusFailedDependency             HttpStatus = 424 // RFC 4918, 11.4
	StatusTooEarly                     HttpStatus = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              HttpStatus = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         HttpStatus = 428 // RFC 6585, 3
	StatusTooManyRequests              HttpStatus = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  HttpStatus = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   HttpStatus = 451 // RFC 7725, 3

	StatusInternalServerError           HttpStatus = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                HttpStatus = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    HttpStatus = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            HttpStatus = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                HttpStatus = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       HttpStatus = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         HttpStatus = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           HttpStatus = 507 // RFC 4918, 11.5
	StatusLoopDetected                  HttpStatus = 508 // RFC 5842, 7.2
	StatusNotExtended                   HttpStatus = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired HttpStatus = 511 // RFC 6585, 6
)
