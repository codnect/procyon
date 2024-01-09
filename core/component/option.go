package component

import (
	"codnect.io/procyon/core/container"
)

type Option interface {
	applyOption(component *Component)
}

type containerOption container.Option

func (d containerOption) applyOption(component *Component) {
	component.containerOptions = append(component.containerOptions, container.Option(d))
}

func Name(name string) Option {
	return containerOption(container.Name(name))
}

func Optional[T any]() Option {
	return containerOption(container.Optional[T]())
}

func OptionalAt(index int) Option {
	return containerOption(container.OptionalAt(index))
}

func Qualifier[T any](name string) Option {
	return containerOption(container.Qualifier[T](name))
}

func QualifierAt(index int, name string) Option {
	return containerOption(container.QualifierAt(index, name))
}

func Scoped(scope string) Option {
	return containerOption(container.Scoped(scope))
}
