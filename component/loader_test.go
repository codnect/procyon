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

package component

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConditionalLoader(t *testing.T) {
	testCases := []struct {
		name       string
		ctx        context.Context
		container  Container
		components []*Component
		wantPanic  error
	}{
		{
			name:       "nil container",
			ctx:        context.Background(),
			container:  nil,
			components: []*Component{},
			wantPanic:  errors.New("nil container"),
		},
		{
			name:       "valid container",
			ctx:        context.Background(),
			container:  NewStandardContainer(),
			components: []*Component{},
		},
		{
			name:       "nil components",
			ctx:        context.Background(),
			container:  NewStandardContainer(),
			components: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewConditionalLoader(tc.container, tc.components)
				})
				return
			}

			loader := NewConditionalLoader(tc.container, tc.components)

			// then
			require.NotNil(t, loader)
		})
	}
}

func TestConditionalLoader_Load(t *testing.T) {
	anyComponentDef, _ := MakeDefinition(NewAnyPointerComponent)
	require.NotNil(t, anyComponentDef)

	anotherComponentDef, _ := MakeDefinition(NewAnyDependentComponent)
	require.NotNil(t, anotherComponentDef)

	testCases := []struct {
		name       string
		ctx        context.Context
		container  Container
		components []*Component
		wantErr    error
		wantTypes  []reflect.Type
	}{
		{
			name:       "nil context",
			ctx:        nil,
			container:  NewStandardContainer(),
			components: []*Component{},
			wantErr:    errors.New("nil context"),
		},
		{
			name:      "load component",
			ctx:       context.Background(),
			container: NewStandardContainer(),
			components: []*Component{
				Create(anyComponentDef),
			},
			wantTypes: []reflect.Type{
				reflect.TypeOf(&AnyPointerComponent{}),
			},
		},
		{
			name:      "skip component",
			ctx:       context.Background(),
			container: NewStandardContainer(),
			components: []*Component{
				Create(anyComponentDef, AnyCondition{matches: false}),
				Create(anotherComponentDef),
			},
			wantTypes: []reflect.Type{
				reflect.TypeOf(&AnyDependentComponent{}),
			},
		},
		{
			name:      "already exists",
			ctx:       context.Background(),
			container: NewStandardContainer(),
			components: []*Component{
				Create(anyComponentDef),
				Create(anyComponentDef),
			},
			wantErr: fmt.Errorf("load component \"anyPointerComponent\": register definition \"anyPointerComponent\": duplicate definition"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			loader := NewConditionalLoader(tc.container, tc.components)

			// when
			err := loader.Load(tc.ctx)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			for _, wantType := range tc.wantTypes {
				result := tc.container.CanResolveType(wantType)
				require.True(t, result)
			}
		})
	}
}
