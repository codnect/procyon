package app

import (
	"fmt"
	"github.com/procyon-projects/procyon/app/env"
	"github.com/procyon-projects/procyon/app/env/property"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type configResourceResolver struct {
	environment   env.Environment
	sourceLoaders []property.SourceLoader
	configName    string
}

func newConfigResourceResolver(environment env.Environment, sourceLoaders []property.SourceLoader) *configResourceResolver {
	if environment == nil {
		panic("app: environment cannot be nil")
	}

	if len(sourceLoaders) == 0 {
		panic("app: sourceLoaders cannot be empty")
	}

	resolver := &configResourceResolver{
		environment:   environment,
		sourceLoaders: sourceLoaders,
	}

	configNameProperty := environment.PropertyResolver().PropertyOrDefault("procyon.config.name", "procyon")
	resolver.configName = strings.TrimSpace(configNameProperty)

	if resolver.configName == "" {
		panic("app: configName cannot be empty or blank")
	}

	return resolver
}

func (r *configResourceResolver) Resolve(location string) ([]*configResource, error) {
	return r.ResolveProfiles(nil, location)
}

func (r *configResourceResolver) ResolveProfiles(profiles []string, location string) ([]*configResource, error) {
	resources := make([]*configResource, 0)
	if profiles == nil {
		resources = append(resources, r.getResources("", location)...)
		return resources, nil
	}

	for _, profile := range profiles {
		if profile == "default" {
			resources = append(resources, r.getResources("", location)...)
		} else {
			resources = append(resources, r.getResources(profile, location)...)
		}
	}

	return resources, nil
}

func (r *configResourceResolver) getResources(profile string, location string) []*configResource {
	resources := make([]*configResource, 0)
	var configFile fs.File

	for _, loader := range r.sourceLoaders {
		extensions := loader.FileExtensions()

		for _, extension := range extensions {
			filePath := ""

			if profile == "" {
				filePath = filepath.Join(location, fmt.Sprintf("%s.%s", r.configName, extension))
			} else {
				filePath = filepath.Join(location, fmt.Sprintf("%s-%s.%s", r.configName, profile, extension))
			}

			if _, err := os.Stat(filePath); err == nil {
				configFile, err = os.Open(filePath)

				if err != nil {
					continue
				}

				resources = append(resources, newConfigResource(filePath, configFile, loader))
			}
		}
	}

	return resources
}
