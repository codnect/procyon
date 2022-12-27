package env

type Customizer interface {
	CustomizeEnvironment(environment Environment) error
}

type environmentCustomizers struct {
	customizers []Customizer
}

func newEnvironmentCustomizers(customizers []Customizer) *environmentCustomizers {
	return &environmentCustomizers{
		customizers: customizers,
	}
}

func (c *environmentCustomizers) invokeCustomizers(environment Environment) error {
	for _, customizer := range c.customizers {
		err := customizer.CustomizeEnvironment(environment)

		if err != nil {
			return err
		}
	}

	return nil
}
