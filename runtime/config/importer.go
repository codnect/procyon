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
)

// Importer interface represents a configuration importer.
type Importer interface {
	Import(ctx context.Context, location string, profiles ...string) ([]*Data, error)
}

// DefaultImporter struct is responsible for importing configurations.
// It uses ResourceResolvers to resolve resources and Loaders to load configurations from these resources.
type DefaultImporter struct {
	resolvers []LocationResolver // The resolvers used to resolve resources.
	loaders   []Loader           // The loaders used to load configurations from resources.
}

// NewDefaultImporter function creates a new Importer with the provided resolvers and loaders.
func NewDefaultImporter(resolvers []LocationResolver, loaders []Loader) *DefaultImporter {
	return &DefaultImporter{
		resolvers: resolvers,
		loaders:   loaders,
	}
}

// Import method imports configurations from a location for specific profiles.
// It first resolves resources from the location and then loads configurations from these resources.
func (i *DefaultImporter) Import(ctx context.Context, location string, profiles []string) ([]*Data, error) {
	resources, err := i.resolve(ctx, location, profiles)
	if err != nil {
		return nil, err
	}

	return i.load(ctx, resources)
}

// resolve method resolves resources from a location for specific profiles using the resolvers
func (i *DefaultImporter) resolve(ctx context.Context, location string, profiles []string) ([]Resource, error) {
	resources := make([]Resource, 0)

	for _, resolver := range i.resolvers {
		if !resolver.IsResolvable(location) {
			continue
		}

		resolved, err := resolver.Resolve(ctx, location, profiles...)

		if err != nil {
			return nil, err
		}

		resources = append(resources, resolved...)
	}

	return resources, nil
}

// load method loads configurations from the resolved resources using the loaders.
func (i *DefaultImporter) load(ctx context.Context, resources []Resource) ([]*Data, error) {
	loaded := make([]*Data, 0)

	for _, resource := range resources {
		loader, err := i.findLoader(resource)

		if err != nil {
			return nil, err
		}

		var data *Data
		data, err = loader.Load(ctx, resource)

		loaded = append(loaded, data)
	}

	return loaded, nil
}

// findLoader method finds a suitable loader for a resource.
// It returns an error if no loader or multiple loaders are found for the resource.
func (i *DefaultImporter) findLoader(resource Resource) (Loader, error) {
	var result Loader

	for _, loader := range i.loaders {
		if loader.IsLoadable(resource) {

			if result != nil {
				return nil, errors.New("multiple loaders")
			}

			result = loader
		}
	}

	if result == nil {
		return nil, errors.New("no loaders")
	}

	return result, nil
}
