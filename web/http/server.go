package http

import (
	"codnect.io/procyon/web"
	"context"
	"fmt"
	"net/http"
	"sync"
)

// DefaultServer is the default implementation of the Server interface.
type DefaultServer struct {
	httpServer *http.Server
	props      *web.ServerProperties

	routeRegistry *RouteRegistry
	contextPool   sync.Pool
}

// NewDefaultServer creates a new DefaultServer with the provided properties and route registry.
func NewDefaultServer(props *web.ServerProperties, routeRegistry *RouteRegistry) *DefaultServer {
	if props == nil {
		panic("nil web server props")
	}

	if routeRegistry == nil {
		panic("nil route registry")
	}

	return &DefaultServer{
		props:         props,
		routeRegistry: routeRegistry,
	}
}

// Start starts the server.
func (s *DefaultServer) Start(ctx context.Context) error {
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

// Stop stops the server.
func (s *DefaultServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Port returns the port of the server.
func (s *DefaultServer) Port() int {
	return s.props.Port
}

// ServeHTTP serves the HTTP request.
func (s *DefaultServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := s.contextPool.Get().(*defaultServerContext)
	ctx.reset(writer, request)

	defer func() {
		s.contextPool.Put(ctx)
	}()

	router, ok := s.routeRegistry.Find(ctx)
	if ok {
		ctx.handlerChain = router.HandlerChain()
		ctx.Invoke(ctx)

		if ctx.Err() != nil {
		}
	} else {
	}
}
