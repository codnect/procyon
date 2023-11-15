package rest

import "net/http"

type ResponseEntity struct {
	status  int
	body    any
	headers http.Header
}

func (re ResponseEntity) Body() any {
	return re.body
}

func (re ResponseEntity) Headers() http.Header {
	return re.headers
}

func (re ResponseEntity) HasBody() bool {
	return re.body != nil
}

func (re ResponseEntity) Status() int {
	return re.status
}
