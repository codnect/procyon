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
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateConstructor(t *testing.T) {
	testCases := []struct {
		name          string
		constructorFn ConstructorFunc

		wantOutType reflect.Type
		wantArgs    []Arg
		wantErr     error
	}{
		{
			name:          "nil constructor function",
			constructorFn: nil,
			wantErr:       errors.New("nil constructor function"),
		},
		{
			name:          "invalid constructor function",
			constructorFn: "string value",
			wantErr:       errors.New("constructor is not a function"),
		},
		{
			name: "multi result constructor function",
			constructorFn: func() (string, int, error) {
				return "", -1, nil
			},
			wantErr: errors.New("constructor must return exactly one result"),
		},
		{
			name:          "valid constructor function",
			constructorFn: NewAnyDependentComponent,
			wantOutType:   reflect.TypeFor[*AnyDependentComponent](),
			wantArgs: []Arg{
				{
					0,
					"",
					reflect.TypeFor[AnySimpleComponent](),
					false,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			constructor, err := createConstructor(tc.constructorFn)

			// then
			if tc.wantErr != nil {
				require.Equal(t, Constructor{}, constructor)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantOutType, constructor.OutType())
			assert.Len(t, constructor.Args(), len(tc.wantArgs))

			if len(tc.wantArgs) == len(constructor.args) {
				for index, wantArg := range tc.wantArgs {
					arg := constructor.args[index]

					assert.Equal(t, wantArg.Index(), arg.Index())
					assert.Equal(t, wantArg.Name(), arg.Name())
					assert.Equal(t, wantArg.Type(), arg.Type())
				}
			}
		})
	}
}

func TestConstructor_Invoke(t *testing.T) {
	testCases := []struct {
		name          string
		constructorFn ConstructorFunc
		inputs        []any

		wantOutType reflect.Type
		wantErr     error
	}{
		{
			name:          "invalid parameter count",
			constructorFn: NewAnyDependentComponent,
			wantErr:       errors.New("invalid argument count: got 0, want 1"),
		},
		{
			name: "invalid variadic parameter count",
			constructorFn: func(dependency AnySimpleComponent, another ...AnyDependentComponent) any {
				return nil
			},
			wantErr: errors.New("invalid argument count: got 0, want at least 1"),
		},
		{
			name:          "invalid parameter",
			constructorFn: NewAnyDependentComponent,
			inputs: []any{
				"invalid parameter",
			},
			wantErr: errors.New("argument 0 has type string, want component.AnySimpleComponent"),
		},
		{
			name:          "nil parameter",
			constructorFn: NewAnyDependentComponent,
			inputs: []any{
				nil,
			},
			wantOutType: reflect.TypeFor[*AnyDependentComponent](),
		},
		{
			name:          "valid constructor",
			constructorFn: NewAnyDependentComponent,
			inputs: []any{
				AnySimpleComponent{},
			},
			wantOutType: reflect.TypeFor[*AnyDependentComponent](),
		},
		{
			name: "invalid variadic parameter",
			constructorFn: func(dependencies ...AnySimpleComponent) *AnyDependentComponent {
				return &AnyDependentComponent{}
			},
			inputs: []any{
				"invalid parameter",
			},
			wantErr: errors.New("argument 0 has type string, want component.AnySimpleComponent"),
		},
		{
			name: "nil variadic parameter",
			constructorFn: func(dependencies ...AnySimpleComponent) *AnyDependentComponent {
				return &AnyDependentComponent{}
			},
			inputs: []any{
				nil,
			},
			wantOutType: reflect.TypeFor[*AnyDependentComponent](),
		},
		{
			name: "valid variadic constructor",
			constructorFn: func(dependencies ...AnySimpleComponent) *AnyDependentComponent {
				return &AnyDependentComponent{}
			},
			inputs: []any{
				AnySimpleComponent{},
				AnySimpleComponent{},
			},
			wantOutType: reflect.TypeFor[*AnyDependentComponent](),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			constructor, err := createConstructor(tc.constructorFn)
			require.NoError(t, err)

			// when
			out, invokeErr := constructor.Invoke(tc.inputs...)

			// then
			if tc.wantErr != nil {
				require.Error(t, invokeErr)
				require.EqualError(t, invokeErr, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantOutType, reflect.TypeOf(out))

		})
	}
}
