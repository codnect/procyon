package container

type Constructor any

type Input struct {
	index    int
	name     string
	optional bool
	typ      Type
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

func (i *Input) Type() Type {
	return i.typ
}
