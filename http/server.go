package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ServerProperties struct {
	Port int `property:"port" default:"8080"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration `property:"read-timeout"`

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If zero, the value of
	// ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	ReadHeaderTimeout time.Duration `property:"read-header-timeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout time.Duration `property:"write-timeout"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If zero, the value
	// of ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	IdleTimeout time.Duration `property:"idle-timeout"`

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int `property:"max-header-bytes"`

	TLS   TLSProperties   `property:"tls"`
	HTTP2 HTTP2Properties `property:"http2"`
}

type TLSProperties struct {
	Enabled  bool   `property:"enabled" default:"false"`
	CertFile string `property:"cert-file"`
	KeyFile  string `property:"key-file"`
}

type HTTP2Properties struct {
	Enabled bool `property:"enabled" default:"false"`

	// MaxConcurrentStreams optionally specifies the number of
	// concurrent streams that a peer may have open at a time.
	// If zero, MaxConcurrentStreams defaults to at least 100.
	MaxConcurrentStreams int `property:"max-concurrent-streams"`

	// MaxDecoderHeaderTableSize optionally specifies an upper limit for the
	// size of the header compression table used for decoding headers sent
	// by the peer.
	// A valid value is less than 4MiB.
	// If zero or invalid, a default value is used.
	MaxDecoderHeaderTableSize int `property:"max-decoder-header-table-size"`

	// MaxEncoderHeaderTableSize optionally specifies an upper limit for the
	// header compression table used for sending headers to the peer.
	// A valid value is less than 4MiB.
	// If zero or invalid, a default value is used.
	MaxEncoderHeaderTableSize int `property:"max-encoder-header-table-size"`

	// MaxReadFrameSize optionally specifies the largest frame
	// this endpoint is willing to read.
	// A valid value is between 16KiB and 16MiB, inclusive.
	// If zero or invalid, a default value is used.
	MaxReadFrameSize int `property:"max-read-frame-size"`

	// MaxReceiveBufferPerConnection is the maximum size of the
	// flow control window for data received on a connection.
	// A valid value is at least 64KiB and less than 4MiB.
	// If invalid, a default value is used.
	MaxReceiveBufferPerConnection int `property:"max-receive-buffer-per-connection"`

	// MaxReceiveBufferPerStream is the maximum size of
	// the flow control window for data received on a stream (request).
	// A valid value is less than 4MiB.
	// If zero or invalid, a default value is used.
	MaxReceiveBufferPerStream int `property:"max-receive-buffer-per-stream"`

	// SendPingTimeout is the timeout after which a health check using a ping
	// frame will be carried out if no frame is received on a connection.
	// If zero, no health check is performed.
	SendPingTimeout time.Duration `property:"send-ping-timeout"`

	// PingTimeout is the timeout after which a connection will be closed
	// if a response to a ping is not received.
	// If zero, a default of 15 seconds is used.
	PingTimeout time.Duration `property:"ping-timeout"`

	// WriteByteTimeout is the timeout after which a connection will be
	// closed if no data can be written to it. The timeout begins when data is
	// available to write, and is extended whenever any bytes are written.
	WriteByteTimeout time.Duration `property:"write-byte-timeout"`

	// PermitProhibitedCipherSuites, if true, permits the use of
	// cipher suites prohibited by the HTTP/2 spec.
	PermitProhibitedCipherSuites bool `property:"permit-prohibited-cipher-suites"`
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
			Handler: NewServerRequestDispatcher(d.router),
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

type ServerRequestDispatcher struct {
	contextPool sync.Pool

	endpointRegistry EndpointRegistry
}

func NewServerRequestDispatcher(endpointRegistry EndpointRegistry) *ServerRequestDispatcher {
	return &ServerRequestDispatcher{
		endpointRegistry: endpointRegistry,
	}
}

func (s *ServerRequestDispatcher) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := s.contextPool.Get().(*Context)
	ctx.reset(writer, request)

	defer func() {
		s.contextPool.Put(ctx)
	}()

	//endpoint, ok := s.endpointRegistry.Match(ctx)
}
