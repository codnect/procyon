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

package procyon

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AnyWriter struct {
	mock.Mock
	buf bytes.Buffer
}

func (w *AnyWriter) Write(p []byte) (int, error) {
	args := w.Called(p)
	err := args.Error(1)

	if err == nil {
		w.buf.Write(p)
	}

	return args.Int(0), err
}

func (w *AnyWriter) String() string {
	return w.buf.String()
}

func TestBannerPrinter_PrintBanner(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(writer *AnyWriter)

		wantOutput string
		wantErr    error
	}{
		{
			name: "no error",
			preCondition: func(writer *AnyWriter) {
				writer.On("Write", mock.Anything).Return(0, nil)
			},
			wantOutput: strings.Join(bannerText, "") + fmt.Sprintf(versionFormat, "(", Version),
		},
		{
			name: "text write error",
			preCondition: func(writer *AnyWriter) {
				writer.On("Write", mock.Anything).Return(0, errors.New("write error"))
			},
			wantErr:    errors.New("write error"),
			wantOutput: "",
		},
		{
			name: "version write error",
			preCondition: func(writer *AnyWriter) {
				for _, line := range bannerText {
					writer.On("Write", []byte(line)).Return(len(line), nil).Once()
				}

				writer.On("Write", []byte(fmt.Sprintf(versionFormat, "(", Version))).
					Return(0, errors.New("write error"))
			},
			wantErr:    errors.New("write error"),
			wantOutput: strings.Join(bannerText, ""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			anyWriter := &AnyWriter{}
			if tc.preCondition != nil {
				tc.preCondition(anyWriter)
			}

			bannerPrinter := NewBannerPrinter()

			// when
			err := bannerPrinter.PrintBanner(anyWriter)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
			}

			assert.Equal(t, tc.wantOutput, anyWriter.String())
		})
	}
}
