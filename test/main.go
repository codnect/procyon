package main

import (
	"os"

	"codnect.io/procyon"
	"codnect.io/procyon/component"
	"codnect.io/procyon/http"
)

type Test interface {
	WriteS(s any)
}

type HelloController struct {
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (h *HelloController) WriteS[T any](s T) {

}

func (h *HelloController) ConfigureEndpoints(endpoints http.Endpoints) {
	endpoints.MapGet("/hello", http.Handle(h.sayHello))
}

func (h *HelloController) sayHello(ctx *http.EndpointContext[string]) error {
	response := ctx.Response()
	response.Writer().Write([]byte("Hello, World!"))
	return nil
}

func init() {
	component.Register(NewHelloController)
}

func main() {
	if err := procyon.Run(); err != nil {
		os.Exit(1)
	}
}
