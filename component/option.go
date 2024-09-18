package component

import (
	"codnect.io/procyon/component/condition"
	"codnect.io/procyon/component/container"
	"fmt"
)

// Option is a function type that modifies a Component.
type Option func(component *Component) error

// WithName sets the name of the component.
func WithName(name string) Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.Named(name))
		return nil
	}
}

// WithScope sets the scope of the component.
func WithScope(scope string) Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.Scoped(scope))
		return nil
	}
}

// WithSingletonScope sets the scope of the component as singleton.
func WithSingletonScope() Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.Scoped(container.SingletonScope))
		return nil
	}
}

// WithPrototypeScope sets the scope of the component as prototype.
func WithPrototypeScope() Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.Scoped(container.PrototypeScope))
		return nil
	}
}

// WithPriority sets the priority of the component.
func WithPriority(priority int) Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.Prioritized(priority))
		return nil
	}
}

// WithQualifier sets the name of the constructor's input parameter.
func WithQualifier[T any](name string) Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.Qualifier[T](name))
		return nil
	}
}

// WithQualifierAt sets the name of the constructor's input parameter at the given index.
func WithQualifierAt(index int, name string) Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.QualifierAt(index, name))
		return nil
	}
}

// WithOptional sets the constructor's input parameter as optional.
func WithOptional[T any]() Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.Optional[T]())
		return nil
	}
}

// WithOptionalAt sets the constructor's input parameter at the given index as optional.
func WithOptionalAt(index int) Option {
	return func(component *Component) error {
		component.definitionOptions = append(component.definitionOptions, container.OptionalAt(index))
		return nil
	}
}

// WithCondition adds a condition for component.
func WithCondition(condition condition.Condition) Option {
	return func(component *Component) error {
		if condition == nil {
			return fmt.Errorf("nil condition")
		}

		component.conditions = append(component.conditions, condition)
		return nil
	}
}
