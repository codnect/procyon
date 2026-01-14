package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"codnect.io/procyon/http"
)

type CreateUserRequest struct {
	Body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `from:"body"`
}

type UserResponse struct {
}

type MyRouter struct {
}

func (m MyRouter) MapRoutes(r http.Router) {
	r.Get("/users", http.HandlerFunc(m.listUsers))
	r.Post("/users", http.Typed(m.createUser))
}

func (m MyRouter) listUsers(ctx *http.Context) (http.Result, error) {
	req, err := http.FromBody[CreateUserRequest](ctx)
	if err != nil {
		return nil, err
	}

	return http.BodyResult{
		Value:  []UserResponse{},
		Status: http.StatusOK,
		Header: map[string]string{
			"X-Custom-Header": "CustomValue",
		},
	}, nil

	/*
		return http.EmptyResult{
			Status: http.StatusNoContent,
			Header: map[string]string{
				"X-Custom-Header": "CustomValue",
			},
		}, nil

		return http.ViewResult{
			Name: "users/list.html",
			Model: map[string]any{
				"users": []UserResponse{},
			},
			Status: http.StatusOK,
		}, nil*/
}

func (m MyRouter) createUser(ctx *http.RequestContext[http.Void]) (http.ResponseEntity[UserResponse], error) {
	return http.EntityNotFound[UserResponse](), nil
}

func main() {
	my := MyRouter{}
	http.Typed(my.createUser)

	serverProps := http.ServerProperties{
		Port:        8080,
		ReadTimeout: time.Second * 5,
		TLS: http.TLSProperties{
			Enabled:  true,
			CertFile: "path/to/cert.pem",
			KeyFile:  "path/to/key.pem",
		},
		HTTP2: http.HTTP2Properties{
			Enabled:              true,
			MaxConcurrentStreams: 1000,
		},
	}

	defaultServer := http.NewDefaultServer(serverProps, nil)
	if err := defaultServer.Start(context.Background()); err != nil {
		panic(err)
	}
}
