package http

import (
	"io"
)

type Result interface {
	StatusCode() Status
	Body() any
	Headers() map[string]string
}

type EmptyResult struct {
	Status Status
	Header map[string]string
}

func OkEmpty() Result {
	return EmptyResult{
		Status: StatusOK,
	}
}

func CreatedEmpty(location string) Result {
	return EmptyResult{
		Status: StatusCreated,
		Header: map[string]string{
			HeaderLocation: location,
		},
	}
}

func NoContent() Result {
	return EmptyResult{
		Status: StatusNoContent,
	}
}

func NotFound() Result {
	return EmptyResult{
		Status: StatusNotFound,
	}
}

func Unauthorized() Result {
	return EmptyResult{
		Status: StatusUnauthorized,
	}
}

func Forbidden() Result {
	return EmptyResult{
		Status: StatusForbidden,
	}
}

func StatusCode(status Status) Result {
	return EmptyResult{
		Status: status,
	}
}

func (s EmptyResult) StatusCode() Status {
	return s.Status
}

func (s EmptyResult) Body() any {
	return nil
}

func (s EmptyResult) Headers() map[string]string {
	return s.Header
}

type BodyResult struct {
	Status Status
	Value  any
	Header map[string]string
}

func Ok(val any) Result {
	return BodyResult{
		Status: StatusOK,
		Value:  val,
	}
}

func Created(val any) Result {
	return BodyResult{
		Status: StatusCreated,
		Value:  val,
	}
}

func CreatedAt(location string, val any) Result {
	return BodyResult{
		Status: StatusCreated,
		Value:  val,
		Header: map[string]string{
			HeaderLocation: location,
		},
	}
}

func (b BodyResult) StatusCode() Status {
	return b.Status
}

func (b BodyResult) Body() any {
	return b.Value
}

func (b BodyResult) Headers() map[string]string {
	return b.Header
}

type ProblemResult struct {
	Title  string
	Detail string
	Status Status
}

func Problem() Result {
	return ProblemResult{
		Status: StatusInternalServerError,
	}
}

func ProblemStatus(status Status) Result {
	return ProblemResult{
		Status: status,
	}
}

func (p ProblemResult) StatusCode() Status {
	return p.Status
}

func (p ProblemResult) Body() any {
	return nil
}

func (p ProblemResult) Headers() map[string]string {
	return nil
}

func Text(content, contentType string) Result {
	return nil
}

type FileResult struct {
	Content      []byte
	ContentType  string
	DownloadName string
	Path         string
}

func File(path string) Result {
	return nil
}

func Attachment(path string, fileName string) Result {
	return nil
}

func Inline(path string, fileName string) Result {
	return nil
}

func Blob(content []byte, contentType string) Result {
	return nil
}

func BlobInline(content []byte, contentType, fileName string) Result {
	return nil
}

func BlobAttachment(content []byte, contentType, fileName string) Result {
	return nil
}

func Stream(reader io.Reader, contentType string) Result {
	return nil
}

func StreamInline(reader io.Reader, contentType, fileName string) Result {
	return nil
}

func StreamAttachment(reader io.Reader, contentType, fileName string) Result {
	return nil
}

func (f FileResult) StatusCode() Status {
	//TODO implement me
	panic("implement me")
}

func (f FileResult) Body() any {
	//TODO implement me
	panic("implement me")
}

func (f FileResult) Headers() map[string]string {
	//TODO implement me
	panic("implement me")
}

type ViewResult struct {
	Name   string
	Model  any
	Status Status
}

func View(name string) Result {
	return ViewResult{
		Name:   name,
		Status: StatusOK,
	}
}

func ViewModel(name string, model any) Result {
	return ViewResult{
		Name:   name,
		Model:  model,
		Status: StatusOK,
	}
}

func (v ViewResult) StatusCode() Status {
	return v.Status
}

func (v ViewResult) Body() any {
	return v.Model
}

func (v ViewResult) Headers() map[string]string {
	return nil
}

/*
type ContentResult struct {
	content         string
	contentType     string
	contentEncoding string
}

type JsonResult struct {
	Value      any
	StatusCode Status
}


*/
