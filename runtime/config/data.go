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
	"os"
	"path/filepath"
	"slices"
	"strings"

	"codnect.io/procyon/io"
)

const (
	// DefaultFileName is the default name of the configuration file.
	DefaultFileName = "procyon"
	// DefaultProfile is the default profile name.
	DefaultProfile = "default"
)

// Data is a struct that represents a configuration data.
type Data struct {
	propertySource PropertySource
	profile        string
}

// DataLoader is an interface that represents a data loader.
type DataLoader interface {
	// Load method loads configuration data from the given location for specific profiles.
	Load(ctx context.Context, location string, profiles ...string) ([]Data, error)
}

// NewData function creates a new Data with the provided property source and profile.
func NewData(propertySource PropertySource, profile string) Data {
	if propertySource == nil {
		panic("nil property source")
	}

	if profile == "" {
		profile = DefaultProfile
	}

	return Data{
		propertySource: propertySource,
		profile:        profile,
	}
}

// PropertySource method returns the property source of the data.
func (d Data) PropertySource() PropertySource {
	return d.propertySource
}

// Profile method returns the profile of the data.
func (d Data) Profile() string {
	return d.profile
}

// StandardDataLoader is a struct that represents a default location resolver.
type StandardDataLoader struct {
	resourceResolver io.ResourceResolver
	loaders          []PropertySourceLoader
}

// NewStandardDataLoader function creates a new StandardDataLoader with the provided property source loaders.
func NewStandardDataLoader(resourceResolver io.ResourceResolver, loaders ...PropertySourceLoader) *StandardDataLoader {
	if resourceResolver == nil {
		panic("nil resource resolver")
	}

	return &StandardDataLoader{
		resourceResolver: resourceResolver,
		loaders:          loaders,
	}
}

// Load method loads configuration data from the given location for specific profiles.
func (r *StandardDataLoader) Load(ctx context.Context, location string, profiles ...string) ([]Data, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	location = strings.TrimSpace(location)
	if location == "" {
		return nil, errors.New("empty or blank location")
	}

	data := make([]Data, 0)
	locations := strings.Split(location, ";")

	if len(profiles) == 0 {
		profiles = append(profiles, DefaultProfile)
	}

	for _, profile := range profiles {
		if profile == DefaultProfile {
			profile = ""
		}

		loadedData, err := r.loadData(ctx, profile, locations)
		if err != nil {
			return nil, err
		}

		data = append(data, loadedData...)
	}

	return data, nil
}

// loadData method retrieves configuration data from the given locations for a specific profile.
func (r *StandardDataLoader) loadData(ctx context.Context, profile string, locations []string) ([]Data, error) {
	resources := make([]Data, 0)

	for _, location := range locations {
		isDirectory := strings.HasSuffix(location, "/") || strings.HasSuffix(location, string(os.PathSeparator))

		if isDirectory {
			dirResources, err := r.loadFromDir(ctx, profile, location)
			if err != nil {
				return nil, err
			}

			resources = append(resources, dirResources...)
		} else {
			fileResources, err := r.loadFromFile(ctx, profile, location)
			if err != nil {
				return nil, err
			}

			resources = append(resources, fileResources...)
		}
	}

	return resources, nil
}

// loadFromDir method loads configuration data from a directory for a specific profile.
func (r *StandardDataLoader) loadFromDir(ctx context.Context, profile string, location string) ([]Data, error) {

	data := make([]Data, 0)

	for _, loader := range r.loaders {
		extensions := loader.Extensions()

		for _, extension := range extensions {
			filePath := ""

			if profile == "" {
				filePath = fmt.Sprintf("%s%s.%s", location, DefaultFileName, extension)
			} else {
				filePath = fmt.Sprintf("%s%s-%s.%s", location, DefaultFileName, profile, extension)
			}

			resource, err := r.resourceResolver.Resolve(ctx, filePath)
			if err != nil {
				return nil, err
			}

			if !resource.Exists() {
				continue
			}

			propSource, loadErr := loader.Load(ctx, filePath, resource)
			if loadErr != nil {
				return nil, loadErr
			}

			data = append(data, NewData(propSource, profile))
		}
	}

	return data, nil
}

// loadFromFile method loads configuration data from a file for a specific profile.
func (r *StandardDataLoader) loadFromFile(ctx context.Context, profile string, file string) ([]Data, error) {

	extension := filepath.Ext(file)
	if extension != "" {
		extension = extension[1:]
	}

	data := make([]Data, 0)

	for _, loader := range r.loaders {
		if slices.Contains(loader.Extensions(), extension) {
			resource, err := r.resourceResolver.Resolve(ctx, file)
			if err != nil {
				return nil, err
			}

			if !resource.Exists() {
				continue
			}

			propSource, loadErr := loader.Load(ctx, file, resource)
			if loadErr != nil {
				return nil, loadErr
			}

			data = append(data, NewData(propSource, profile))
		}
	}

	return data, nil

}
