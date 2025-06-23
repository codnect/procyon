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
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
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
			container:  NewDefaultContainer(),
			components: []*Component{},
		},
		{
			name:       "nil components",
			ctx:        context.Background(),
			container:  NewDefaultContainer(),
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
	anyComponentDef, _ := MakeDefinition(NewAnyComponent)
	require.NotNil(t, anyComponentDef)

	anotherComponentDef, _ := MakeDefinition(NewAnotherComponent)
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
			name:      "load component",
			container: NewDefaultContainer(),
			components: []*Component{
				createComponent(anyComponentDef),
			},
			wantTypes: []reflect.Type{
				reflect.TypeOf(&AnyComponent{}),
			},
		},
		{
			name:      "skip component",
			ctx:       context.Background(),
			container: NewDefaultContainer(),
			components: []*Component{
				createComponent(anyComponentDef, AnyCondition{matches: false}),
				createComponent(anotherComponentDef),
			},
			wantTypes: []reflect.Type{
				reflect.TypeOf(&AnotherComponent{}),
			},
		},
		{
			name:      "already exists",
			container: NewDefaultContainer(),
			components: []*Component{
				createComponent(anyComponentDef),
				createComponent(anyComponentDef),
			},
			wantErr: fmt.Errorf("failed to register component \"anyComponent\": %s", ErrDefinitionAlreadyExists),
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
