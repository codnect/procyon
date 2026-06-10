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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name       string
		def        *Definition
		conditions []Condition

		wantPanic error
	}{
		{
			name:       "nil definition",
			def:        nil,
			conditions: []Condition{},
			wantPanic:  errors.New("component: nil definition"),
		},
		{
			name:       "valid definition",
			def:        &Definition{},
			conditions: []Condition{},
		},
		{

			name: "valid definition with conditions",
			def:  &Definition{},
			conditions: []Condition{
				AnyCondition{matches: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					Create(tc.def, tc.conditions...)
				})
				return
			}

			comp := Create(tc.def, tc.conditions...)
			require.NotNil(t, comp, "nil component")
		})
	}
}

func TestRegister(t *testing.T) {
	testCases := []struct {
		name          string
		preCondition  func()
		constructorFn ConstructorFunc
		opts          []DefinitionOption
		conditions    []Condition

		wantName       string
		wantScope      string
		wantType       reflect.Type
		wantConditions []Condition
		wantPanic      error
	}{
		{
			name:          "nil constructor function",
			constructorFn: nil,
			wantPanic:     errors.New("component: nil constructor function"),
		},
		{
			name: "already exists",
			preCondition: func() {
				components["anySimpleComponent"] = &Component{}
			},
			constructorFn: NewAnySimpleComponent,
			wantPanic:     errors.New("component: duplicate component name 'anySimpleComponent'"),
		},
		{
			name:          "without options",
			constructorFn: NewAnySimpleComponent,
			wantName:      "anySimpleComponent",
			wantScope:     SingletonScope,
			wantType:      reflect.TypeFor[AnySimpleComponent](),
		},
		{
			name:          "with custom name",
			constructorFn: NewAnySimpleComponent,
			opts: []DefinitionOption{
				WithName("customName"),
			},
			wantName:  "customName",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[AnySimpleComponent](),
		},
		{
			name:          "with prototype scope",
			constructorFn: NewAnySimpleComponent,
			opts: []DefinitionOption{
				WithScope(PrototypeScope),
			},
			wantName:  "anySimpleComponent",
			wantScope: PrototypeScope,
			wantType:  reflect.TypeFor[AnySimpleComponent](),
		},
		{
			name:          "with custom scope",
			constructorFn: NewAnySimpleComponent,
			opts: []DefinitionOption{
				WithScope("anyScope"),
			},
			wantName:  "anySimpleComponent",
			wantScope: "anyScope",
			wantType:  reflect.TypeFor[AnySimpleComponent](),
		},
		{
			name:          "with conditions",
			constructorFn: NewAnySimpleComponent,
			conditions: []Condition{
				AnyCondition{},
			},
			wantName:  "anySimpleComponent",
			wantScope: SingletonScope,
			wantType:  reflect.TypeFor[AnySimpleComponent](),
			wantConditions: []Condition{
				AnyCondition{},
			},
		},
		{
			name: "multi return values",
			constructorFn: func() (AnySimpleComponent, error) {
				return AnySimpleComponent{}, nil
			},
			wantPanic: errors.New("component: constructor must return exactly one result"),
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
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
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

func TestLoad_NotFound(t *testing.T) {
	// given
	clear(components)

	// when
	instance, err := Load[*AnyPointerComponent]("anyPointerComponent")

	// then
	require.Nil(t, instance)
	require.NotNil(t, err)

	assert.Equal(t, "load component \"anyPointerComponent\": not found", err.Error())
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestLoad_NotConvertible(t *testing.T) {
	// given
	clear(components)

	Register(NewAnySimpleComponent)

	// when
	instance, err := Load[AnyPointerComponent]("anySimpleComponent")

	// then
	require.Zero(t, instance)
	require.NotNil(t, err)

	assert.Equal(t, "load component \"anySimpleComponent\": component.AnySimpleComponent is not convertible to component.AnyPointerComponent: type mismatch", err.Error())
	assert.ErrorIs(t, err, ErrTypeMismatch)

}

func TestLoad_InvalidArgument(t *testing.T) {
	// given
	clear(components)

	Register(NewAnyPointerComponent)

	// when
	instance, err := Load[*AnyPointerComponent]("anyPointerComponent", context.Background())

	// then
	require.Nil(t, instance)
	require.NotNil(t, err)

	assert.Equal(t, "load component \"anyPointerComponent\": invalid argument count: got 1, want 0", err.Error())
}

func TestLoad(t *testing.T) {
	// given
	clear(components)

	Register(NewAnyPointerComponent)

	// when
	instance, err := Load[*AnyPointerComponent]("anyPointerComponent")

	// then
	require.NotNil(t, instance)
	require.Nil(t, err)
}

func TestList(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnySimpleComponent).Conditional(AnyCondition{})

	// when
	componentList := List()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anySimpleComponent", def.Name())

	// Condition check
	require.Len(t, component.Conditions(), 1)
}

func TestList_MultipleComponent(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnySimpleComponent).Conditional(AnyCondition{})
	Register(NewAnyPointerComponent).Conditional(nil)

	// when
	componentList := List()

	// then
	require.Len(t, componentList, 2)

	assert.Contains(t, components, "anySimpleComponent")
	assert.Contains(t, components, "anyPointerComponent")
}

func TestListOf_NonRegisteredType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnySimpleComponent)

	// when
	componentList := ListOf[*AnyPointerComponent]()

	// then
	require.Len(t, componentList, 0)
}

func TestListOf_PointerStructType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyPointerComponent)

	// when
	componentList := ListOf[*AnyPointerComponent]()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anyPointerComponent", def.Name())
}

func TestListOf_NonPointerStructType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyPointerComponent)

	// when
	componentList := ListOf[AnyPointerComponent]()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anyPointerComponent", def.Name())
}

func TestListOf_InterfaceType(t *testing.T) {
	// cleanup
	clear(components)

	// given
	Register(NewAnyPointerComponent)

	// when
	componentList := ListOf[AnyComponent]()

	// then
	require.Len(t, componentList, 1)

	component := componentList[0]
	def := component.Definition()

	require.NotNil(t, def, "nil definition")
	assert.Equal(t, "anyPointerComponent", def.Name())
}
