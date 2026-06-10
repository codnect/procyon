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

package procyon

import (
	"bytes"

	"github.com/stretchr/testify/mock"
)

type AnyMockWriter struct {
	mock.Mock
	buf bytes.Buffer
}

func (w *AnyMockWriter) Write(p []byte) (int, error) {
	args := w.Called(p)
	err := args.Error(1)

	if err == nil {
		w.buf.Write(p)
	}

	return args.Int(0), err
}

func (w *AnyMockWriter) String() string {
	return w.buf.String()
}
