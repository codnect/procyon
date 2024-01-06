package container

import "codnect.io/reflector"

type Constructor any

type PostConstructor interface {
	PostConstruct() error
}

type Input struct {
	index    int
	name     string
	optional bool
	typ      reflector.Type
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

func (i *Input) Type() reflector.Type {
	return i.typ
}
