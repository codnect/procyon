package server

type WebServerException struct {
	message string
}

func NewWebServerException(errMessage string) WebServerException {
	return WebServerException{
		message: errMessage,
	}
}

func (e WebServerException) Error() string {
	return e.message
}

type WebServer interface {
	Start() WebServerException
	Stop() WebServerException
	Port() int
}

type DefaultWebServer struct {
}

func (server DefaultWebServer) Start() WebServerException {
	return NewWebServerException("")
}

func (server DefaultWebServer) Stop() WebServerException {
	return NewWebServerException("")
}

func (server DefaultWebServer) Port() int {
	return 8080
}
