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
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"

	"codnect.io/logy"
	"go.codnect.io/procyon/runtime"
)

// Server represents an HTTP server managed by the application runtime.
type Server interface {
	runtime.Server

	// Port returns the port the server is bound to, or the configured port
	// if the server has not been started.
	Port() int
}

// ServerProperties defines the configuration properties for an HTTP server.
type ServerProperties struct {
	// Port specifies the port the HTTP server listens on.
	Port int `property:"port,default=8080"`
}

// newServerProperties creates a new ServerProperties.
func newServerProperties() *ServerProperties {
	return &ServerProperties{}
}

// Prefix returns the configuration property prefix for the HTTP server.
func (s *ServerProperties) Prefix() string {
	return "server"
}

// serverAdapter adapts net/http requests to the HTTP request handling pipeline.
type serverAdapter struct {
	contextPool sync.Pool
	dispatcher  Dispatcher
}

// newServerAdapter creates a new serverAdapter with the given dispatcher.
func newServerAdapter(dispatcher Dispatcher) *serverAdapter {
	if dispatcher == nil {
		panic("nil dispatcher")
	}

	return &serverAdapter{
		contextPool: sync.Pool{
			New: func() any {
				return newContext(nil, nil)
			},
		},
		dispatcher: dispatcher,
	}
}

// ServeHTTP dispatches an incoming HTTP request through the request handling
// pipeline.
func (a *serverAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := a.contextPool.Get().(*Context)
	ctx.reset(r, w)

	defer a.contextPool.Put(ctx)

	if err := a.dispatcher.Dispatch(ctx); err != nil {
		logy.Get().Error("HTTP request dispatch failed: {}", err)
	}
}

// DefaultServer is the default HTTP server implementation.
type DefaultServer struct {
	props         *ServerProperties
	httpServer    *http.Server
	serverAdapter *serverAdapter
	boundPort     int
}

// NewDefaultServer creates a new DefaultServer with the given properties and
// dispatcher.
func NewDefaultServer(props *ServerProperties, dispatcher Dispatcher) *DefaultServer {
	if props == nil {
		panic("nil server properties")
	}

	if dispatcher == nil {
		panic("nil dispatcher")
	}

	return &DefaultServer{
		props:         props,
		serverAdapter: newServerAdapter(dispatcher),
	}
}

// Start begins listening for HTTP requests on the configured port.
func (s *DefaultServer) Start(ctx context.Context) error {
	addr := net.JoinHostPort("", strconv.Itoa(s.props.Port))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{
		Handler: s.serverAdapter,
	}

	s.boundPort = listener.Addr().(*net.TCPAddr).Port

	go s.serve(listener)

	log.Info("HTTP server started on port {}", s.boundPort)
	return nil
}

// Shutdown gracefully shuts down the server without interrupting active
// connections.
func (s *DefaultServer) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Info("HTTP server stopped")
	return nil

}

// Port returns the bound port, or the configured port before the server has
// been started.
func (s *DefaultServer) Port() int {
	if s.boundPort != 0 {
		return s.boundPort
	}

	return s.props.Port
}

// serve serves HTTP requests using the given listener.
func (s *DefaultServer) serve(listener net.Listener) {
	if err := s.httpServer.Serve(listener); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Error("HTTP server stopped unexpectedly: {}", err)
	}
}
