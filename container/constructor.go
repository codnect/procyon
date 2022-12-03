package container

import "github.com/procyon-projects/reflector"

type Constructor any

type PostConstructor interface {
	PostConstruct() error
}

type Input struct {
	index    int
	name     string
	optional bool
	typ      *Type
}

func (i *Input) Index() int {
	return i.index
}

func (i *Input) Name() string {
	return i.name
}

func (i *Input) IsOptional() bool {
	return i.optional
}

func (i *Input) Type() *Type {
	return i.typ
}

func (i *Input) reflectorType() reflector.Type {
	return i.typ.typ
}
