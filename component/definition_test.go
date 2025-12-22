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
			wantErr:       errors.New("nil constructor"),
		},
		{
			name: "invalid out type",
			constructorFn: func() string {
				return ""
			},
			wantErr: errors.New("invalid constructor output: expected struct, pointer, or interface; got string"),
		},
		{
			name:          "without options",
			constructorFn: NewAnyComponent,
			wantName:      "anyComponent",
			wantScope:     SingletonScope,
			wantType:      reflect.TypeFor[*AnyComponent](),
		},
		{
			name:          "with custom name",
			constructorFn: NewAnyComponent,
			opts: []DefinitionOption{
				WithName("customName"),
			},
			wantName:  "customName",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnyComponent](),
		},
		{
			name:          "with singleton scope",
			constructorFn: NewAnyComponent,
			opts: []DefinitionOption{
				AsSingleton(),
			},
			wantName:  "anyComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnyComponent](),
		},
		{
			name:          "with prototype scope",
			constructorFn: NewAnyComponent,
			opts: []DefinitionOption{
				AsPrototype(),
			},
			wantName:  "anyComponent",
			wantScope: PrototypeScope,
			wantType:  reflect.TypeFor[*AnyComponent](),
		},
		{
			name:          "with custom scope",
			constructorFn: NewAnyComponent,
			opts: []DefinitionOption{
				WithScope("anyScope"),
			},
			wantName:  "anyComponent",
			wantScope: "anyScope",
			wantType:  reflect.TypeFor[*AnyComponent](),
		},
		{
			name:          "with qualifier",
			constructorFn: NewAnotherComponent,
			opts: []DefinitionOption{
				WithQualifierFor[DependentComponent]("anyQualifier"),
			},
			wantName:  "anotherComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnotherComponent](),
			wantArgNames: []string{
				"anyQualifier",
			},
		},
		{
			name:          "with qualifier for missing constructor input",
			constructorFn: NewAnotherComponent,
			opts: []DefinitionOption{
				WithQualifierFor[string]("anyQualifier"),
			},
			wantName:  "anotherComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnotherComponent](),
			wantErr:   errors.New("no constructor input of type string"),
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
