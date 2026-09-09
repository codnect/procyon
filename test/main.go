package main

import (
	"os"

	"codnect.io/procyon"
	"codnect.io/procyon/http"
	"codnect.io/procyon/web"
)

type WelcomeController struct {
}

func NewWelcomeController() *WelcomeController {
	return &WelcomeController{}
}

func (w *WelcomeController) ConfigureEndpoints(endpoints http.Endpoints) {
	endpoints.MapGet("/hello", http.HandleResult(w.hello))
}

func (w *WelcomeController) hello(ctx *http.Context) (web.ViewResult, error) {
	return web.View("hello.html"), nil
}

func main() {
	if err := procyon.New().Run(); err != nil {
		os.Exit(1)
	}
}
