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
	"fmt"
	"reflect"
	"slices"
)

// ConstructorFunc represents a function that can be used as a constructor.
type ConstructorFunc any

// Constructor represents a constructor function along with its arguments.
type Constructor struct {
	fnType  reflect.Type  // The type of the constructor function.
	fnValue reflect.Value // The value of the constructor function.
	args    []Arg         // The arguments of the constructor function.
}

// createConstructor validates the given ConstructorFunc and builds a Constructor metadata struct.
// It ensures that the function is non-nil, is of kind Func, and returns exactly one result.
// If valid, it extracts the argument types and returns a populated Constructor.
func createConstructor(fn ConstructorFunc) (Constructor, error) {
	if fn == nil {
		return Constructor{}, fmt.Errorf("nil constructor")
	}

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return Constructor{}, fmt.Errorf("constructor must be a function")
	}

	if fnType.NumOut() != 1 {
		return Constructor{}, fmt.Errorf("constructor must only return one result")
	}

	return Constructor{
		fnType:  fnType,
		fnValue: reflect.ValueOf(fn),
		args:    extractConstructorArgs(fnType),
	}, nil
}

// OutType returns the type of the constructor function's output.
func (f Constructor) OutType() reflect.Type {
	return f.fnType.Out(0)
}

// Args returns a copy of the arguments of the constructor function.
func (f Constructor) Args() []Arg {
	return slices.Clone(f.args)
}

// Invoke invokes the constructor function with the provided arguments.
// It returns the results of the function invocation and an error if the invocation fails.
func (f Constructor) Invoke(args ...any) ([]any, error) {
	numIn := f.fnType.NumIn()
	numOut := f.fnType.NumOut()
	isVariadic := f.fnType.IsVariadic()

	// Check if the number of arguments matches the number of parameters in the function.
	if (isVariadic && len(args) < numIn) || (!isVariadic && len(args) != numIn) {
		return nil, fmt.Errorf("invalid parameter count, expected %d but got %d", numIn, len(args))
	}

	var variadicType reflect.Type
	inputs := make([]reflect.Value, 0)

	if isVariadic {
		variadicType = f.fnType.In(numOut - 1)
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

		expectedArgType := f.fnType.In(index)

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
	results := f.fnValue.Call(inputs)

	for _, result := range results {
		outputs = append(outputs, result.Interface())
	}

	return outputs, nil
}

// Arg represents an argument of a constructor function.
type Arg struct {
	index int          // The index of the argument in the function parameter list.
	name  string       // The name of the argument.
	typ   reflect.Type // The type of the argument.
}

// Index returns the index of the argument in the function parameter list.
func (a Arg) Index() int {
	return a.index
}

// Name returns the name of the argument.
func (a Arg) Name() string {
	return a.name
}

// Type returns the type of the argument.
func (a Arg) Type() reflect.Type {
	return a.typ
}

// extractConstructorArgs returns the list of argument metadata for a given function type.
func extractConstructorArgs(fnType reflect.Type) []Arg {
	numIn := fnType.NumIn()
	args := make([]Arg, 0, numIn)

	for index := 0; index < numIn; index++ {
		args = append(args, Arg{
			index: index,
			typ:   fnType.In(index),
		})
	}

	return args
}
