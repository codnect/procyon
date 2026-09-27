// Copyright 2025 Codnect
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

package condition

import (
	"context"
	stdio "io"

	"github.com/stretchr/testify/mock"
	"go.codnect.io/procyon/io"
)

type AnyMockResource struct {
	mock.Mock
}

func (a *AnyMockResource) Name() string {
	result := a.Called()
	if result.Get(0) == nil {
		return ""
	}

	return result.Get(0).(string)
}

func (a *AnyMockResource) Location() string {
	result := a.Called()
	if result.Get(0) == nil {
		return ""
	}

	return result.Get(0).(string)
}

func (a *AnyMockResource) Exists() bool {
	result := a.Called()
	return result.Bool(0)
}

func (a *AnyMockResource) Reader() (stdio.ReadCloser, error) {
	result := a.Called()
	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(stdio.ReadCloser), result.Error(1)
}

type AnyMockResourceResolver struct {
	mock.Mock
}

func (a *AnyMockResourceResolver) Resolve(ctx context.Context, location string) (io.Resource, error) {
	result := a.Called(ctx, location)
	if result.Get(0) == nil {
		return nil, result.Error(1)
	}

	return result.Get(0).(io.Resource), nil
}
