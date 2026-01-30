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

// Result represents an HTTP response produced by a handler.
type Result interface {
	// StatusCode returns the HTTP status code of the response.
	StatusCode() Status
	// BodyValue returns the body value of the response.
	BodyValue() any
	// Header returns the HTTP headers of the response.
	Header() Header
}

// TypedResult represents a typed HTTP response with a generic body type.
// It implements the Result interface and provides a convenient way
// to return structured responses from handlers.
type TypedResult[T any] struct {
	// Body contains the response body that will be serialized.
	Body T
	// Status is the HTTP status code of the response.
	Status Status
	// Header contains the HTTP headers to be sent with the response.
	Headers Header
}

// StatusCode returns the HTTP status code of the response.
// If Status is not set, it defaults to StatusOK.
func (t TypedResult[T]) StatusCode() Status {
	if t.Status == 0 {
		return StatusOK
	}

	return t.Status
}

// BodyValue returns the body as an any type for serialization.
func (t TypedResult[T]) BodyValue() any {
	return t.Body
}

// Header returns the HTTP headers of the response.
func (t TypedResult[T]) Header() Header {
	return t.Headers
}
