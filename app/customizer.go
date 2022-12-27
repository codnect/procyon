package app

type ContextCustomizer interface {
	CustomizeContext(ctx Context) error
}

type contextCustomizers []ContextCustomizer

func (c contextCustomizers) invoke(ctx Context) error {
	for _, customizer := range c {
		err := customizer.CustomizeContext(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}
