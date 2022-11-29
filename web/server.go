package web

import (
	"fmt"
	"github.com/procyon-projects/procyon/env/property"
	"net/http"
)

type ServerProperties struct {
	property.Properties `prefix:"server"`

	Port int `prop:"port" default:"8080"`
}

type Server interface {
	Start() error
	Stop() error
	Port() int
	ShutDownGracefully()
}

type DefaultServer struct {
	props  *ServerProperties
	router *Router
}

func NewDefaultServer(props *ServerProperties) *DefaultServer {
	return &DefaultServer{
		props:  props,
		router: newRouter(),
	}
}

func (s *DefaultServer) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.props.Port), s)
}

func (s *DefaultServer) Stop() error {
	return nil
}

func (s *DefaultServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	s.router.Route(response, request)
}
