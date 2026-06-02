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

func TestMakeDefinition(t *testing.T) {
	testCases := []struct {
		name          string
		constructorFn ConstructorFunc
		opts          []DefinitionOption

		wantName     string
		wantScope    string
		wantType     reflect.Type
		wantErr      error
		wantArgNames []string
	}{
		{
			name:          "nil constructor",
			constructorFn: nil,
			wantErr:       errors.New("nil constructor function"),
		},
		{
			name: "invalid out type",
			constructorFn: func() string {
				return ""
			},
			wantErr: errors.New("constructor must return a struct, pointer to struct, or interface, got string"),
		},
		{
			name:          "without options",
			constructorFn: NewAnyPointerComponent,
			wantName:      "anyPointerComponent",
			wantScope:     SingletonScope,
			wantType:      reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name:          "with custom name",
			constructorFn: NewAnyPointerComponent,
			opts: []DefinitionOption{
				WithName("customName"),
			},
			wantName:  "customName",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name:          "with singleton scope",
			constructorFn: NewAnyPointerComponent,
			opts: []DefinitionOption{
				AsSingleton(),
			},
			wantName:  "anyPointerComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name:          "with prototype scope",
			constructorFn: NewAnyPointerComponent,
			opts: []DefinitionOption{
				AsPrototype(),
			},
			wantName:  "anyPointerComponent",
			wantScope: PrototypeScope,
			wantType:  reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name:          "with custom scope",
			constructorFn: NewAnyPointerComponent,
			opts: []DefinitionOption{
				WithScope("anyScope"),
			},
			wantName:  "anyPointerComponent",
			wantScope: "anyScope",
			wantType:  reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name:          "with qualifier",
			constructorFn: NewAnyDependentComponent,
			opts: []DefinitionOption{
				WithQualifierFor[AnySimpleComponent]("anyQualifier"),
			},
			wantName:  "anyDependentComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnyDependentComponent](),
			wantArgNames: []string{
				"anyQualifier",
			},
		},
		{
			name:          "with qualifier for missing constructor input",
			constructorFn: NewAnyDependentComponent,
			opts: []DefinitionOption{
				WithQualifierFor[string]("anyQualifier"),
			},
			wantName:  "anotherComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnyDependentComponent](),
			wantErr:   errors.New("constructor has no parameter of type string"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			def, err := MakeDefinition(tc.constructorFn, tc.opts...)

			// then
			if tc.wantErr != nil {
				require.Nil(t, def)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, def)

			assert.Equal(t, tc.wantName, def.Name())
			assert.Equal(t, tc.wantScope, def.Scope())

			if tc.wantScope == SingletonScope {
				assert.True(t, def.IsSingleton())
			} else if tc.wantScope == PrototypeScope {
				assert.True(t, def.IsPrototype())
			}

			assert.Equal(t, tc.wantType, def.Type())

			// constructor
			constructor := def.Constructor()
			require.NotNil(t, constructor)
			assert.Equal(t, tc.wantType, constructor.OutType())

			// constructor args
			args := constructor.Args()
			require.Equal(t, len(tc.wantArgNames), len(args))
			for index, wantArg := range tc.wantArgNames {
				assert.Equal(t, index, args[index].Index())
				assert.Equal(t, wantArg, args[index].Name())
			}
		})
	}
}
