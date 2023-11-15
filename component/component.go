package component

import (
	"github.com/procyon-projects/procyon/app/condition"
	"github.com/procyon-projects/procyon/container"
)

type Component struct {
	containerOptions []container.Option
	definition       *container.Definition
	conditions       []condition.Condition
}

func New(constructor Constructor, options ...Option) (*Component, error) {
	c := &Component{}

	for _, opt := range options {
		opt.applyOption(c)
	}

	definition, err := container.MakeDefinition(constructor, c.containerOptions...)
	if err != nil {
		return nil, err
	}

	c.definition = definition
	return c, nil
}

func (c *Component) Definition() *container.Definition {
	return c.definition
}

func (c *Component) ConditionalOn(condition condition.Condition) *Component {
	c.conditions = append(c.conditions, condition)
	return c
}

func (c *Component) Conditions() []condition.Condition {
	conditions := make([]condition.Condition, len(c.conditions))
	copy(conditions, c.conditions)
	return conditions
}
