package config

import (
	"codnect.io/procyon/core/env"
	"codnect.io/procyon/core/env/property"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Resolver interface {
	Resolve(location string) ([]Resource, error)
	ResolveProfiles(profiles []string, location string) ([]Resource, error)
}

type FileResolver struct {
	environment   env.Environment
	sourceLoaders []property.SourceLoader
	configName    string
}

func NewFileResolver(environment env.Environment, sourceLoaders []property.SourceLoader) *FileResolver {
	if environment == nil {
		panic("app: environment cannot be nil")
	}

	if len(sourceLoaders) == 0 {
		panic("app: sourceLoaders cannot be empty")
	}

	resolver := &FileResolver{
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

func (r *FileResolver) Resolve(location string) ([]Resource, error) {
	return r.ResolveProfiles(nil, location)
}

func (r *FileResolver) ResolveProfiles(profiles []string, location string) ([]Resource, error) {
	resources := make([]Resource, 0)
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

func (r *FileResolver) getResources(profile string, location string) []Resource {
	resources := make([]Resource, 0)
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

				resources = append(resources, NewFileResource(filePath, configFile, loader))
			}
		}
	}

	return resources
}
