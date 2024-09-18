package config

import (
	"codnect.io/procyon/runtime/property"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	FileName = "procyon"
)

// ResourceResolver is an interface that represents a resource resolver.
type ResourceResolver interface {
	ResolveResources(ctx context.Context, location string, profiles []string) ([]Resource, error)
}

type DefaultResourceResolver struct {
	loaders    []property.SourceLoader
	configName string
}

// NewDefaultResourceResolver function creates a new DefaultResourceResolver with the provided loaders.
func NewDefaultResourceResolver(loaders []property.SourceLoader) *DefaultResourceResolver {
	return &DefaultResourceResolver{
		loaders:    loaders,
		configName: FileName,
	}
}

// ResolveResources method resolves resources from a location for specific profiles.
// It returns a list of resources and an error if the resolution fails.
func (r *DefaultResourceResolver) ResolveResources(ctx context.Context, location string, profiles []string) ([]Resource, error) {
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

// getResources method gets resources from a location for a specific profile.
// It returns a list of resources.
func (r *DefaultResourceResolver) getResources(profile string, location string) []Resource {
	var (
		configFile fs.File
		resources  = make([]Resource, 0)
	)

	for _, loader := range r.loaders {
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

				resources = append(resources, newFileResource(filePath, configFile, loader))
			}
		}
	}

	return resources
}
