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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_String(t *testing.T) {
	testCases := []struct {
		name string
		code Status
		want string
	}{
		{name: "100 Continue", code: StatusContinue, want: "Continue"},
		{name: "101 Switching Protocols", code: StatusSwitchingProtocols, want: "Switching Protocols"},
		{name: "102 Processing", code: StatusProcessing, want: "Processing"},
		{name: "103 Early Hints", code: StatusEarlyHints, want: "Early Hints"},
		{name: "200 OK", code: StatusOK, want: "OK"},
		{name: "201 Created", code: StatusCreated, want: "Created"},
		{name: "202 Accepted", code: StatusAccepted, want: "Accepted"},
		{name: "203 Non-Authoritative Information", code: StatusNonAuthoritativeInfo, want: "Non-Authoritative Information"},
		{name: "204 No Content", code: StatusNoContent, want: "No Content"},
		{name: "205 Reset Content", code: StatusResetContent, want: "Reset Content"},
		{name: "206 Partial Content", code: StatusPartialContent, want: "Partial Content"},
		{name: "207 Multi-Status", code: StatusMultiStatus, want: "Multi-Status"},
		{name: "208 Already Reported", code: StatusAlreadyReported, want: "Already Reported"},
		{name: "226 IM Used", code: StatusIMUsed, want: "IM Used"},
		{name: "300 Multiple Choices", code: StatusMultipleChoices, want: "Multiple Choices"},
		{name: "301 Moved Permanently", code: StatusMovedPermanently, want: "Moved Permanently"},
		{name: "302 Found", code: StatusFound, want: "Found"},
		{name: "303 See Other", code: StatusSeeOther, want: "See Other"},
		{name: "304 Not Modified", code: StatusNotModified, want: "Not Modified"},
		{name: "305 Use Proxy", code: StatusUseProxy, want: "Use Proxy"},
		{name: "307 Temporary Redirect", code: StatusTemporaryRedirect, want: "Temporary Redirect"},
		{name: "308 Permanent Redirect", code: StatusPermanentRedirect, want: "Permanent Redirect"},
		{name: "400 Bad Request", code: StatusBadRequest, want: "Bad Request"},
		{name: "401 Unauthorized", code: StatusUnauthorized, want: "Unauthorized"},
		{name: "402 Payment Required", code: StatusPaymentRequired, want: "Payment Required"},
		{name: "403 Forbidden", code: StatusForbidden, want: "Forbidden"},
		{name: "404 Not Found", code: StatusNotFound, want: "Not Found"},
		{name: "405 Method Not Allowed", code: StatusMethodNotAllowed, want: "Method Not Allowed"},
		{name: "406 Not Acceptable", code: StatusNotAcceptable, want: "Not Acceptable"},
		{name: "407 Proxy Authentication Required", code: StatusProxyAuthRequired, want: "Proxy Authentication Required"},
		{name: "408 Request Timeout", code: StatusRequestTimeout, want: "Request Timeout"},
		{name: "409 Conflict", code: StatusConflict, want: "Conflict"},
		{name: "410 Gone", code: StatusGone, want: "Gone"},
		{name: "411 Length Required", code: StatusLengthRequired, want: "Length Required"},
		{name: "412 Precondition Failed", code: StatusPreconditionFailed, want: "Precondition Failed"},
		{name: "413 Request Entity Too Large", code: StatusRequestEntityTooLarge, want: "Request Entity Too Large"},
		{name: "414 Request-URI Too Long", code: StatusRequestURITooLong, want: "Request-URI Too Long"},
		{name: "415 Unsupported Media Type", code: StatusUnsupportedMediaType, want: "Unsupported Media Type"},
		{name: "416 Requested Range Not Satisfiable", code: StatusRequestedRangeNotSatisfiable, want: "Requested Range Not Satisfiable"},
		{name: "417 Expectation Failed", code: StatusExpectationFailed, want: "Expectation Failed"},
		{name: "418 I'm a teapot", code: StatusTeapot, want: "I'm a teapot"},
		{name: "421 Misdirected Request", code: StatusMisdirectedRequest, want: "Misdirected Request"},
		{name: "422 Unprocessable Entity", code: StatusUnprocessableEntity, want: "Unprocessable Entity"},
		{name: "423 Locked", code: StatusLocked, want: "Locked"},
		{name: "424 Failed Dependency", code: StatusFailedDependency, want: "Failed Dependency"},
		{name: "425 Too Early", code: StatusTooEarly, want: "Too Early"},
		{name: "426 Upgrade Required", code: StatusUpgradeRequired, want: "Upgrade Required"},
		{name: "428 Precondition Required", code: StatusPreconditionRequired, want: "Precondition Required"},
		{name: "429 Too Many Requests", code: StatusTooManyRequests, want: "Too Many Requests"},
		{name: "431 Request Header Fields Too Large", code: StatusRequestHeaderFieldsTooLarge, want: "Request Header Fields Too Large"},
		{name: "451 Unavailable For Legal Reasons", code: StatusUnavailableForLegalReasons, want: "Unavailable For Legal Reasons"},
		{name: "500 Internal Server Error", code: StatusInternalServerError, want: "Internal Server Error"},
		{name: "501 Not Implemented", code: StatusNotImplemented, want: "Not Implemented"},
		{name: "502 Bad Gateway", code: StatusBadGateway, want: "Bad Gateway"},
		{name: "503 Service Unavailable", code: StatusServiceUnavailable, want: "Service Unavailable"},
		{name: "504 Gateway Timeout", code: StatusGatewayTimeout, want: "Gateway Timeout"},
		{name: "505 HTTP Version Not Supported", code: StatusHTTPVersionNotSupported, want: "HTTP Version Not Supported"},
		{name: "506 Variant Also Negotiates", code: StatusVariantAlsoNegotiates, want: "Variant Also Negotiates"},
		{name: "507 Insufficient Storage", code: StatusInsufficientStorage, want: "Insufficient Storage"},
		{name: "508 Loop Detected", code: StatusLoopDetected, want: "Loop Detected"},
		{name: "510 Not Extended", code: StatusNotExtended, want: "Not Extended"},
		{name: "511 Network Authentication Required", code: StatusNetworkAuthenticationRequired, want: "Network Authentication Required"},
		{name: "unknown status code", code: Status(999), want: "Unknown Status"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			text := tc.code.String()

			// then
			assert.Equal(t, tc.want, text)
		})
	}
}
