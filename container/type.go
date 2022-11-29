package container

import (
	"github.com/procyon-projects/reflector"
	"reflect"
)

type Type struct {
	typ reflector.Type
}

func (t Type) Name() string {
	return t.typ.Name()
}

func (t Type) PackageName() string {
	return t.typ.PackageName()
}

func (t Type) ReflectType() reflect.Type {
	return t.typ.ReflectType()
}

func TypeOf[T any]() *Type {
	typ := reflector.TypeOf[T]()
	return &Type{
		typ: typ,
	}
}

func TypeOfAny[T any](obj T) *Type {
	typ := reflector.TypeOfAny(obj)
	return &Type{
		typ: typ,
	}
}
