package httpx

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type ServerProperties struct {
	Port int `property:"port" default:"8080"`
}

// Server interface represents a server that can be started and stopped.
type Server interface {
	// Start method starts the server.
	Start(ctx context.Context) error
	// Stop method stops the server.
	Stop(ctx context.Context) error
	// Port method returns the port the server is running on.
	Port() int
}

type DefaultServer struct {
	props  ServerProperties
	router EndpointRegistry

	httpServer *http.Server
}

func NewDefaultServer(props ServerProperties, router EndpointRegistry) *DefaultServer {
	return &DefaultServer{
		props:  props,
		router: router,
	}
}

func (d *DefaultServer) Start(ctx context.Context) error {
	if d.httpServer == nil {
		d.httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", d.props.Port),
			Handler: NewRequestDispatcher(d.router),
		}
	}

	if err := d.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (d *DefaultServer) Stop(ctx context.Context) error {
	return d.httpServer.Shutdown(ctx)
}

func (d *DefaultServer) Port() int {
	return d.props.Port
}

type RequestDispatcher struct {
	contextPool sync.Pool

	endpointRegistry EndpointRegistry
}

func NewRequestDispatcher(endpointRegistry EndpointRegistry) *RequestDispatcher {
	return &RequestDispatcher{
		endpointRegistry: endpointRegistry,
	}
}

func (s *RequestDispatcher) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := s.contextPool.Get().(*defaultContext)
	ctx.reset(writer, request)

	defer func() {
		s.contextPool.Put(ctx)
	}()

	endpoint, ok := s.endpointRegistry.Match(ctx)

}
