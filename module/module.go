package module

import (
	"fmt"
	"reflect"
)

// Module is an interface that represents a module.
// It provides a method to initialize the module.
type Module interface {
	// InitModule method initializes the module.
	// It returns an error if the initialization fails.
	InitModule() error
}

// Use creates a new instance of the module and initializes it.
func Use[M Module]() {
	moduleType := reflect.TypeFor[M]()
	if moduleType.Kind() == reflect.Struct {
		moduleValue := reflect.New(moduleType)

		m := moduleValue.Interface().(Module)
		err := m.InitModule()
		if err != nil {
			panic(fmt.Errorf("failed to initialize the module '%s': %e", moduleType.Name(), err))
		}
	}
}
