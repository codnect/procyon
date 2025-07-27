package config

import (
	"codnect.io/procyon/runtime/prop"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

const (
	DefaultFileName = "procyon"
)

// ResourceResolver is an interface that represents a resource resolver.
type ResourceResolver interface {
	ResolveResources(ctx context.Context, location string, profiles []string) ([]Resource, error)
}

type FileResourceResolver struct {
	loaders        []prop.SourceLoader
	configFileName string
	fileSystem     fileSystem
}

// NewFileResourceResolver function creates a new FileResourceResolver with the provided loaders.
func NewFileResourceResolver(loaders []prop.SourceLoader) *FileResourceResolver {
	if len(loaders) == 0 {
		panic("no loaders")
	}

	return &FileResourceResolver{
		loaders:        loaders,
		configFileName: DefaultFileName,
		fileSystem:     newOsFileSystem(),
	}
}

// ResolveResources method resolves resources from a location for specific profiles.
// It returns a list of resources and an error if the resolution fails.
func (r *FileResourceResolver) ResolveResources(ctx context.Context, location string, profiles []string) ([]Resource, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if location == "" {
		return nil, errors.New("empty location")
	}

	resources := make([]Resource, 0)

	if len(profiles) == 0 {
		profiles = append(profiles, "default")
	}

	for _, profile := range profiles {
		if profile == "default" {
			profile = ""
		}

		fileResources, err := r.getFileResources(profile, location)
		if err != nil {
			return nil, err
		}

		resources = append(resources, fileResources...)
	}

	return resources, nil
}

// getFileResources method gets resources from a location for a specific profile.
// It returns a list of resources.
func (r *FileResourceResolver) getFileResources(profile string, location string) ([]Resource, error) {
	var (
		configFile fs.File
		resources  = make([]Resource, 0)
	)

	for _, loader := range r.loaders {
		extensions := loader.FileExtensions()

		for _, extension := range extensions {
			filePath := ""

			if profile == "" {
				filePath = filepath.Join(location, fmt.Sprintf("%s.%s", r.configFileName, extension))
			} else {
				filePath = filepath.Join(location, fmt.Sprintf("%s-%s.%s", r.configFileName, profile, extension))
			}

			if _, err := r.fileSystem.Stat(filePath); err == nil {
				configFile, err = r.fileSystem.Open(filePath)

				if err != nil {
					return nil, err
				}

				resources = append(resources, newFileResource(filePath, configFile, loader))
			}
		}
	}

	return resources, nil
}
