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
	"context"
	"fmt"
	"net/http"
	"sync"
)

// ServerProperties defines the configuration properties for the Server component.
type ServerProperties struct {
	Port int `property:"port" default:"8080"`
}

// Server is the HTTP server that listens for incoming requests and
// dispatches them through the configured Dispatcher.
//
// It implements http.Handler and uses a sync.Pool for Context reuse
// to minimize allocations per request.
type Server struct {
	props       ServerProperties
	httpServer  *http.Server
	contextPool sync.Pool
	dispatcher  Dispatcher
}

// NewServer creates a new Server with the given properties and dispatcher.
// The dispatcher is invoked for every incoming request to route it through
// the middleware pipeline to the appropriate endpoint handler.
func NewServer(props ServerProperties, dispatcher Dispatcher) *Server {
	return &Server{
		props: props,
		contextPool: sync.Pool{
			New: func() any {
				return &Context{
					req:    &ServerRequest{},
					res:    &ServerResponse{},
					values: map[any]any{},
				}
			},
		},
		dispatcher: dispatcher,
	}
}

// Start begins listening for HTTP requests on the configured port.
// It blocks until the server is shut down or an error occurs.
func (s *Server) Start(ctx context.Context) error {
	if s.httpServer == nil {
		s.httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", s.props.Port),
			Handler: s,
		}
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Stop gracefully shuts down the server without interrupting
// any active connections.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Port returns the port number the server is configured to listen on.
func (s *Server) Port() int {
	return s.props.Port
}

// ServeHTTP handles an incoming HTTP request by obtaining a pooled
// Context, dispatching it through the middleware pipeline, and
// returning the Context to the pool when done.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := s.contextPool.Get().(*Context)
	ctx.reset(w, r)

	defer func() {
		s.contextPool.Put(ctx)
	}()

	_ = s.dispatcher.Dispatch(ctx)
}
