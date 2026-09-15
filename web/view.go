package web

import "codnect.io/procyon/http"

type ViewResult struct {
	ViewName   string
	Model      any
	StatusCode http.Status
}

func View(name string) ViewResult {
	return ViewResult{
		ViewName:   name,
		StatusCode: http.StatusOK,
	}
}

func ViewModel(name string, model any) ViewResult {
	return ViewResult{
		ViewName:   name,
		Model:      model,
		StatusCode: http.StatusOK,
	}
}

func ViewStatus(name string, status http.Status) ViewResult {
	return ViewResult{
		ViewName:   name,
		StatusCode: status,
	}
}

func (v ViewResult) Status() http.Status {
	return v.StatusCode
}

func (v ViewResult) Header() http.Header {
	return nil
}

type ViewResolver interface {
}
