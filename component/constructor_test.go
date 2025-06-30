package component

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
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
			wantErr:       errors.New("nil constructor"),
		},
		{
			name:          "invalid constructor function",
			constructorFn: "string value",
			wantErr:       errors.New("constructor must be a function"),
		},
		{
			name: "multi result constructor function",
			constructorFn: func() (string, int, error) {
				return "", -1, nil
			},
			wantErr: errors.New("constructor must only return one result"),
		},
		{
			name:          "valid constructor function",
			constructorFn: NewAnotherComponent,
			wantOutType:   reflect.TypeFor[*AnotherComponent](),
			wantArgs: []Arg{
				{
					0,
					"",
					reflect.TypeFor[DependentComponent](),
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
			constructorFn: NewAnotherComponent,
			wantErr:       errors.New("invalid parameter count, expected 1 but got 0"),
		},
		{
			name:          "invalid parameter",
			constructorFn: NewAnotherComponent,
			inputs: []any{
				"invalid parameter",
			},
			wantErr: errors.New("expected DependentComponent but got string at index 0"),
		},
		{
			name:          "nil parameter",
			constructorFn: NewAnotherComponent,
			inputs: []any{
				nil,
			},
			wantOutType: reflect.TypeFor[*AnotherComponent](),
		},
		{
			name:          "valid constructor",
			constructorFn: NewAnotherComponent,
			inputs: []any{
				DependentComponent{},
			},
			wantOutType: reflect.TypeFor[*AnotherComponent](),
		},
		{
			name: "invalid variadic parameter",
			constructorFn: func(dependents ...DependentComponent) *AnotherComponent {
				return &AnotherComponent{}
			},
			inputs: []any{
				"invalid parameter",
			},
			wantErr: errors.New("expected DependentComponent but got string at index 0"),
		},
		{
			name: "nil variadic parameter",
			constructorFn: func(dependents ...DependentComponent) *AnotherComponent {
				return &AnotherComponent{}
			},
			inputs: []any{
				nil,
			},
			wantOutType: reflect.TypeFor[*AnotherComponent](),
		},
		{
			name: "valid variadic constructor",
			constructorFn: func(dependents ...DependentComponent) *AnotherComponent {
				return &AnotherComponent{}
			},
			inputs: []any{
				DependentComponent{},
				DependentComponent{},
			},
			wantOutType: reflect.TypeFor[*AnotherComponent](),
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
