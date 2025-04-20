package component

import "reflect"

// Binder defines a method for binding an instance to a specific type.
type Binder interface {
	// Bind associates the given instance with the specified type.
	// Returns an error if the binding fails.
	Bind(typ reflect.Type, instance any) error
}
