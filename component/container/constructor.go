package container

import (
	"fmt"
	"reflect"
)

// ConstructorFunc represents a function that can be used as a constructor.
type ConstructorFunc any

// Constructor represents a constructor function along with its arguments.
type Constructor struct {
	funcType  reflect.Type          // The type of the constructor function.
	funcValue reflect.Value         // The value of the constructor function.
	arguments []ConstructorArgument // The arguments of the constructor function.
}

// Name returns the name of the constructor function.
func (f Constructor) Name() string {
	return f.funcType.Name()
}

// Arguments returns a copy of the arguments of the constructor function.
func (f Constructor) Arguments() []ConstructorArgument {
	args := make([]ConstructorArgument, len(f.arguments))
	copy(args, f.arguments)
	return args
}

// Invoke invokes the constructor function with the provided arguments.
// It returns the results of the function invocation and an error if the invocation fails.
func (f Constructor) Invoke(args ...any) ([]any, error) {
	numIn := f.funcType.NumIn()
	numOut := f.funcType.NumOut()
	isVariadic := f.funcType.IsVariadic()

	// Check if the number of arguments matches the number of parameters in the function.
	if (isVariadic && len(args) < numIn) || (!isVariadic && len(args) != numIn) {
		return nil, fmt.Errorf("invalid parameter count, expected %d but got %d", numIn, len(args))
	}

	var variadicType reflect.Type
	inputs := make([]reflect.Value, 0)

	if isVariadic {
		variadicType = f.funcType.In(numOut - 1)
	}

	for index, arg := range args {
		argType := reflect.TypeOf(arg)

		if isVariadic && index > numOut {
			if arg == nil {
				inputs = append(inputs, reflect.New(variadicType.Elem()).Elem())
				continue
			} else if !argType.ConvertibleTo(variadicType.Elem()) {
				return nil, fmt.Errorf("expected %s but got %s at index %d", variadicType.Elem().Name(), argType.Name(), index)
			}

			inputs = append(inputs, reflect.ValueOf(arg))
			continue
		}

		expectedArgType := f.funcType.In(index)

		if arg == nil {
			inputs = append(inputs, reflect.New(expectedArgType).Elem())
		} else {
			if !argType.ConvertibleTo(expectedArgType) {
				return nil, fmt.Errorf("expected %s but got %s at index %d", expectedArgType.Name(), expectedArgType.Name(), index)
			}

			inputs = append(inputs, reflect.ValueOf(arg))
		}
	}

	// Call the function and collect the results.
	outputs := make([]any, 0)
	results := f.funcValue.Call(inputs)

	for _, result := range results {
		outputs = append(outputs, result.Interface())
	}

	return outputs, nil
}

// ConstructorArgument represents an argument of a constructor function.
type ConstructorArgument struct {
	index    int          // The index of the argument in the function parameter list.
	name     string       // The name of the argument.
	typ      reflect.Type // The type of the argument.
	optional bool         // Indicates whether the argument is optional.
}

// ArgumentIndex returns the index of the argument in the function parameter list.
func (a ConstructorArgument) ArgumentIndex() int {
	return a.index
}

// Name returns the name of the argument.
func (a ConstructorArgument) Name() string {
	return a.name
}

// Type returns the type of the argument.
func (a ConstructorArgument) Type() reflect.Type {
	return a.typ
}

// IsOptional returns whether the argument is optional.
func (a ConstructorArgument) IsOptional() bool {
	return a.optional
}
