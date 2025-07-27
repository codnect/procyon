package component

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type AnyScope struct {
	mock.Mock

	useFactory bool
}

func (a *AnyScope) Resolve(ctx context.Context, name string, fn FactoryFunc) (any, error) {
	if a.useFactory {
		return fn(ctx)
	}

	result := a.Called(ctx, name, fn)
	return result.Get(0).(any), result.Error(1)
}

func (a *AnyScope) Remove(ctx context.Context, name string) error {
	result := a.Called(ctx, name)
	return result.Error(0)
}

type AnyPreProcessor struct {
	mock.Mock
}

func (a *AnyPreProcessor) ProcessBeforeInit(ctx context.Context, instance any) (any, error) {
	result := a.Called(ctx, instance)
	return result.Get(0).(any), result.Error(1)
}

type AnyPostProcessor struct {
	mock.Mock
}

func (a *AnyPostProcessor) ProcessAfterInit(ctx context.Context, instance any) (any, error) {
	result := a.Called(ctx, instance)
	return result.Get(0).(any), result.Error(1)
}

func TestDefaultContainer_RegisterDefinition(t *testing.T) {
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
			wantErr: ErrDefinitionAlreadyExists,
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
			container := NewDefaultContainer()

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

func TestDefaultContainer_UnregisterDefinition(t *testing.T) {
	testCases := []struct {
		name           string
		preCondition   func(container Container)
		definitionName string

		wantErr error
	}{
		{
			name:           "no definition",
			definitionName: "anyDefinitionName",
			wantErr:        ErrDefinitionNotFound,
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
			container := NewDefaultContainer()

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

func TestDefaultContainer_Definition(t *testing.T) {
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
			container := NewDefaultContainer()

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

func TestDefaultContainer_ContainsDefinition(t *testing.T) {
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
			container := NewDefaultContainer()

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

func TestDefaultContainer_Definitions(t *testing.T) {

}

func TestDefaultContainer_DefinitionsOf(t *testing.T) {

}

func TestDefaultContainer_RegisterSingleton(t *testing.T) {
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
			wantErr:      errors.New("empty name"),
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
				_ = container.RegisterSingleton("anyInstanceName", AnyComponent{})
			},
			instanceName: "anyInstanceName",
			instance:     AnyComponent{},
			wantErr:      ErrInstanceAlreadyExists,
		},
		{
			name:         "valid singleton",
			instanceName: "anyInstanceName",
			instance:     &AnyComponent{},
			wantErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

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

func TestDefaultContainer_ContainsSingleton(t *testing.T) {
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
				_ = container.RegisterSingleton("anySingletonName", AnyComponent{})
			},
			singletonName: "anySingletonName",
			wantResult:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

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

func TestDefaultContainer_Singleton(t *testing.T) {
	anySingleton := &AnyComponent{}

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
			container := NewDefaultContainer()

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

func TestDefaultContainer_RemoveSingleton(t *testing.T) {
	testCases := []struct {
		name          string
		preCondition  func(container Container)
		singletonName string

		wantErr error
	}{
		{
			name:          "no singleton",
			singletonName: "anySingletonName",
			wantErr:       ErrInstanceNotFound,
		},
		{
			name: "valid singleton",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anySingletonName", AnyComponent{})
			},
			singletonName: "anySingletonName",
			wantErr:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

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

func TestDefaultContainer_CanResolve(t *testing.T) {
	var testCases = []struct {
		name         string
		preCondition func(container Container)
		instanceName string

		wantResult bool
	}{
		{
			name: "can resolve instance singleton",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", &AnyComponent{})
			},
			instanceName: "anyInstanceName",
			wantResult:   true,
		},
		{
			name: "can resolve instance definition",
			preCondition: func(container Container) {
				_ = container.RegisterDefinition(&Definition{
					name: "anyInstanceName",
				})
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
			container := NewDefaultContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			result := container.CanResolve(tc.instanceName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestDefaultContainer_CanResolveType(t *testing.T) {
	var testCases = []struct {
		name         string
		preCondition func(container Container)
		instanceType reflect.Type

		wantResult bool
	}{
		{
			name: "can resolve instance singleton",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", &AnyComponent{})
			},
			instanceType: reflect.TypeFor[AnyComponent](),
			wantResult:   true,
		},
		{
			name: "can resolve instance definition",
			preCondition: func(container Container) {
				constructor, _ := createConstructor(NewAnyComponent)
				_ = container.RegisterDefinition(&Definition{
					name:        "anyInstanceName",
					constructor: constructor,
				})
			},
			instanceType: reflect.TypeFor[AnyComponent](),
			wantResult:   true,
		},
		{
			name:         "cannot resolve instance",
			instanceType: reflect.TypeFor[AnyComponent](),
			wantResult:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

			if tc.preCondition != nil {
				tc.preCondition(container)
			}

			// when
			result := container.CanResolveType(tc.instanceType)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestDefaultContainer_Resolve(t *testing.T) {
	var testCases = []struct {
		name         string
		ctx          context.Context
		preCondition func(container Container)
		instanceName string

		wantErr error
		wantTyp reflect.Type
	}{
		{
			name:         "empty name",
			instanceName: "",
			wantErr:      errors.New("empty name"),
		},
		{
			name: "resolve singleton already in container",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyComponent{})
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[AnyComponent](),
		},
		{
			name: "resolve from singleton definition",
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyComponent](),
		},
		{
			name:         "no singleton/definition",
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyComponent](),
			wantErr:      ErrDefinitionNotFound,
		},
		{
			name: "resolve from prototype definition",
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyComponent, WithName("anyInstanceName"), WithScope(PrototypeScope))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyComponent](),
		},
		{
			name: "no scope",
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyComponent, WithName("anyInstanceName"), WithScope("anyScope"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantErr:      ErrScopeNotFound,
		},
		{
			name: "resolve from custom scope",
			preCondition: func(container Container) {
				scope := &AnyScope{
					useFactory: true,
				}
				_ = container.RegisterScope("anyScope", scope)

				def, _ := MakeDefinition(NewAnyComponent, WithName("anyInstanceName"), WithScope("anyScope"))
				_ = container.RegisterDefinition(def)
			},
			instanceName: "anyInstanceName",
			wantTyp:      reflect.TypeFor[*AnyComponent](),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

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

func TestDefaultContainer_ResolveType(t *testing.T) {
	var testCases = []struct {
		name         string
		ctx          context.Context
		preCondition func(container Container)
		instanceType reflect.Type

		wantErr error
		wantTyp reflect.Type
	}{
		{
			name:    "nil type",
			wantErr: errors.New("nil type"),
		},
		{
			name: "multiple singletons",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyComponent{})
				_ = container.RegisterSingleton("anotherInstanceName", AnyComponent{})
			},
			instanceType: reflect.TypeFor[AnyComponent](),
			wantErr:      errors.New("multiple singletons found"),
		},
		{
			name: "multi definitions",
			preCondition: func(container Container) {
				def, _ := MakeDefinition(NewAnyComponent, WithName("anyInstanceName"))
				_ = container.RegisterDefinition(def)

				def, _ = MakeDefinition(NewAnyComponent, WithName("anotherInstanceName"))
				_ = container.RegisterDefinition(def)
			},
			instanceType: reflect.TypeFor[*AnyComponent](),
			wantErr:      errors.New("multiple definitions found"),
		},
		{
			name:         "no singleton/definition",
			instanceType: reflect.TypeFor[AnyComponent](),
			wantErr:      ErrInstanceNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

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

func TestDefaultContainer_ResolveAs(t *testing.T) {
	var testCases = []struct {
		name         string
		ctx          context.Context
		preCondition func(container Container)
		instanceName string
		instanceType reflect.Type

		wantErr error
	}{
		{
			name:         "empty name",
			wantErr:      errors.New("empty name"),
			instanceName: "",
		},
		{
			name:         "nil type",
			wantErr:      errors.New("nil type"),
			instanceName: "anyInstanceName",
			instanceType: nil,
		},
		{
			name: "assignable type",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyComponent{})
			},
			instanceName: "anyInstanceName",
			instanceType: reflect.TypeFor[AnyComponent](),
			wantErr:      nil,
		},
		{
			name: "not assignable type",
			preCondition: func(container Container) {
				_ = container.RegisterSingleton("anyInstanceName", AnyComponent{})
			},
			instanceName: "anyInstanceName",
			instanceType: reflect.TypeFor[AnotherComponent](),
			wantErr:      errors.New("component \"anyInstanceName\" is not assignable to component.AnotherComponent"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()
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

func TestDefaultContainer_ResolveAll(t *testing.T) {

}

func TestDefaultContainer_RegisterResolvable(t *testing.T) {

}

func TestDefaultContainer_RegisterScope(t *testing.T) {
	testCases := []struct {
		name      string
		scopeName string
		scope     Scope

		wantErr error
	}{
		{
			name:    "invalid scope name",
			wantErr: ErrInvalidScopeName,
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
			scope:     &AnyScope{},
			wantErr:   ErrScopeReplacementNotAllowed,
		},
		{
			name:      "prototype scope replacement not allowed",
			scopeName: PrototypeScope,
			scope:     &AnyScope{},
			wantErr:   ErrScopeReplacementNotAllowed,
		},
		{
			name:      "valid scope",
			scopeName: "anyScopeName",
			scope:     &AnyScope{},
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

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

func TestDefaultContainer_Scope(t *testing.T) {
	anyScope := &AnyScope{}

	testCases := []struct {
		name         string
		preCondition func(container Container)
		scopeName    string

		wantResult bool
		wantScope  Scope
	}{
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
			container := NewDefaultContainer()

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

func TestDefaultContainer_UsePreProcessor(t *testing.T) {
	testCases := []struct {
		name         string
		preProcessor PreProcessor

		wantErr error
		wantLen int
	}{
		{
			name:    "nil pre processor",
			wantErr: errors.New("nil processor"),
			wantLen: 0,
		},
		{
			name:         "valid pre processor",
			preProcessor: &AnyPreProcessor{},
			wantErr:      nil,
			wantLen:      1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

			// when
			err := container.UsePreProcessor(tc.preProcessor)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			assert.Len(t, container.preProcessors, tc.wantLen)
		})
	}
}

func TestDefaultContainer_UsePostProcessor(t *testing.T) {
	testCases := []struct {
		name          string
		postProcessor PostProcessor

		wantErr error
		wantLen int
	}{
		{
			name:    "nil pre processor",
			wantErr: errors.New("nil processor"),
			wantLen: 0,
		},
		{
			name:          "valid pre processor",
			postProcessor: &AnyPostProcessor{},
			wantErr:       nil,
			wantLen:       1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			container := NewDefaultContainer()

			// when
			err := container.UsePostProcessor(tc.postProcessor)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			assert.Len(t, container.postProcessors, tc.wantLen)
		})
	}
}
