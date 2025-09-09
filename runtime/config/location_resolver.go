// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"codnect.io/procyon/runtime/property"
)

const (
	// DefaultFileName is the default name of the configuration file.
	DefaultFileName = "procyon"
)

// LocationResolver is an interface that represents a location resolver.
type LocationResolver interface {
	// IsResolvable method checks if the location can be resolved.
	IsResolvable(location string) bool
	// Resolve method resolves the location for specific profiles and returns a list of resources.
	Resolve(ctx context.Context, location string, profiles ...string) ([]Resource, error)
}

// DefaultLocationResolver is a struct that represents a default location resolver.
type DefaultLocationResolver struct {
	loaders    []property.SourceLoader
	fsProvider func(path string) fs.FS
}

// NewDefaultLocationResolver function creates a new DefaultLocationResolver with the provided property source loaders.
func NewDefaultLocationResolver(loaders ...property.SourceLoader) *DefaultLocationResolver {
	return &DefaultLocationResolver{
		loaders:    loaders,
		fsProvider: os.DirFS,
	}
}

// IsResolvable method checks if the location can be resolved.
func (r *DefaultLocationResolver) IsResolvable(location string) bool {
	return true
}

// Resolve method resolves the location for specific profiles and returns a list of resources.
func (r *DefaultLocationResolver) Resolve(ctx context.Context, location string, profiles ...string) ([]Resource, error) {
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

		loadedResources, err := r.getResources(profile, locations)
		if err != nil {
			return nil, err
		}

		resources = append(resources, loadedResources...)
	}

	return resources, nil
}

// getResources method retrieves resources from the given locations for a specific profile.
func (r *DefaultLocationResolver) getResources(profile string, locations []string) ([]Resource, error) {
	resources := make([]Resource, 0)

	for _, location := range locations {
		isDirectory := strings.HasSuffix(location, "/") || strings.HasSuffix(location, string(os.PathSeparator))

		if isDirectory {
			dirResources, err := r.getDirectoryResources(profile, location)
			if err != nil {
				return nil, err
			}

			resources = append(resources, dirResources...)
		} else {
			fileResources, err := r.getFileResources(profile, location)
			if err != nil {
				return nil, err
			}

			resources = append(resources, fileResources...)
		}
	}

	return resources, nil
}

// getDirectoryResources method retrieves resources from a directory for a specific profile.
func (r *DefaultLocationResolver) getDirectoryResources(profile string, location string) ([]Resource, error) {
	resources := make([]Resource, 0)

	for _, loader := range r.loaders {
		extensions := loader.Extensions()

		for _, extension := range extensions {
			filePath := ""

			if profile == "" {
				filePath = fmt.Sprintf("%s%s.%s", location, DefaultFileName, extension)
			} else {
				filePath = fmt.Sprintf("%s%s-%s.%s", location, DefaultFileName, profile, extension)
			}

			resource, err := r.loadResource(filePath, profile, loader)
			if err != nil {
				return nil, err
			}

			if resource.Exists() {
				resources = append(resources, resource)
			}
		}
	}

	return resources, nil
}

// getFileResources method retrieves resources from a file for a specific profile.
func (r *DefaultLocationResolver) getFileResources(profile string, file string) ([]Resource, error) {
	extension := filepath.Ext(file)
	if extension != "" {
		extension = extension[1:]
	}

	resources := make([]Resource, 0)

	for _, loader := range r.loaders {
		if slices.Contains(loader.Extensions(), extension) {
			resource, err := r.loadResource(file, profile, loader)
			if err != nil {
				return nil, err
			}

			if resource.Exists() {
				resources = append(resources, resource)
			}
		}
	}

	return resources, nil
}

// loadResource method loads a resource from the given location and profile using the provided property source loader.
func (r *DefaultLocationResolver) loadResource(location, profile string, loader property.SourceLoader) (Resource, error) {
	scheme := "file"

	configUrl, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	if configUrl.Scheme == "file" {
		location = strings.ReplaceAll(configUrl.String(), "file://", "")
	} else if configUrl.Scheme != "" {
		scheme = configUrl.Scheme
	}

	if scheme == "file" {
		dir, fileName := filepath.Split(location)
		fsys := r.fsProvider(dir)
		return newFileResource(fsys, fileName, location, profile, loader), nil
	}

	return newURLResource(configUrl, profile, loader), nil
}
