package app

import (
	"github.com/procyon-projects/procyon/app/env"
	"github.com/procyon-projects/procyon/app/env/property"
)

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

type environmentCustomizer struct {
	sourceLoaders []property.SourceLoader
}

func newEnvironmentCustomizer() *environmentCustomizer {
	return &environmentCustomizer{
		sourceLoaders: []property.SourceLoader{
			property.NewYamlPropertySourceLoader(),
		},
	}
}

func (c *environmentCustomizer) CustomizeEnvironment(environment env.Environment) error {
	resolver := newConfigResourceResolver(environment, c.sourceLoaders)
	importer := newConfigImporter([]*configResourceResolver{resolver})

	defaultConfigs, err := importer.Load(environment.DefaultProfiles(), "resources")
	if err != nil {
		return err
	}

	for _, config := range defaultConfigs {
		environment.PropertySources().AddLast(config.PropertySource())
	}

	activeProfiles := environment.ActiveProfiles()
	if len(activeProfiles) != 0 {
		for _, activeProfile := range activeProfiles {
			err = environment.AddActiveProfile(activeProfile)

			if err != nil {
				return err
			}
		}

		return c.loadActiveProfiles(importer, environment, activeProfiles)
	}

	return nil
}

func (c *environmentCustomizer) loadActiveProfiles(importer *configImporter, environment env.Environment, activeProfiles []string) error {
	configs, err := importer.Load(activeProfiles, "resources")
	if err != nil {
		return err
	}

	for _, config := range configs {
		environment.PropertySources().AddLast(config.PropertySource())
	}

	return nil
}
