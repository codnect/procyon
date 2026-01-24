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

// HandlerFunc represents a function that handles an HTTP request
// and produces a Result.
type HandlerFunc func(ctx *Context) (Result, error)

// Handle processes the HTTP request contained in the Context
// and returns a Result or an error.
func (f HandlerFunc) Handle(ctx *Context) (Result, error) {
	return f(ctx)
}

// Handler represents a component that can handle an HTTP request
// and produce a Result.
type Handler interface {
	// Handle processes the HTTP request contained in the Context
	// and returns a Result or an error.
	Handle(ctx *Context) (Result, error)
}
