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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type AnyInterface interface {
	AnyMethod()
}

type AnyComponent struct{}

func (a AnyComponent) AnyMethod() {}

func NewAnyComponent() *AnyComponent {
	return &AnyComponent{}
}

type AnotherComponent struct{}

func (a AnotherComponent) AnyMethod() {}

func NewAnotherComponent(component DependentComponent) *AnotherComponent {
	return &AnotherComponent{}
}

type DependentComponent struct {
}

func (d DependentComponent) AnyMethod() {

}

func TestRegister(t *testing.T) {
	testCases := []struct {
		name           string
		preCondition   func()
		constructorFn  ConstructorFunc
		opts           []DefinitionOption
		conditions     []Condition
		wantName       string
		wantScope      string
		wantType       reflect.Type
		wantConditions []Condition
		wantPanic      error
	}{
		{
			name:          "nil constructor",
			constructorFn: nil,
			wantPanic:     errors.New("nil constructor"),
		},
		{
			name: "already exists",
			preCondition: func() {
				components["anyComponent"] = &Component{}
			},
			constructorFn: NewAnyComponent,
			wantPanic:     errors.New("component with name 'anyComponent' already exists"),
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
			name:          "with prototype scope",
			constructorFn: NewAnyComponent,
			opts: []DefinitionOption{
				WithScope(PrototypeScope),
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
			name:          "with conditions",
			constructorFn: NewAnyComponent,
			conditions: []Condition{
				AnyCondition{},
			},
			wantName:  "anyComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[*AnyComponent](),
			wantConditions: []Condition{
				AnyCondition{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// cleanup
			clear(components)

			// given
			if tc.preCondition != nil {
				tc.preCondition()
			}

			// when
			if tc.wantPanic != nil {
				require.PanicsWithError(t, tc.wantPanic.Error(), func() {
					Register(tc.constructorFn, tc.opts...)
				})
				return
			}

			reg := Register(tc.constructorFn, tc.opts...)
			require.NotNil(t, reg, "nil registration")

			// Apply conditions if any
			for _, cond := range tc.conditions {
				reg.Conditional(cond)
			}

			// then
			require.Contains(t, components, tc.wantName)
			component := components[tc.wantName]

			// Definition checks
			def := component.Definition()
			require.NotNil(t, def, "nil definition")
			assert.Equal(t, tc.wantName, def.Name())
			assert.Equal(t, tc.wantType, def.Type())
			assert.Equal(t, tc.wantScope, def.Scope())
			require.NotNil(t, def.Constructor(), "nil constructor")

			// Condition check
			require.Len(t, component.Conditions(), len(tc.wantConditions))
		})
	}
}

func TestList(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyComponent).Conditional(AnyCondition{})

	// when
	componentList := List()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anyComponent", def.Name())

	// Condition check
	require.Len(t, component.Conditions(), 1)
}

func TestList_MultipleComponent(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyComponent).Conditional(AnyCondition{})
	Register(NewAnotherComponent).Conditional(nil)

	// when
	componentList := List()

	// then
	require.Len(t, componentList, 2)

	assert.Contains(t, components, "anyComponent")
	assert.Contains(t, components, "anotherComponent")
}

func TestListOf_WithNonRegisteredType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnotherComponent)

	// when
	componentList := ListOf[*AnyComponent]()

	// then
	require.Len(t, componentList, 0)
}

func TestListOf_WithPointerStructType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyComponent)

	// when
	componentList := ListOf[*AnyComponent]()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anyComponent", def.Name())
}

func TestListOf_WithNonPointerStructType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyComponent)

	// when
	componentList := ListOf[AnyComponent]()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anyComponent", def.Name())
}

func TestListOf_WithInterfaceType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyComponent)

	// when
	componentList := ListOf[AnyInterface]()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anyComponent", def.Name())
}
