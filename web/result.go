package web

import (
	"codnect.io/procyon/http"
)

type ViewResultExecutor struct {
}

func (v *ViewResultExecutor) CanExecute(result http.Result) bool {
	_, ok := result.(ViewResult)
	return ok
}

func (v *ViewResultExecutor) Execute(ctx *http.Context, result http.Result) error {
	//TODO implement me
	panic("implement me")
}
