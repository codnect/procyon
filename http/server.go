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

package http

import (
	stdhttp "net/http"
	"strings"
)

// Server is the default HTTP server implementation using the standard net/http stack.
type Server struct {
	router     *Router
	middleware []Middleware
}

// NewServer creates a Server with a router ready to register handlers.
func NewServer() *Server {
	return &Server{
		router: NewRouter(),
	}
}

// Router exposes the underlying router to allow advanced configuration.
func (s *Server) Router() *Router {
	return s.router
}

// Use registers middleware executed for every request.
func (s *Server) Use(middleware ...Middleware) {
	s.middleware = append(s.middleware, middleware...)
}

// Handle registers a handler for the given method and pattern.
func (s *Server) Handle(method, pattern string, handler Handler, middleware ...Middleware) error {
	return s.router.Handle(method, pattern, handler, middleware...)
}

// ServeHTTP satisfies stdlib's http.Handler and dispatches requests through the router.
func (s *Server) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	method := strings.ToUpper(r.Method)
	path := r.URL.Path

	handler, middleware, params := s.router.Resolve(method, path)
	request := &standardRequest{req: r}
	response := &standardResponse{writer: w}

	security := NewUnauthenticatedSecurityContext()
	ctx := newStandardContext(request, response, security, params, append(s.middleware, middleware...), handler)
	ctx.Next()
}

// ListenAndServe starts the HTTP server on the given address.
func (s *Server) ListenAndServe(addr string) error {
	return stdhttp.ListenAndServe(addr, s)
}

// ListenAndServeTLS starts the HTTP server with TLS configuration.
func (s *Server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	return stdhttp.ListenAndServeTLS(addr, certFile, keyFile, s)
}
