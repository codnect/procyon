package cfg

import (
	"codnect.io/procyon/runtime/property"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	DefaultFileName = "procyon"
)

type ResourceResolver interface {
	IsResolvable(location string) bool
	ResolveResources(ctx context.Context, location string, profiles ...string) ([]Resource, error)
}

type StandardResourceResolver struct {
	fileSystem fileSystem
	loaders    []property.SourceLoader
}

func NewStandardResourceResolver(loaders []property.SourceLoader) *StandardResourceResolver {
	return &StandardResourceResolver{
		loaders: loaders,
	}
}

func (r *StandardResourceResolver) IsResolvable(location string) bool {
	return true
}

func (r *StandardResourceResolver) ResolveResources(ctx context.Context, location string, profiles ...string) ([]Resource, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	location = strings.TrimSpace(location)
	if location == "" {
		return nil, errors.New("empty or blank location")
	}

	resources := make([]Resource, 0)
	locations := strings.Split(location, ";")

	if len(profiles) == 0 {
		profiles = append(profiles, "default")
	}

	for _, profile := range profiles {
		if profile == "default" {
			profile = ""
		}

		fileResources, err := r.getResources(profile, locations)
		if err != nil {
			return nil, err
		}

		resources = append(resources, fileResources...)
	}

	return nil, nil
}

func (r *StandardResourceResolver) getResources(profile string, locations []string) ([]Resource, error) {
	resources := make([]Resource, 0)

	for _, location := range locations {
		isDirectory := strings.HasSuffix(location, "/") || strings.HasSuffix(location, string(os.PathSeparator))

		if isDirectory {
			resolved, err := r.getDirectoryResources(profile, location)
			if err != nil {
				return nil, err
			}

			resources = append(resources, resolved...)
		} else {
			resolved, ok := r.getFileResource(profile, location)
			if ok {
				resources = append(resources, resolved)
			}
		}
	}

	return nil, nil
}

// getFileResources method gets resources from a location for a specific profile.
// It returns a list of resources.
func (r *StandardResourceResolver) getFileResource(profile string, file string) (Resource, bool) {
	ext := filepath.Ext(file)

	var (
		configFile fs.File
		resources  = make([]Resource, 0)
	)

	for _, loader := range r.loaders {
		if slices.Contains(loader.FileExtensions(), ext) {
			if _, err := r.fileSystem.Stat(file); err == nil {
				configFile, err = r.fileSystem.Open(file)

				if err != nil {
					return nil, false
				}

				resources = append(resources, newFileResource(file, configFile, loader))
			}
		}
	}

	return nil, false
}

// getDirectoryResources method gets resources from a location for a specific profile.
// It returns a list of resources.
func (r *StandardResourceResolver) getDirectoryResources(profile string, location string) ([]Resource, error) {
	var (
		configFile fs.File
		resources  = make([]Resource, 0)
	)

	for _, loader := range r.loaders {
		extensions := loader.FileExtensions()

		for _, extension := range extensions {
			filePath := ""

			if profile == "" {
				filePath = filepath.Join(location, fmt.Sprintf("%s.%s", DefaultFileName, extension))
			} else {
				filePath = filepath.Join(location, fmt.Sprintf("%s-%s.%s", DefaultFileName, profile, extension))
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
