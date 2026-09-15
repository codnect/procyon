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
	"fmt"
	"net"
	"net/http"
	"sync"

	"codnect.io/logy"
)

// ServerProperties defines the configuration properties for the Server component.
type ServerProperties struct {
	Port int `property:"port,default=8080"`
}

func newServerProperties() *ServerProperties {
	return &ServerProperties{}
}

func (s *ServerProperties) Prefix() string {
	return "server"
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
	mu          sync.RWMutex
	running     bool
	stopping    bool
	boundPort   int
}

// NewServer creates a new Server with the given properties and dispatcher.
// The dispatcher is invoked for every incoming request to route it through
// the middleware pipeline to the appropriate endpoint handler.
func NewServer(props ServerProperties, dispatcher Dispatcher) *Server {
	if dispatcher == nil {
		panic("nil dispatcher")
	}

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

// newServer adapts the pointer properties component to the public constructor.
func newServer(props *ServerProperties, dispatcher Dispatcher) *Server {
	return NewServer(*props, dispatcher)
}

// Start binds synchronously, then serves in the background. A bind failure is
// returned before lifecycle startup succeeds.
func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopping {
		return fmt.Errorf("HTTP server is stopping")
	}
	if s.running {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", fmt.Sprintf(":%d", s.props.Port))
	if err != nil {
		return err
	}
	server := &http.Server{Handler: s}
	s.httpServer = server
	s.boundPort = listener.Addr().(*net.TCPAddr).Port
	s.running = true
	go func() {
		err := server.Serve(listener)
		s.mu.Lock()
		if s.httpServer == server {
			s.running = false
		}
		s.mu.Unlock()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logy.Get().Error("HTTP server stopped unexpectedly", err)
		}
	}()
	logy.Get().Info("HTTP server started on port {}", s.boundPort)
	return nil
}

// Stop drains active requests; a timeout also closes remaining connections.
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	server := s.httpServer
	if server == nil {
		s.mu.Unlock()
		return nil
	}
	s.stopping = true
	s.mu.Unlock()
	err := server.Shutdown(ctx)
	if err != nil {
		err = errors.Join(err, server.Close())
	}
	s.mu.Lock()
	s.running = false
	s.stopping = false
	s.mu.Unlock()
	return err
}

func (s *Server) Port() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.boundPort != 0 {
		return s.boundPort
	}
	return s.props.Port
}

// ServeHTTP handles an incoming HTTP request by obtaining a pooled
// Context, dispatching it through the middleware pipeline, and
// returning the Context to the pool when done.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := s.contextPool.Get().(*Context)
	ctx.reset(r, w)

	defer func() {
		s.contextPool.Put(ctx)
	}()

	if err := s.dispatcher.Dispatch(ctx); err != nil {
		logy.Get().Error("HTTP request failed", err)
		if !ctx.Response().IsCommitted() {
			_ = ctx.Response().Reset()
			ctx.Response().SetStatus(StatusInternalServerError)
		}
	} else if ctx.Endpoint() == nil && !ctx.Response().IsCommitted() && ctx.Response().Status() == StatusOK {
		ctx.Response().SetStatus(StatusNotFound)
	}
	ctx.Response().writeHeaders()
}

// serverLifecycle keeps lifecycle discovery separate from runtime.Server.
type serverLifecycle struct{ server *Server }

func newServerLifecycle(server *Server) *serverLifecycle {
	return &serverLifecycle{server: server}
}

func (s *serverLifecycle) Start(ctx context.Context) error { return s.server.Start(ctx) }
func (s *serverLifecycle) Stop(ctx context.Context) error  { return s.server.Stop(ctx) }
func (s *serverLifecycle) IsRunning() bool {
	s.server.mu.RLock()
	defer s.server.mu.RUnlock()
	return s.server.running
}
