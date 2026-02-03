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
