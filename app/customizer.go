package app

type ContextCustomizer interface {
	CustomizeContext(ctx Context) error
}

type contextCustomizers struct {
	customizers []ContextCustomizer
}

func newContextCustomizers(customizers []ContextCustomizer) *contextCustomizers {
	return &contextCustomizers{
		customizers: customizers,
	}
}

func (c *contextCustomizers) invokeCustomizers(ctx Context) error {
	for _, customizer := range c.customizers {
		err := customizer.CustomizeContext(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}
