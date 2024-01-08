package env

type Customizer interface {
	CustomizeEnvironment(environment Environment) error
}
