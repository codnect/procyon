package procyon

import (
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/config"
	"codnect.io/procyon-core/runtime/property"
	"context"
	"strings"
)

type environmentConfigurer struct {
	loaders  []property.SourceLoader
	importer *config.Importer
}

func newEnvironmentConfigurer(loaders []property.SourceLoader, importer *config.Importer) *environmentConfigurer {
	return &environmentConfigurer{
		loaders:  loaders,
		importer: importer,
	}
}

func (c *environmentConfigurer) ConfigureEnvironment(ctx context.Context, environment runtime.Environment) error {
	return c.importConfig(environment)
}

func (c *environmentConfigurer) importConfig(environment runtime.Environment) error {
	defaultConfigs, err := c.importer.Import(context.Background(), "resources", environment.DefaultProfiles())

	if err != nil {
		return err
	}

	sources := property.NewSources()

	for _, defaultConfig := range defaultConfigs {
		sources.AddLast(defaultConfig.PropertySource())
	}

	activeProfiles := environment.ActiveProfiles()

	if len(activeProfiles) == 0 {
		resolver := property.NewSourcesResolver(sources.ToSlice()...)
		value, ok := resolver.Property("procyon.profiles.active")

		if ok {
			activeProfiles = strings.Split(strings.TrimSpace(value.(string)), ",")
		}
	}

	if len(activeProfiles) != 0 {
		err = environment.SetActiveProfiles(activeProfiles...)
		if err != nil {
			return err
		}

		err = c.loadActiveProfiles(environment, sources, activeProfiles)
		if err != nil {
			return err
		}
	}

	c.mergeSources(environment, sources)
	return nil
}

func (c *environmentConfigurer) loadActiveProfiles(environment runtime.Environment, sourceList *property.Sources, activeProfiles []string) error {
	configs, err := c.importer.Import(context.Background(), "config", activeProfiles)
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		propertySource := cfg.PropertySource()
		sourceList.AddFirst(propertySource)

		err = c.activateIncludeProfiles(environment, sourceList, propertySource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *environmentConfigurer) activateIncludeProfiles(environment runtime.Environment, sourceList *property.Sources, source property.Source) error {
	value, ok := source.Property("procyon.profiles.include")

	if ok {
		profiles := strings.Split(strings.TrimSpace(value.(string)), ",")

		for _, profile := range profiles {
			err := environment.AddActiveProfile(strings.TrimSpace(profile))
			if err != nil {
				return err
			}
		}

		err := c.loadActiveProfiles(environment, sourceList, profiles)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *environmentConfigurer) mergeSources(environment runtime.Environment, sourceList *property.Sources) {
	for _, source := range sourceList.ToSlice() {
		environment.PropertySources().AddLast(source)
	}
}
