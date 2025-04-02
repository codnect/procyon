package procyon

import (
	"codnect.io/procyon/runtime"
	"fmt"
	"reflect"
)

// Use creates a new instance of the module and initializes it.
func Use[M runtime.Module]() {
	moduleType := reflect.TypeFor[M]()
	if moduleType.Kind() == reflect.Struct {
		moduleValue := reflect.New(moduleType)

		m := moduleValue.Interface().(runtime.Module)
		err := m.InitModule()
		if err != nil {
			panic(fmt.Errorf("failed to initialize the module '%s': %e", moduleType.Name(), err))
		}
	}
}
