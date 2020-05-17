package server

type WebServerFactory interface {
	GetWebServer() *WebServer
}

type DefaultWebServerFactory struct {
}

func NewDefaultWebServerFactory() WebServerFactory {
	return &DefaultWebServerFactory{}
}

func (factory DefaultWebServerFactory) GetWebServer() *WebServer {
	return nil
}
