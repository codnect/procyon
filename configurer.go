package procyon

import (
	"codnect.io/procyon/runtime"
	"codnect.io/procyon/runtime/config"
	"codnect.io/procyon/runtime/property"
	"context"
	"strings"
)

type configContextConfigurer struct {
	loaders  []property.SourceLoader
	importer *config.Importer
}

func newConfigContextConfigurer(loaders []property.SourceLoader, importer *config.Importer) *configContextConfigurer {
	return &configContextConfigurer{
		loaders:  loaders,
		importer: importer,
	}
}

func (c *configContextConfigurer) ConfigureContext(ctx runtime.Context) error {
	return c.importConfig(ctx.Environment())
}

func (c *configContextConfigurer) importConfig(environment runtime.Environment) error {
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
		resolver := property.NewSourcesResolver(sources)
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

func (c *configContextConfigurer) loadActiveProfiles(environment runtime.Environment, sourceList *property.Sources, activeProfiles []string) error {
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

func (c *configContextConfigurer) activateIncludeProfiles(environment runtime.Environment, sourceList *property.Sources, source property.Source) error {
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

func (c *configContextConfigurer) mergeSources(environment runtime.Environment, sourceList *property.Sources) {
	for _, source := range sourceList.ToSlice() {
		environment.PropertySources().AddLast(source)
	}
}
