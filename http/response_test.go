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

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerResponse_Context(t *testing.T) {
	// given
	requestCtx := &Context{}
	serverResponse := &ServerResponse{
		ctx: requestCtx,
	}

	// when
	ctx := serverResponse.Context()

	// then
	assert.Equal(t, requestCtx, ctx)
}

func TestServerResponse_AddCookie(t *testing.T) {
	testCases := []struct {
		name          string
		cookie        *Cookie
		preCondition  func(serverResponse *ServerResponse)
		postCondition func(t *testing.T, serverResponse *ServerResponse)
	}{
		{
			name: "headers already written",
			cookie: &Cookie{
				Name:  "anotherTestCookie",
				Value: "anotherTestCookieValue",
			},
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.AddCookie(&Cookie{
					Name:  "testCookie",
					Value: "testCookieValue",
				})
				serverResponse.Writer()
			},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("Set-Cookie")
				assert.True(t, ok)
				assert.Equal(t, "testCookie=testCookieValue; Path=/", header)
			},
		},
		{
			name: "headers not written",
			cookie: &Cookie{
				Name:  "testCookie",
				Value: "testCookieValue",
			},
			preCondition: func(serverResponse *ServerResponse) {},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("Set-Cookie")
				assert.True(t, ok)
				assert.Equal(t, "testCookie=testCookieValue; Path=/", header)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			serverResponse.AddCookie(tc.cookie)

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, serverResponse)
			}
		})
	}
}

func TestServerResponse_AddHeader(t *testing.T) {
	testCases := []struct {
		name          string
		headerName    string
		headerValue   string
		preCondition  func(serverResponse *ServerResponse)
		postCondition func(t *testing.T, serverResponse *ServerResponse)
	}{
		{
			name:        "headers already written",
			headerName:  "anotherTestHeader",
			headerValue: "anotherTestHeaderValue",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.AddHeader("testHeader", "testValue")
				serverResponse.Writer()
			},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("testHeader")
				assert.True(t, ok)
				assert.Equal(t, "testValue", header)

				header, ok = serverResponse.Header("anotherTestHeader")
				assert.False(t, ok)
				assert.Empty(t, header)
			},
		},
		{
			name:         "headers not written",
			headerName:   "testHeader",
			headerValue:  "testValue",
			preCondition: func(serverResponse *ServerResponse) {},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("testHeader")
				assert.True(t, ok)
				assert.Equal(t, "testValue", header)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			serverResponse.AddHeader(tc.headerName, tc.headerValue)

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, serverResponse)
			}
		})
	}
}

func TestServerResponse_SetHeader(t *testing.T) {
	testCases := []struct {
		name          string
		headerName    string
		headerValue   string
		preCondition  func(serverResponse *ServerResponse)
		postCondition func(t *testing.T, serverResponse *ServerResponse)
	}{
		{
			name:        "headers already written",
			headerName:  "anotherTestHeader",
			headerValue: "anotherTestHeaderValue",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetHeader("testHeader", "testValue")
				serverResponse.Writer()
			},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("testHeader")
				assert.True(t, ok)
				assert.Equal(t, "testValue", header)

				header, ok = serverResponse.Header("anotherTestHeader")
				assert.False(t, ok)
				assert.Empty(t, header)
			},
		},
		{
			name:         "headers not written",
			headerName:   "testHeader",
			headerValue:  "testValue",
			preCondition: func(serverResponse *ServerResponse) {},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("testHeader")
				assert.True(t, ok)
				assert.Equal(t, "testValue", header)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			serverResponse.SetHeader(tc.headerName, tc.headerValue)

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, serverResponse)
			}
		})
	}
}

func TestServerResponse_DeleteHeader(t *testing.T) {
	testCases := []struct {
		name          string
		headerName    string
		preCondition  func(serverResponse *ServerResponse)
		postCondition func(t *testing.T, serverResponse *ServerResponse)
	}{
		{
			name:       "headers already written",
			headerName: "testHeader",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetHeader("testHeader", "testValue")
				serverResponse.Writer()
			},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("testHeader")
				assert.True(t, ok)
				assert.Equal(t, "testValue", header)
			},
		},
		{
			name:         "headers not written",
			headerName:   "testHeader",
			preCondition: func(serverResponse *ServerResponse) {},
			postCondition: func(t *testing.T, serverResponse *ServerResponse) {
				header, ok := serverResponse.Header("testHeader")
				assert.False(t, ok)
				assert.Empty(t, header)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			serverResponse.DeleteHeader(tc.headerName)

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, serverResponse)
			}
		})
	}
}

func TestServerResponse_Header(t *testing.T) {
	testCases := []struct {
		name            string
		headerName      string
		preCondition    func(serverResponse *ServerResponse)
		wantHeaderValue string
		wantExists      bool
	}{
		{
			name:            "empty header name",
			headerName:      "",
			preCondition:    func(serverResponse *ServerResponse) {},
			wantHeaderValue: "",
			wantExists:      false,
		},
		{
			name:            "no header",
			preCondition:    func(serverResponse *ServerResponse) {},
			headerName:      "testHeader",
			wantHeaderValue: "",
			wantExists:      false,
		},
		{
			name: "header exists",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetHeader("testHeader", "testValue")
			},
			headerName:      "testHeader",
			wantHeaderValue: "testValue",
			wantExists:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			val, ok := serverResponse.Header(tc.headerName)

			// then
			assert.Equal(t, tc.wantExists, ok)
			assert.Equal(t, tc.wantHeaderValue, val)
		})
	}
}

func TestServerResponse_HeaderValues(t *testing.T) {
	testCases := []struct {
		name             string
		headerName       string
		preCondition     func(serverResponse *ServerResponse)
		wantHeaderValues []string
		wantExists       bool
	}{
		{
			name:             "empty header name",
			headerName:       "",
			preCondition:     func(serverResponse *ServerResponse) {},
			wantHeaderValues: []string{},
			wantExists:       false,
		},
		{
			name:             "no header",
			preCondition:     func(serverResponse *ServerResponse) {},
			headerName:       "testHeader",
			wantHeaderValues: []string{},
			wantExists:       false,
		},
		{
			name: "header exists",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.AddHeader("testHeader", "testValue")
				serverResponse.AddHeader("testHeader", "anotherTestValue")
			},
			headerName:       "testHeader",
			wantHeaderValues: []string{"testValue", "anotherTestValue"},
			wantExists:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			val := serverResponse.HeaderValues(tc.headerName)

			// then
			assert.ElementsMatch(t, tc.wantHeaderValues, val)
		})
	}
}

func TestServerResponse_Status(t *testing.T) {
	// given
	serverResponse := &ServerResponse{
		status: StatusOK,
	}

	// when
	status := serverResponse.Status()

	// then
	assert.Equal(t, StatusOK, status)
}

func TestServerResponse_SetStatus(t *testing.T) {
	testCases := []struct {
		name         string
		status       Status
		preCondition func(serverResponse *ServerResponse)
		wantStatus   Status
	}{
		{
			name: "headers already written",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.Writer()
			},
			status:     StatusBadRequest,
			wantStatus: StatusOK,
		},
		{
			name:       "headers not written",
			status:     StatusBadRequest,
			wantStatus: StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			serverResponse.SetStatus(tc.status)

			// then
			assert.Equal(t, tc.wantStatus, serverResponse.Status())
		})
	}
}

func TestServerResponse_Flush(t *testing.T) {
	testCases := []struct {
		name          string
		preCondition  func(serverResponse *ServerResponse)
		postCondition func(t *testing.T, recoder *httptest.ResponseRecorder)
		wantErr       error
	}{
		{
			name: "already committed",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetStatus(400)
				serverResponse.Writer()
			},
			wantErr: errors.New("response already committed"),
		},
		{
			name: "not committed",
			postCondition: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.True(t, recoder.Flushed)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			err := serverResponse.Flush()

			// then
			if tc.postCondition != nil {
				tc.postCondition(t, recorder)
			}

			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestServerResponse_IsCommitted(t *testing.T) {
	testCases := []struct {
		name          string
		preCondition  func(serverResponse *ServerResponse)
		wantCommitted bool
	}{
		{
			name:          "no committed",
			preCondition:  func(serverResponse *ServerResponse) {},
			wantCommitted: false,
		},
		{
			name: "committed",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetStatus(StatusOK)
				serverResponse.Writer()
			},
			wantCommitted: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			isCommitted := serverResponse.IsCommitted()

			// then
			assert.Equal(t, tc.wantCommitted, isCommitted)
		})
	}
}

func TestServerResponse_Redirect(t *testing.T) {
	testCases := []struct {
		name         string
		location     string
		status       Status
		preCondition func(serverResponse *ServerResponse)
		wantStatus   Status
		wantErr      error
	}{
		{
			name:     "already committed",
			location: "/committed",
			status:   StatusOK,
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetStatus(400)
				serverResponse.Writer()
			},
			wantStatus: StatusBadRequest,
			wantErr:    errors.New("response already committed"),
		},
		{
			name:     "not committed",
			location: "/committed",
			status:   StatusOK,
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetStatus(400)
			},
			wantStatus: StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			err := serverResponse.Redirect(tc.location, tc.status)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, serverResponse.Status())
		})
	}
}

func TestServerResponse_Reset(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(serverResponse *ServerResponse)
		wantStatus   Status
		wantErr      error
	}{
		{
			name: "already committed",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetStatus(400)
				serverResponse.Writer()
			},
			wantStatus: StatusBadRequest,
			wantErr:    errors.New("response already committed"),
		},
		{
			name: "not committed",
			preCondition: func(serverResponse *ServerResponse) {
				serverResponse.SetStatus(400)
			},
			wantStatus: StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			recorder := httptest.NewRecorder()
			serverResponse := &ServerResponse{
				ctx:     &Context{},
				headers: Header{},
				writer:  recorder,
			}

			if tc.preCondition != nil {
				tc.preCondition(serverResponse)
			}

			// when
			err := serverResponse.Reset()

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, serverResponse.Status())
		})
	}
}
