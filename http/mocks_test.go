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
	"context"

	"github.com/stretchr/testify/mock"
)

type AnyDispatcher struct {
	mock.Mock
}

func (a *AnyDispatcher) Dispatch(ctx *Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

type AnyServer struct {
	mock.Mock
}

func (a *AnyServer) Start(ctx context.Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

func (a *AnyServer) Shutdown(ctx context.Context) error {
	result := a.Called(ctx)
	return result.Error(0)
}

func (a *AnyServer) Port() int {
	result := a.Called()
	return result.Int(0)
}

type AnyHandler struct {
	mock.Mock
}

func (h *AnyHandler) Handle(ctx *Context) (Result, error) {
	result := h.Called(ctx)

	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(Result), result.Error(1)
}

type AnyEndpointRegistrar struct {
	mock.Mock
}

func (r *AnyEndpointRegistrar) Register(endpoint *Endpoint) error {
	args := r.Called(endpoint)
	return args.Error(0)
}
