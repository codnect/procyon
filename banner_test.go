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
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBannerPrinter_PrintBanner(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(writer *AnyMockWriter)

		wantOutput string
		wantErr    error
	}{
		{
			name: "no error",
			preCondition: func(writer *AnyMockWriter) {
				writer.On("Write", mock.Anything).Return(0, nil)
			},
			wantOutput: strings.Join(bannerText, "") + fmt.Sprintf(versionFormat, "(", Version),
		},
		{
			name: "text write error",
			preCondition: func(writer *AnyMockWriter) {
				writer.On("Write", mock.Anything).Return(0, errors.New("write error"))
			},
			wantErr:    errors.New("print banner: write error"),
			wantOutput: "",
		},
		{
			name: "version write error",
			preCondition: func(writer *AnyMockWriter) {
				for _, line := range bannerText {
					writer.On("Write", []byte(line)).Return(len(line), nil).Once()
				}

				writer.On("Write", []byte(fmt.Sprintf(versionFormat, "(", Version))).
					Return(0, errors.New("write error"))
			},
			wantErr:    errors.New("print banner: write error"),
			wantOutput: strings.Join(bannerText, ""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			anyWriter := &AnyMockWriter{}
			if tc.preCondition != nil {
				tc.preCondition(anyWriter)
			}

			bannerPrinter := NewBannerPrinter()

			// when
			err := bannerPrinter.Print(nil, anyWriter)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
			}

			assert.Equal(t, tc.wantOutput, anyWriter.String())
		})
	}
}
