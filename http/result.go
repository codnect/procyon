// Copyright 2026 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import "io"

// Result represents an HTTP response produced by a handler.
type Result interface {
	// Status returns the HTTP status code of the response.
	Status() Status
	// Header returns the HTTP headers of the response.
	Header() Header
}

type ValueResult interface {
	Result

	Value() any
}

type StatusCodeResult struct {
	StatusCode Status
	Headers    Header
}

func StatusCode(status Status) StatusCodeResult {
	return StatusCodeResult{
		StatusCode: status,
	}
}

func NoContent() StatusCodeResult {
	return StatusCodeResult{
		StatusCode: StatusNoContent,
	}
}

func (s StatusCodeResult) Status() Status {
	return s.StatusCode
}

func (s StatusCodeResult) Header() Header {
	return s.Headers
}

type BodyResult[T any] struct {
	Body       T
	StatusCode Status
	Headers    Header
}

func Ok[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusOK,
		Headers:    Header{},
	}

}

func Accepted[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusAccepted,
		Headers:    Header{},
	}
}

func Created[T any](location string, body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusCreated,
		Headers: Header{
			"Location": []string{location},
		},
	}
}

func BadRequest[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusBadRequest,
		Headers:    Header{},
	}
}

func Unauthorized[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusUnauthorized,
		Headers:    Header{},
	}
}

func Forbidden[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusForbidden,
		Headers:    Header{},
	}
}

func NotFound[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusNotFound,
		Headers:    Header{},
	}
}

func Conflict[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusConflict,
		Headers:    Header{},
	}
}

func UnprocessableEntity[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusUnprocessableEntity,
		Headers:    Header{},
	}
}

func InternalServerError[T any](body T) BodyResult[T] {
	return BodyResult[T]{
		Body:       body,
		StatusCode: StatusInternalServerError,
		Headers:    Header{},
	}
}

func (b BodyResult[T]) Header() Header {
	return b.Headers
}

func (b BodyResult[T]) Status() Status {
	return b.StatusCode
}

func (b BodyResult[T]) Value() any {
	return b.Body
}

type RedirectResult struct {
	Location   string
	StatusCode Status
}

func Redirect(url string) RedirectResult {
	return RedirectResult{
		Location:   url,
		StatusCode: StatusFound,
	}
}

func RedirectPermanent(url string) RedirectResult {
	return RedirectResult{
		Location:   url,
		StatusCode: StatusMovedPermanently,
	}
}

func RedirectPreserveMethod(url string) RedirectResult {
	return RedirectResult{
		Location:   url,
		StatusCode: StatusTemporaryRedirect,
	}
}

func RedirectPermanentPreserveMethod(url string) RedirectResult {
	return RedirectResult{
		Location:   url,
		StatusCode: StatusPermanentRedirect,
	}
}

func (r RedirectResult) Status() Status {
	return r.StatusCode
}

func (r RedirectResult) Header() Header {
	return Header{
		"Location": []string{r.Location},
	}
}

type ContentResult struct {
	Content     string
	ContentType string
	StatusCode  Status
}

func Text(content string) ContentResult {
	return ContentResult{
		Content:     content,
		ContentType: "text/plain",
		StatusCode:  StatusOK,
	}
}

func Content(content, contentType string) ContentResult {
	return ContentResult{
		Content:     content,
		ContentType: contentType,
		StatusCode:  StatusOK,
	}
}

func Html(content string) ContentResult {
	return ContentResult{
		Content:     content,
		ContentType: "text/html",
		StatusCode:  StatusOK,
	}
}

func (c ContentResult) Status() Status {
	return c.StatusCode
}

func (c ContentResult) Value() any {
	return c.Content
}

func (c ContentResult) Header() Header {
	return Header{
		"Content-Type": []string{c.ContentType},
	}
}

type JsonResult struct {
	Data       any
	StatusCode Status
	Headers    Header
}

func Json(data any, status Status) JsonResult {
	return JsonResult{
		Data:       data,
		StatusCode: status,
		Headers:    Header{"Content-Type": []string{"application/json"}},
	}
}

func (j JsonResult) Status() Status {
	return j.StatusCode
}

func (j JsonResult) Value() any {
	return j.Data
}

func (j JsonResult) Header() Header {
	return j.Headers
}

type FileDisposition string

const (
	FileDispositionInline     FileDisposition = "inline"
	FileDispositionAttachment FileDisposition = "attachment"
)

type FileResult interface {
	Result

	MediaType() string
	FileDownloadName() string
	FileDisposition() FileDisposition
}

type FileContentResult struct {
	Content      []byte
	ContentType  string
	DownloadName string
	Disposition  FileDisposition
	StatusCode   Status
	EntityTag    string
	LastModified string
}

func File(content []byte, contentType string) FileContentResult {
	return FileContentResult{
		Content:      content,
		ContentType:  contentType,
		DownloadName: "",
		StatusCode:   StatusOK,
	}
}

func FileInline(content []byte, contentType string) FileContentResult {
	return FileContentResult{
		Content:      content,
		ContentType:  contentType,
		DownloadName: "",
		Disposition:  FileDispositionInline,
		StatusCode:   StatusOK,
	}
}

func FileAttachment(content []byte, contentType, downloadName string) FileContentResult {
	return FileContentResult{
		Content:      content,
		ContentType:  contentType,
		DownloadName: downloadName,
		Disposition:  FileDispositionAttachment,
		StatusCode:   StatusOK,
	}
}

func (f FileContentResult) Status() Status {
	return f.StatusCode
}

func (f FileContentResult) Header() Header {

	h := Header{
		"Content-Type": []string{f.ContentType},
	}
	if f.Disposition != "" {
		value := string(f.Disposition)
		if f.DownloadName != "" {
			value += `; filename="` + f.DownloadName + `"`
		}
		h.Set("Content-Disposition", value)
	}
	if f.EntityTag != "" {
		h.Set("ETag", f.EntityTag)
	}
	if f.LastModified != "" {
		h.Set("Last-Modified", f.LastModified)
	}
	return h

}
func (f FileContentResult) Value() any {
	return f.Content
}

func (f FileContentResult) MediaType() string {
	return f.ContentType
}

func (f FileContentResult) FileDownloadName() string {
	return f.DownloadName
}

func (f FileContentResult) FileDisposition() FileDisposition {
	return f.Disposition
}

type PhysicalFileResult struct {
	FilePath     string
	ContentType  string
	DownloadName string
	Disposition  FileDisposition
	StatusCode   Status
}

func PhysicalFile(filePath, contentType string) PhysicalFileResult {
	return PhysicalFileResult{
		FilePath:    filePath,
		ContentType: contentType,
		StatusCode:  StatusOK,
	}
}

func (f PhysicalFileResult) Status() Status {
	return f.StatusCode
}

func (f PhysicalFileResult) Header() Header {
	h := Header{
		"Content-Type": []string{f.ContentType},
	}
	if f.Disposition != "" {
		value := string(f.Disposition)
		if f.DownloadName != "" {
			value += `; filename="` + f.DownloadName + `"`
		}
		h.Set("Content-Disposition", value)
	}
	return h
}

func (f PhysicalFileResult) MediaType() string {
	return f.ContentType
}

func (f PhysicalFileResult) FileDownloadName() string {
	return f.DownloadName
}

func (f PhysicalFileResult) FileDisposition() FileDisposition {
	return f.Disposition
}

type StreamResult struct {
	Reader      io.Reader
	ContentType string
	StatusCode  Status
	Headers     Header
}

func Stream(reader io.Reader, contentType string) StreamResult {
	return StreamResult{
		Reader:      reader,
		ContentType: contentType,
		StatusCode:  StatusOK,
	}
}

func (s StreamResult) Status() Status {
	return s.StatusCode
}

func (s StreamResult) Value() any {
	return s.Reader
}

func (s StreamResult) Header() Header {
	h := s.Headers.Clone()
	h.Set("Content-Type", s.ContentType)
	return h
}

func Problem() Result {
	return nil
}

func ValidationProblem() Result {
	return nil
}

type ResultExecutorRegistry struct {
	executors []ResultExecutor
}

func (r *ResultExecutorRegistry) Register(executor ResultExecutor) {
	r.executors = append(r.executors, executor)

}

func (r *ResultExecutorRegistry) Resolve(result Result) (ResultExecutor, bool) {
	for _, executor := range r.executors {
		if executor.CanExecute(result) {
			return executor, true
		}
	}
	return nil, false
}

type ResultExecutor interface {
	CanExecute(result Result) bool
	Execute(ctx *Context, result Result) error
}
