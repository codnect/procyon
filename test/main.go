package main

import (
	"os"

	"codnect.io/procyon"
	"codnect.io/procyon/component"
	"codnect.io/procyon/http"
)

type UserResponse struct {
	Name string `json:"name"`
}

type WelcomeController struct {
}

func NewWelcomeController() *WelcomeController {
	return &WelcomeController{}
}

func (w *WelcomeController) ConfigureEndpoints(endpoints http.Endpoints) {
	endpoints.MapGet("/hello", http.HandleResult(w.hello))
}

func (w *WelcomeController) hello(ctx *http.Context) (http.Result, error) {
	return http.Json(UserResponse{
		Name: "John Doe",
	}, http.StatusOK), nil
}

func main() {
	component.Register(NewWelcomeController)
	http.ServerProperties{}
	if err := procyon.New().Run(); err != nil {
		os.Exit(1)
	}
}
