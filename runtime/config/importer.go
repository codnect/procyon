package config

import (
	"context"
	"errors"
)

// Importer struct is responsible for importing configurations.
// It uses ResourceResolvers to resolve resources and Loaders to load configurations from these resources.
type Importer struct {
	resolvers []ResourceResolver // The resolvers used to resolve resources.
	loaders   []Loader           // The loaders used to load configurations from resources.
}

// NewImporter function creates a new Importer with the provided resolvers and loaders.
func NewImporter(resolvers []ResourceResolver, loaders []Loader) *Importer {
	return &Importer{
		resolvers: resolvers,
		loaders:   loaders,
	}
}

// Import method imports configurations from a location for specific profiles.
// It first resolves resources from the location and then loads configurations from these resources.
func (i *Importer) Import(ctx context.Context, location string, profiles []string) ([]*Data, error) {
	resources, err := i.resolve(ctx, location, profiles)
	if err != nil {
		return nil, err
	}

	return i.load(ctx, resources)
}

// resolve method resolves resources from a location for specific profiles using the resolvers
func (i *Importer) resolve(ctx context.Context, location string, profiles []string) ([]Resource, error) {
	resources := make([]Resource, 0)

	for _, resolver := range i.resolvers {
		resolved, err := resolver.ResolveResources(ctx, location, profiles)

		if err != nil {
			return nil, err
		}

		resources = append(resources, resolved...)
	}

	return resources, nil
}

// load method loads configurations from the resolved resources using the loaders.
func (i *Importer) load(ctx context.Context, resources []Resource) ([]*Data, error) {
	loaded := make([]*Data, 0)

	for _, resource := range resources {
		loader, err := i.findLoader(resource)

		if err != nil {
			return nil, err
		}

		var data *Data
		data, err = loader.LoadData(ctx, resource)

		loaded = append(loaded, data)
	}

	return loaded, nil
}

// findLoader method finds a suitable loader for a resource.
// It returns an error if no loader or multiple loaders are found for the resource.
func (i *Importer) findLoader(resource Resource) (Loader, error) {
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
