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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStandardContainer_RegisterDefinition(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(container Container)
		definition   *Definition

		wantErr error
	}{
		{
			name:       "nil definition",
			definition: nil,
			wantErr:    errors.New("nil definition"),
		},
		{
			name: "already registered",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(&Definition{
					name: "anyDefinitionName",
				})
			},
			definition: &Definition{
				name: "anyDefinitionName",
			},
			wantErr: errors.New("register definition \"anyDefinitionName\": duplicate definition"),
		},
		{
			name: "valid definition",
			definition: &Definition{
				name: "anyDefinitionName",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			err := container.RegisterDefinition(tc.definition)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStandardContainer_UnregisterDefinition(t *testing.T) {
	testCases := []struct {
		name           string
		preCondition   func(container Container)
		definitionName string

		wantErr error
	}{
		{
			name:           "no definition",
			definitionName: "anyDefinitionName",
			wantErr:        errors.New("unregister definition \"anyDefinitionName\": definition not found"),
		},
		{
			name: "valid definition",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(&Definition{
					name: "anyDefinitionName",
				})
			},
			definitionName: "anyDefinitionName",
			wantErr:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			err := container.UnregisterDefinition(tc.definitionName)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStandardContainer_Definition(t *testing.T) {
	anyDefinition := &Definition{
		name: "anyDefinitionName",
	}

	testCases := []struct {
		name           string
		preCondition   func(container Container)
		definitionName string

		wantResult     bool
		wantDefinition *Definition
	}{
		{
			name:           "no definition",
			definitionName: "anyDefinitionName",
			wantResult:     false,
		},
		{
			name: "valid definition",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(anyDefinition)
			},
			definitionName: "anyDefinitionName",
			wantResult:     true,
			wantDefinition: anyDefinition,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			def, exists := container.Definition(tc.definitionName)

			// then
			require.Equal(t, tc.wantResult, exists)
			require.Equal(t, tc.wantDefinition, def)
		})
	}
}

func TestStandardContainer_ContainsDefinition(t *testing.T) {
	testCases := []struct {
		name           string
		preCondition   func(container Container)
		definitionName string

		wantResult bool
	}{
		{
			name:           "no definition",
			definitionName: "anyDefinitionName",
			wantResult:     false,
		},
		{
			name: "valid definition",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(&Definition{
					name: "anyDefinitionName",
				})
			},
			definitionName: "anyDefinitionName",
			wantResult:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			exists := container.ContainsDefinition(tc.definitionName)

			// then
			require.Equal(t, tc.wantResult, exists)
		})
	}
}

func TestStandardContainer_Definitions(t *testing.T) {
	anyDefinition := &Definition{
		name: "anyDefinitionName",
	}

	anotherDefinition := &Definition{
		name: "anotherDefinitionName",
	}

	testCases := []struct {
		name         string
		preCondition func(container Container)

		wantDefinitions []*Definition
	}{
		{
			name:            "no definition",
			wantDefinitions: []*Definition{},
		},
		{
			name: "return definitions",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(anyDefinition)

				_ = container.RegisterDefinition(anotherDefinition)
			},
			wantDefinitions: []*Definition{anyDefinition, anotherDefinition},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			definitions := container.Definitions()

			// then
			assert.Len(t, definitions, len(tc.wantDefinitions))

			for _, wantDefinition := range tc.wantDefinitions {
				assert.Contains(t, definitions, wantDefinition)
			}
		})
	}
}

func TestStandardContainer_DefinitionsOf(t *testing.T) {
	anyDefinition, _ := MakeDefinition(NewAnyPointerComponent)
	anotherDefinition, _ := MakeDefinition(NewAnyDependentComponent)

	testCases := []struct {
		name         string
		preCondition func(container Container)
		typ          reflect.Type

		wantDefinitions []*Definition
		wantPanic       error
	}{
		{
			name:      "nil definition type",
			typ:       nil,
			wantPanic: errors.New("nil definition type"),
		},
		{
			name:            "no definition",
			typ:             reflect.TypeFor[*AnyPointerComponent](),
			wantDefinitions: []*Definition{},
		},
		{
			name: "definitions by struct type",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(anyDefinition)

				_ = container.RegisterDefinition(anotherDefinition)
			},
			typ:             reflect.TypeFor[*AnyPointerComponent](),
			wantDefinitions: []*Definition{anyDefinition},
		},
		{
			name: "definitions by interface type",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(anyDefinition)

				_ = container.RegisterDefinition(anotherDefinition)
			},
			typ:             reflect.TypeFor[AnyComponent](),
			wantDefinitions: []*Definition{anyDefinition, anotherDefinition},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					// when
					container.DefinitionsOf(tc.typ)

					// then
				})
				return
			}

			assert.NotPanics(t, func() {
				// when
				definitions := container.DefinitionsOf(tc.typ)

				// then
				assert.Len(t, definitions, len(tc.wantDefinitions))

				for _, wantDefinition := range tc.wantDefinitions {
					require.Contains(t, definitions, wantDefinition)
				}
			})

		})
	}
}

func TestStandardContainer_DefinitionNames(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(container Container)

		wantNames []string
	}{
		{
			name:      "no definition",
			wantNames: []string{},
		},
		{
			name: "return definition names",
			preCondition: func(container Container) {
				err := container.RegisterDefinition(&Definition{
					name: "anyDefinitionName",
				})
				require.NoError(t, err)

				err = container.RegisterDefinition(&Definition{
					name: "anotherDefinitionName",
				})
				require.NoError(t, err)
			},
			wantNames: []string{"anyDefinitionName", "anotherDefinitionName"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			names := container.DefinitionNames()

			// then
			assert.ElementsMatch(t, tc.wantNames, names)
		})
	}
}

func TestStandardContainer_DefinitionNamesOf(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(container Container)
		typ          reflect.Type

		wantNames []string
		wantPanic error
	}{
		{
			name:      "nil definition type",
			typ:       nil,
			wantPanic: errors.New("nil definition type"),
		},
		{
			name:      "no definition",
			typ:       reflect.TypeFor[*AnyPointerComponent](),
			wantNames: []string{},
		},
		{
			name: "definition names by struct type",
			preCondition: func(container Container) {
				def, err := MakeDefinition(NewAnyPointerComponent, WithName("anyDefinitionName"))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				def, err = MakeDefinition(NewAnyDependentComponent, WithName("anyDependentName"))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			typ:       reflect.TypeFor[*AnyPointerComponent](),
			wantNames: []string{"anyDefinitionName"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					_ = container.DefinitionNamesOf(tc.typ)
				})
				return
			}

			require.NotPanics(t, func() {
				// when
				names := container.DefinitionNamesOf(tc.typ)

				// then
				assert.ElementsMatch(t, tc.wantNames, names)
			})
		})
	}
}

func TestStandardContainer_RegisterSingleton(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(container Container)
		instanceName string
		instance     any

		wantErr error
	}{
		{
			name:         "empty name",
			instanceName: "",
			wantErr:      errors.New("empty instance name"),
		},
		{
			name:         "nil instance",
			instanceName: "anyInstanceName",
			instance:     nil,
			wantErr:      errors.New("nil instance"),
		},
		{
			name: "already registered",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyPointerComponent{})
			},
			instanceName: "anyInstanceName",
			instance:     AnyPointerComponent{},
			wantErr:      errors.New("register singleton \"anyInstanceName\": duplicate instance"),
		},
		{
			name:         "valid singleton",
			instanceName: "anyInstanceName",
			instance:     &AnyPointerComponent{},
			wantErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			err := container.RegisterSingleton(tc.instanceName, tc.instance)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStandardContainer_ContainsSingleton(t *testing.T) {
	testCases := []struct {
		name          string
		preCondition  func(container Container)
		singletonName string

		wantResult bool
	}{
		{
			name:          "no singleton",
			singletonName: "anySingletonName",
			wantResult:    false,
		},
		{
			name: "valid singleton",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anySingletonName", AnyPointerComponent{})
			},
			singletonName: "anySingletonName",
			wantResult:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			exists := container.ContainsSingleton(tc.singletonName)

			// then
			require.Equal(t, tc.wantResult, exists)
		})
	}
}

func TestStandardContainer_Singleton(t *testing.T) {
	anySingleton := &AnyPointerComponent{}

	testCases := []struct {
		name          string
		preCondition  func(container Container)
		singletonName string

		wantResult    bool
		wantSingleton any
	}{
		{
			name:          "no singleton",
			singletonName: "anySingletonName",
			wantResult:    false,
		},
		{
			name: "valid singleton",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anySingletonName", anySingleton)
			},
			singletonName: "anySingletonName",
			wantResult:    true,
			wantSingleton: anySingleton,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			singleton, exists := container.Singleton(tc.singletonName)

			// then
			require.Equal(t, tc.wantResult, exists)
			require.Equal(t, tc.wantSingleton, singleton)
		})
	}

}

func TestStandardContainer_RemoveSingleton(t *testing.T) {
	testCases := []struct {
		name          string
		preCondition  func(container Container)
		singletonName string

		wantErr error
	}{
		{
			name:          "no singleton",
			singletonName: "anySingletonName",
			wantErr:       errors.New("remove singleton \"anySingletonName\": not found"),
		},
		{
			name: "valid singleton",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anySingletonName", AnyPointerComponent{})
			},
			singletonName: "anySingletonName",
			wantErr:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			err := container.RemoveSingleton(tc.singletonName)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStandardContainer_CanResolve(t *testing.T) {
	var testCases = []struct {
		name         string
		preCondition func(parent, container Container)
		instanceName string

		wantResult bool
	}{
		{
			name:         "cannot resolve empty instance name",
			instanceName: "",
			wantResult:   false,
		},
		{
			name: "can resolve instance singleton",
			preCondition: func(parent, container Container) {
				err := container.RegisterSingleton("anyInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)
			},
			instanceName: "anyInstanceName",
			wantResult:   true,
		},
		{
			name: "can resolve instance definition",
			preCondition: func(parent, container Container) {
				_ = container.RegisterDefinition(&Definition{
					name: "anyInstanceName",
				})
			},
			instanceName: "anyInstanceName",
			wantResult:   true,
		},
		{
			name: "can resolve instance from parent container",
			preCondition: func(parent, container Container) {
				err := parent.RegisterSingleton("anyInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)
			},
			instanceName: "anyInstanceName",
			wantResult:   true,
		},
		{
			name:         "cannot resolve instance",
			instanceName: "anyInstanceName",
			wantResult:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			parentContainer := NewStandardContainer()
			container := NewStandardContainer()
			container.SetParentContainer(parentContainer)

			if tc.preCondition != nil {
				tc.preCondition(parentContainer, container)
			}

			// when
			result := container.CanResolve(tc.instanceName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestStandardContainer_CanResolveType(t *testing.T) {
	var testCases = []struct {
		name         string
		preCondition func(parent, container Container)
		instanceType reflect.Type

		wantResult bool
	}{
		{
			name:         "cannot resolve nil instance type",
			instanceType: nil,
			wantResult:   false,
		},
		{
			name: "can resolve instance singleton",
			preCondition: func(parent, container Container) {
				_ = container.RegisterSingleton("anyInstanceName", &AnyPointerComponent{})
			},
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantResult:   true,
		},
		{
			name: "can resolve instance definition",
			preCondition: func(parent, container Container) {
				constructor, _ := createConstructor(NewAnyPointerComponent)
				_ = container.RegisterDefinition(&Definition{
					name:        "anyInstanceName",
					constructor: constructor,
				})
			},
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantResult:   true,
		},
		{
			name: "can resolve instance from parent container",
			preCondition: func(parent, container Container) {
				err := parent.RegisterSingleton("anyInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)
			},
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantResult:   true,
		},
		{
			name:         "cannot resolve instance",
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantResult:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			parentContainer := NewStandardContainer()
			container := NewStandardContainer()
			container.SetParentContainer(parentContainer)

			if tc.preCondition != nil {
				tc.preCondition(parentContainer, container)
			}

			// when
			result := container.CanResolveType(tc.instanceType)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestStandardContainer_Resolve(t *testing.T) {
	var testCases = []struct {
		name         string
		ctx          context.Context
		preCondition func(parent, container Container)
		instanceName string

		wantErr error
		wantTyp reflect.Type
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: errors.New("nil context"),
			wantTyp: nil,
		},
		{
			ctx:          context.Background(),
			name:         "empty name",
			instanceName: "",
			wantErr:      errors.New("empty instance name"),
		},
		{
			name: "resolve singleton already in container",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyPointerComponent{})
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[AnyPointerComponent](),
		},
		{
			name: "resolve from singleton definition",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyPointerComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name:         "no singleton/definition",
			ctx:          context.Background(),
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyPointerComponent](),
			wantErr:      errors.New("resolve \"anyInstanceName\": not found"),
		},
		{
			name: "resolve from prototype definition",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyPointerComponent, WithName("anyInstanceName"), WithScope(PrototypeScope))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name: "no scope",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyPointerComponent, WithName("anyInstanceName"), WithScope("anyScope"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": scope \"anyScope\" not found"),
		},
		{
			name: "resolve from custom scope",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				scope := &AnyMockScope{
					useFactory: true,
				}
				_ = container.RegisterScope("anyScope", scope)

				def, _ := MakeDefinition(NewAnyPointerComponent, WithName("anyInstanceName"), WithScope("anyScope"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name: "scope error",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				scope := &AnyMockScope{
					useFactory: false,
				}

				scope.On("Resolve", mock.AnythingOfType("*context.valueCtx"), "anyInstanceName", mock.Anything).
					Return(nil, errors.New("custom scope error"))
				_ = container.RegisterScope("anyScope", scope)

				def, _ := MakeDefinition(NewAnyPointerComponent, WithName("anyInstanceName"), WithScope("anyScope"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": scope \"anyScope\": custom scope error"),
		},
		{
			name: "resolve singleton with dependencies",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyDependentComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnySimpleComponent, WithName("anyDependentName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyDependentComponent](),
		},
		{
			name: "cannot resolve singleton with dependencies",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyDependentComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": create \"anyInstanceName\" (*component.AnyDependentComponent): unsatisfied dependency for argument 0 (component.AnySimpleComponent): resolve type component.AnySimpleComponent: not found"),
		},
		{
			name: "resolve singleton with named dependencies",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyDependentComponent, WithName("anyInstanceName"),
					WithQualifierFor[AnySimpleComponent]("anyDependentName"))
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnySimpleComponent, WithName("anyDependentName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyDependentComponent](),
		},
		{
			name: "cannot resolve singleton with named dependencies",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyDependentComponent, WithName("anyInstanceName"),
					WithQualifierFor[AnySimpleComponent]("anyDependentName"))
				_ = container.RegisterDefinition(def)

			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": create \"anyInstanceName\" (*component.AnyDependentComponent): unsatisfied dependency for argument 0 \"anyDependentName\" (component.AnySimpleComponent): resolve \"anyDependentName\": not found"),
		},
		{
			name: "resolve singleton with slice dependencies",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyIndexedComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnyPointerComponent)
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnyDependentComponent)
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnySimpleComponent)
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyIndexedComponent](),
		},
		{
			name: "resolve instance from parent container",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				err := parent.RegisterSingleton("anyInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name: "cannot resolve singleton with slice dependencies",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(NewAnyIndexedComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnyPointerComponent)
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnyDependentComponent)
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": create \"anyInstanceName\" (*component.AnyIndexedComponent): unsatisfied dependency for argument 0 ([]component.AnyComponent): resolve \"anyDependentComponent\": create \"anyDependentComponent\" (*component.AnyDependentComponent): unsatisfied dependency for argument 0 (component.AnySimpleComponent): resolve type component.AnySimpleComponent: not found"),
		},
		{
			name: "constructor error",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, _ := MakeDefinition(func() *AnySimpleComponent {
					panic("any constructor error")
					return nil
				}, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": invoke constructor \"anyInstanceName\" (*component.AnySimpleComponent): constructor panic: any constructor error"),
		},
		{
			name: "resolve instance with circular dependencies",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, err := MakeDefinition(func(dep *AnyDependentComponent) *AnyPointerComponent {
					return &AnyPointerComponent{}
				}, WithName("anyInstanceName"))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				def, err = MakeDefinition(func(dep *AnyPointerComponent) *AnyDependentComponent {
					return &AnyDependentComponent{}
				}, WithName("anyDependentName"))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": create \"anyInstanceName\" (*component.AnyPointerComponent): unsatisfied dependency for argument 0 (*component.AnyDependentComponent): resolve \"anyDependentName\": create \"anyDependentName\" (*component.AnyDependentComponent): unsatisfied dependency for argument 0 (*component.AnyPointerComponent): resolve \"anyInstanceName\": circular dependency detected"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			parentContainer := NewStandardContainer()
			container := NewStandardContainer()
			container.SetParentContainer(parentContainer)

			if tc.preCondition != nil {
				tc.preCondition(parentContainer, container)
			}

			// when
			result, err := container.Resolve(tc.ctx, tc.instanceName)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantTyp, reflect.TypeOf(result))
		})
	}
}

func TestStandardContainer_ResolveType(t *testing.T) {
	var testCases = []struct {
		name         string
		ctx          context.Context
		preCondition func(container Container)
		instanceType reflect.Type

		wantErr error
		wantTyp reflect.Type
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: errors.New("nil context"),
			wantTyp: nil,
		},
		{
			name:    "nil type",
			ctx:     context.Background(),
			wantErr: errors.New("nil instance type"),
		},
		{
			name: "multiple singletons",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyPointerComponent{})
				_ = container.RegisterSingleton("anotherInstanceName", AnyPointerComponent{})
			},
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantErr:      errors.New("resolve type component.AnyPointerComponent: ambiguous match"),
		},
		{
			name: "multi definitions",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyPointerComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnyPointerComponent, WithName("anotherInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceType: reflect.TypeFor[*AnyPointerComponent](),
			wantErr:      errors.New("resolve type *component.AnyPointerComponent: ambiguous match"),
		},
		{
			name:         "no singleton/definition",
			ctx:          context.Background(),
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantErr:      errors.New("resolve type component.AnyPointerComponent: not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			result, err := container.ResolveType(tc.ctx, tc.instanceType)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantTyp, reflect.TypeOf(result))
		})
	}
}

func TestStandardContainer_ResolveAs(t *testing.T) {
	var testCases = []struct {
		name         string
		ctx          context.Context
		preCondition func(container Container)
		instanceName string
		instanceType reflect.Type

		wantErr error
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: errors.New("nil context"),
		},
		{
			name:         "empty name",
			ctx:          context.Background(),
			instanceName: "",
			wantErr:      errors.New("empty instance name"),
		},
		{
			name:         "nil type",
			ctx:          context.Background(),
			instanceName: "anyInstanceName",
			instanceType: nil,
			wantErr:      errors.New("nil instance type"),
		},
		{
			name: "assignable type",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyPointerComponent{})
			},
			instanceName: "anyInstanceName",
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantErr:      nil,
		},
		{
			name: "not assignable type",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyPointerComponent{})
			},
			instanceName: "anyInstanceName",
			instanceType: reflect.TypeFor[AnyDependentComponent](),
			wantErr:      errors.New("resolve \"anyInstanceName\": component.AnyPointerComponent is not convertible to component.AnyDependentComponent: type mismatch"),
		},
		{
			name: "resolve error",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(func() *AnySimpleComponent {
					panic("any constructor error")
					return nil
				}, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			instanceType: reflect.TypeFor[*AnySimpleComponent](),
			wantErr:      errors.New("resolve \"anyInstanceName\": invoke constructor \"anyInstanceName\" (*component.AnySimpleComponent): constructor panic: any constructor error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()
			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			result, err := container.ResolveAs(tc.ctx, tc.instanceName, tc.instanceType)
			// when

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestStandardContainer_ResolveAll(t *testing.T) {
	var testCases = []struct {
		name         string
		ctx          context.Context
		preCondition func(parent, container Container)
		instanceType reflect.Type

		wantLen   int
		wantTypes []reflect.Type
		wantErr   error
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: errors.New("nil context"),
		},
		{
			name:         "nil type",
			ctx:          context.Background(),
			instanceType: nil,
			wantErr:      errors.New("nil instance type"),
		},
		{
			name:         "no singleton/definition",
			ctx:          context.Background(),
			instanceType: reflect.TypeFor[AnyPointerComponent](),
			wantLen:      0,
			wantTypes:    []reflect.Type{},
			wantErr:      nil,
		},
		{
			name: "resolve singletons with definition",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, err := MakeDefinition(NewAnySimpleComponent, WithName("anyInstanceName"))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				_, err = container.Resolve(context.Background(), "anyInstanceName")
				require.NoError(t, err)
			},
			instanceType: reflect.TypeFor[AnySimpleComponent](),
			wantLen:      1,
			wantTypes:    []reflect.Type{reflect.TypeFor[AnySimpleComponent]()},
		},
		{
			name: "resolve singleton without definition",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				err := container.RegisterSingleton("anyInstanceName", &AnySimpleComponent{})
				require.NoError(t, err)

				_, err = container.Resolve(context.Background(), "anyInstanceName")
				require.NoError(t, err)
			},
			instanceType: reflect.TypeFor[AnySimpleComponent](),
			wantLen:      1,
			wantTypes:    []reflect.Type{reflect.TypeFor[*AnySimpleComponent]()},
		},
		{
			name: "resolve instances by interface type",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, err := MakeDefinition(NewAnySimpleComponent, WithName("anyInstanceName"))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				err = container.RegisterSingleton("anotherInstanceName", &AnyDependentComponent{})
				require.NoError(t, err)

				err = container.RegisterSingleton("anyInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)

				_, err = container.Resolve(context.Background(), "anyInstanceName")
				require.NoError(t, err)
			},
			instanceType: reflect.TypeFor[AnyComponent](),
			wantLen:      2,
			wantTypes:    []reflect.Type{reflect.TypeFor[*AnyDependentComponent](), reflect.TypeFor[*AnyPointerComponent]()},
		},
		{
			name: "resolve singleton/prototype instances",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, err := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"), WithScope(PrototypeScope))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				err = container.RegisterSingleton("anotherInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)
			},
			instanceType: reflect.TypeFor[AnyComponent](),
			wantLen:      2,
			wantTypes:    []reflect.Type{reflect.TypeFor[*AnyInitializableComponent](), reflect.TypeFor[*AnyPointerComponent]()},
		},
		{
			name: "resolve instances from parent container",
			ctx:  context.Background(),
			preCondition: func(parent, container Container) {
				def, err := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"), WithScope(PrototypeScope))
				require.NoError(t, err)

				err = parent.RegisterDefinition(def)
				require.NoError(t, err)

				err = parent.RegisterSingleton("anotherInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)
			},
			instanceType: reflect.TypeFor[AnyComponent](),
			wantLen:      2,
			wantTypes:    []reflect.Type{reflect.TypeFor[*AnyInitializableComponent](), reflect.TypeFor[*AnyPointerComponent]()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			parentContainer := NewStandardContainer()
			container := NewStandardContainer()
			container.SetParentContainer(parentContainer)

			if tc.preCondition != nil {
				tc.preCondition(parentContainer, container)
			}

			// when
			results, err := container.ResolveAll(tc.ctx, tc.instanceType)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tc.wantLen)

			gotTypes := make([]reflect.Type, len(results))
			for i, r := range results {
				gotTypes[i] = reflect.TypeOf(r)
			}
			require.ElementsMatch(t, tc.wantTypes, gotTypes)
		})
	}
}

func TestStandardContainer_ParentContainer(t *testing.T) {
	parentContainer := NewStandardContainer()

	testCases := []struct {
		name         string
		preCondition func(container *StandardContainer)
		wantParent   Container
	}{
		{
			name: "no parent container",
			preCondition: func(container *StandardContainer) {
				// no parent container registered
			},
			wantParent: nil,
		},
		{
			name: "parent container registered",
			preCondition: func(container *StandardContainer) {
				container.SetParentContainer(parentContainer)
			},
			wantParent: parentContainer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			parent := container.ParentContainer()

			// then
			require.Equal(t, tc.wantParent, parent)
		})
	}
}

func TestStandardContainer_SetParentContainer(t *testing.T) {
	parentContainer := NewStandardContainer()

	testCases := []struct {
		name         string
		preCondition func(container *StandardContainer)
		parent       Container

		wantPanic error
	}{
		{
			name:      "set nil parent container",
			parent:    nil,
			wantPanic: errors.New("nil parent container"),
		},
		{
			name:      "set valid parent container",
			parent:    parentContainer,
			wantPanic: nil,
		},
		{
			name: "set parent container when already associated with the same parent container",
			preCondition: func(container *StandardContainer) {
				container.SetParentContainer(parentContainer)
			},
			parent:    parentContainer,
			wantPanic: nil,
		},
		{
			name: "set parent container when already associated with a parent container",
			preCondition: func(container *StandardContainer) {
				container.SetParentContainer(parentContainer)
			},
			parent:    NewStandardContainer(),
			wantPanic: errors.New("already associated with a parent container"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					container.SetParentContainer(tc.parent)
				})
				return
			}

			assert.NotPanics(t, func() {
				container.SetParentContainer(tc.parent)
			})
		})
	}
}

func TestStandardContainer_RegisterResolvable(t *testing.T) {
	rAnyType := reflect.TypeFor[*AnyPointerComponent]()
	anyComponent := &AnyPointerComponent{}

	testCases := []struct {
		name     string
		typ      reflect.Type
		instance any

		wantErr         error
		wantResolvables map[reflect.Type]any
	}{
		{
			name:    "nil type",
			typ:     nil,
			wantErr: errors.New("nil instance type"),
		},
		{
			name:    "nil value",
			typ:     rAnyType,
			wantErr: errors.New("nil instance"),
		},
		{
			name:     "valid resolvable",
			typ:      rAnyType,
			instance: anyComponent,
			wantErr:  nil,
			wantResolvables: map[reflect.Type]any{
				rAnyType: anyComponent,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			// when
			err := container.RegisterDependency(tc.typ, tc.instance)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			for wantResolvableType, wantResolvable := range tc.wantResolvables {
				resolvable, ok := container.resolvableDependencies[wantResolvableType]
				require.True(t, ok)
				assert.Equal(t, wantResolvable, resolvable)
			}
		})
	}
}

func TestStandardContainer_RegisterScope(t *testing.T) {
	testCases := []struct {
		name      string
		scopeName string
		scope     Scope

		wantErr error
	}{
		{
			name:    "empty scope name",
			wantErr: errors.New("empty scope name"),
		},
		{
			name:      "nil scope",
			scopeName: "anyScopeName",
			scope:     nil,
			wantErr:   errors.New("nil scope"),
		},
		{
			name:      "singleton scope replacement not allowed",
			scopeName: SingletonScope,
			scope:     &AnyMockScope{},
			wantErr:   errors.New("register scope \"singleton\": reserved scope"),
		},
		{
			name:      "prototype scope replacement not allowed",
			scopeName: PrototypeScope,
			scope:     &AnyMockScope{},
			wantErr:   errors.New("register scope \"prototype\": reserved scope"),
		},
		{
			name:      "valid scope",
			scopeName: "anyScopeName",
			scope:     &AnyMockScope{},
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			// when
			err := container.RegisterScope(tc.scopeName, tc.scope)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStandardContainer_DestroySingletons(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(container Container)

		wantErr error
	}{
		{
			name: "destroy singletons with no singletons",
			preCondition: func(container Container) {
				// no singletons registered
			},
			wantErr: nil,
		},
		{
			name: "destroy singletons",
			preCondition: func(container Container) {
				componentList := []struct {
					name string
					fn   ConstructorFunc
					typ  reflect.Type
				}{
					{
						name: "disposable",
						fn:   NewAnyDisposableComponent,
					},
					{
						name: "simple",
						fn:   NewAnySimpleComponent,
					},
					{
						name: "dependent",
						fn:   NewAnyDependentComponent,
					},
				}

				for _, c := range componentList {
					def, err := MakeDefinition(c.fn, WithName(c.name))
					require.NoError(t, err)

					err = container.RegisterDefinition(def)
					require.NoError(t, err)
				}

				for _, c := range componentList {
					_, err := container.Resolve(context.Background(), c.name)
					require.NoError(t, err)
				}
			},
			wantErr: nil,
		},
		{
			name: "dispose error",
			preCondition: func(container Container) {
				def, err := MakeDefinition(func() *AnyDisposableComponent {
					return &AnyDisposableComponent{
						disposeError: errors.New("failed to dispose"),
					}
				}, WithName("disposable"))
				require.NoError(t, err)

				err = container.RegisterDefinition(def)
				require.NoError(t, err)

				_, err = container.Resolve(context.Background(), "disposable")
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			container.DestroySingletons()

			// then
		})
	}
}

func TestStandardContainer_SingletonNames(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(container Container)

		wantNames []string
	}{
		{
			name: "no singletons",
			preCondition: func(container Container) {
				// no singletons registered
			},
			wantNames: []string{},
		},
		{
			name: "multiple singletons",
			preCondition: func(container Container) {
				err := container.RegisterSingleton("anyInstanceName", &AnyPointerComponent{})
				require.NoError(t, err)

				err = container.RegisterSingleton("anotherInstanceName", &AnySimpleComponent{})
				require.NoError(t, err)
			},
			wantNames: []string{"anyInstanceName", "anotherInstanceName"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			names := container.SingletonNames()

			// then
			assert.ElementsMatch(t, tc.wantNames, names)
		})
	}
}

func TestStandardContainer_Scope(t *testing.T) {
	anyScope := &AnyMockScope{}

	testCases := []struct {
		name         string
		preCondition func(container Container)
		scopeName    string

		wantResult bool
		wantScope  Scope
	}{
		{
			name:       "empty scope name",
			scopeName:  "",
			wantResult: false,
			wantScope:  nil,
		},
		{
			name:       "no scope",
			scopeName:  "anyScopeName",
			wantResult: false,
		},
		{
			name: "valid scope",
			preCondition: func(container Container) {
				_ = container.RegisterScope("anyScopeName", anyScope)
			},
			scopeName:  "anyScopeName",
			wantResult: true,
			wantScope:  anyScope,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			scope, exists := container.Scope(tc.scopeName)

			// then
			require.Equal(t, tc.wantResult, exists)
			require.Equal(t, tc.wantScope, scope)
		})
	}
}

func TestStandardContainer_UseBeforeInitProcessor(t *testing.T) {
	testCases := []struct {
		name                string
		beforeInitProcessor BeforeInitProcessor

		wantErr error
		wantLen int
	}{
		{
			name:    "nil before-init processor",
			wantErr: errors.New("nil before-init processor"),
			wantLen: 0,
		},
		{
			name:                "valid before-init processor",
			beforeInitProcessor: &AnyMockBeforeInitProcessor{},
			wantErr:             nil,
			wantLen:             1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			// when
			err := container.UseBeforeInitProcessor(tc.beforeInitProcessor)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			assert.Len(t, container.beforeInitProcessors, tc.wantLen)
		})
	}
}

func TestStandardContainer_UseAfterInitProcessor(t *testing.T) {
	testCases := []struct {
		name               string
		afterInitProcessor AfterInitProcessor

		wantErr error
		wantLen int
	}{
		{
			name:    "nil after-init processor",
			wantErr: errors.New("nil after-init processor"),
			wantLen: 0,
		},
		{
			name:               "valid after-init processor",
			afterInitProcessor: &AnyMockAfterInitProcessor{},
			wantErr:            nil,
			wantLen:            1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			// when
			err := container.UseAfterInitProcessor(tc.afterInitProcessor)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			assert.Len(t, container.afterInitProcessors, tc.wantLen)
		})
	}
}

func TestStandardContainer_Initialize(t *testing.T) {

	testCases := []struct {
		name         string
		ctx          context.Context
		preCondition func(container Container)
		instanceName string

		wantErr error
		wantTyp reflect.Type
	}{
		{
			name: "initialize error",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(func() *AnyInitializableComponent {
					return &AnyInitializableComponent{
						initError: errors.New("failed to initialize"),
					}
				}, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": initialize \"anyInstanceName\" (*component.AnyInitializableComponent): invoke init: failed to initialize"),
		},
		{
			name: "no initialize error",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyInitializableComponent](),
		},
		{
			name: "before-init processor returns error",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				anyBeforeInitProcessor := &AnyMockBeforeInitProcessor{}
				_ = container.UseBeforeInitProcessor(anyBeforeInitProcessor)
				anyBeforeInitProcessor.On("ProcessBeforeInit", mock.AnythingOfType("*context.valueCtx"), "anyInstanceName", mock.AnythingOfType("*component.AnyInitializableComponent")).
					Return(nil, errors.New("before initialization error"))
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": initialize \"anyInstanceName\" (*component.AnyInitializableComponent): apply before-init processors: before-init processor (*component.AnyMockBeforeInitProcessor): before initialization error"),
		},
		{
			name: "before-init processor returns nil value",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				anyBeforeInitProcessor := &AnyMockBeforeInitProcessor{}
				_ = container.UseBeforeInitProcessor(anyBeforeInitProcessor)
				anyBeforeInitProcessor.On("ProcessBeforeInit", mock.AnythingOfType("*context.valueCtx"), "anyInstanceName", mock.AnythingOfType("*component.AnyInitializableComponent")).
					Return(nil, nil)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": initialize \"anyInstanceName\" (*component.AnyInitializableComponent): apply before-init processors: before-init processor (*component.AnyMockBeforeInitProcessor) returned nil"),
		},
		{
			name: "apply before-init processor",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				anyBeforeInitProcessor := &AnyMockBeforeInitProcessor{}
				_ = container.UseBeforeInitProcessor(anyBeforeInitProcessor)
				anyBeforeInitProcessor.On("ProcessBeforeInit", mock.AnythingOfType("*context.valueCtx"), "anyInstanceName", mock.AnythingOfType("*component.AnyInitializableComponent")).
					Return(&AnyPointerComponent{}, nil)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyPointerComponent](),
		},
		{
			name: "after-init processor returns error",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				anyAfterInitProcessor := &AnyMockAfterInitProcessor{}
				_ = container.UseAfterInitProcessor(anyAfterInitProcessor)
				anyAfterInitProcessor.On("ProcessAfterInit", mock.AnythingOfType("*context.valueCtx"), "anyInstanceName", mock.AnythingOfType("*component.AnyInitializableComponent")).
					Return(nil, errors.New("after initialization error"))
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": initialize \"anyInstanceName\" (*component.AnyInitializableComponent): apply after-init processors: after-init processor (*component.AnyMockAfterInitProcessor): after initialization error"),
		},
		{
			name: "after-init processor returns nil value",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				anyAfterInitProcessor := &AnyMockAfterInitProcessor{}
				_ = container.UseAfterInitProcessor(anyAfterInitProcessor)
				anyAfterInitProcessor.On("ProcessAfterInit", mock.AnythingOfType("*context.valueCtx"), "anyInstanceName", mock.AnythingOfType("*component.AnyInitializableComponent")).
					Return(nil, nil)
			},
			instanceName: "anyInstanceName",
			wantErr:      errors.New("resolve \"anyInstanceName\": initialize \"anyInstanceName\" (*component.AnyInitializableComponent): apply after-init processors: after-init processor (*component.AnyMockAfterInitProcessor) returned nil"),
		},
		{
			name: "apply after-init processor",
			ctx:  context.Background(),
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyInitializableComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				anyAfterInitProcessor := &AnyMockAfterInitProcessor{}
				_ = container.UseAfterInitProcessor(anyAfterInitProcessor)
				anyAfterInitProcessor.On("ProcessAfterInit", mock.AnythingOfType("*context.valueCtx"), "anyInstanceName", mock.AnythingOfType("*component.AnyInitializableComponent")).
					Return(&AnyPointerComponent{}, nil)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyPointerComponent](),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewStandardContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			result, err := container.Resolve(tc.ctx, tc.instanceName)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantTyp, reflect.TypeOf(result))
		})
	}
}
