package web

import (
	"codnect.io/procyon/web/http"
	"codnect.io/procyon/web/http/router"
	"context"
	"errors"
	"fmt"
	stdhttp "net/http"
	"sync"
)

type Server interface {
	Start() error
	Stop() error
	Port() int
	ShutDownGracefully(ctx context.Context) error
}

type ServerProperties struct {
	//property.Properties `prefix:"server"`
	Port int `prop:"port" default:"8080"`
}

type DefaultServer struct {
	props  ServerProperties
	server *stdhttp.Server

	contextPool sync.Pool

	mappingRegistry *router.MappingRegistry
	errorHandler    http.ErrorHandler
}

func NewDefaultServer() *DefaultServer {
	return &DefaultServer{
		contextPool: sync.Pool{
			New: func() any {
				return newServerContext()
			},
		},
	}
}

func (s *DefaultServer) Start() error {
	s.server = &stdhttp.Server{
		Addr:    fmt.Sprintf(":%d", s.props.Port),
		Handler: s,
	}

	return nil
}

func (s *DefaultServer) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *DefaultServer) Port() int {
	return s.props.Port
}

func (s *DefaultServer) ShutDownGracefully(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s *DefaultServer) ServeHTTP(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
	ctx := s.contextPool.Get().(*ServerContext)
	ctx.Reset(request, writer)

	defer func() {
		if runtimeError := recover(); runtimeError != nil {
			s.handlerRuntimeError(ctx, runtimeError)
			s.contextPool.Put(ctx)
		}
	}()

	chain, exists := s.mappingRegistry.GetHandler(ctx)

	if exists {
		ctx.HandlerChain = chain
		ctx.Invoke(ctx)

		if ctx.Err() != nil {
			s.errorHandler.HandleError(ctx, ctx.Err())
		}
	} else {
		s.errorHandler.HandleError(ctx, &http.NotFoundError{})
	}

	s.contextPool.Put(ctx)
}

func (s *DefaultServer) handlerRuntimeError(ctx http.Context, runtimeError any) {
	switch err := runtimeError.(type) {
	case error:
		s.errorHandler.HandleError(ctx, err)
	case string:
		s.errorHandler.HandleError(ctx, errors.New(err))
	default:
		s.errorHandler.HandleError(ctx, fmt.Errorf("unknown error: %v", err))
	}
}
