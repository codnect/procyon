package procyon

import (
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/env"
)

type contextCustomizers []runtime.ContextCustomizer

func (c contextCustomizers) invoke(ctx runtime.Context) error {
	for _, customizer := range c {
		err := customizer.CustomizeContext(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

type eventCustomizers []env.Customizer

func (e eventCustomizers) invoke(environment env.Environment) error {
	for _, customizer := range e {
		err := customizer.CustomizeEnvironment(environment)

		if err != nil {
			return err
		}
	}

	return nil
}
