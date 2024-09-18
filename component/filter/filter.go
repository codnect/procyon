package filter

import (
	"reflect"
)

// Filter function type represents a filter function that filters components.
type Filter func(filters *Filters)

// Filters struct contains criteria used to filter components.
type Filters struct {
	Name string
	Type reflect.Type
}

// Of function applies a list of filters to a Filters struct and returns it.
func Of(filters ...Filter) *Filters {
	filterOpts := &Filters{}

	for _, filter := range filters {
		filter(filterOpts)
	}

	return filterOpts
}

// ByName function returns a filter that filters components by name.
func ByName(name string) Filter {
	return func(filters *Filters) {
		filters.Name = name
	}
}

// ByTypeOf function returns a filter that filters components by type.
func ByTypeOf[T any]() Filter {
	return func(filters *Filters) {
		typ := reflect.TypeFor[T]()
		filters.Type = typ
	}
}

// ByType function returns a filter that filters components by type.
func ByType(typ reflect.Type) Filter {
	return func(filters *Filters) {
		filters.Type = typ
	}
}
