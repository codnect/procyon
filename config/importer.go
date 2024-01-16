package config

import (
	"codnect.io/procyon-core/env"
	"codnect.io/procyon-core/env/property"
)

type Importer interface {
	Load(profiles []string, location string) ([]*Data, error)
}

type FileImporter struct {
	resolver *FileResolver
}

func NewFileImporter(environment env.Environment) *FileImporter {
	resolver := NewFileResolver(environment, []property.SourceLoader{
		property.NewYamlPropertySourceLoader(),
	})

	return &FileImporter{
		resolver,
	}
}

func (i *FileImporter) Load(profiles []string, location string) ([]*Data, error) {
	resources, err := i.resolveResources(profiles, location)
	if err != nil {
		return nil, err
	}

	return i.loadResources(resources)
}

func (i *FileImporter) resolveResources(profiles []string, location string) ([]Resource, error) {
	configResources := make([]Resource, 0)

	if len(profiles) == 0 {
		resources, err := i.resolver.Resolve(location)

		if err != nil {
			return nil, err
		}

		configResources = append(configResources, resources...)
	} else {
		resources, err := i.resolver.ResolveProfiles(profiles, location)

		if err != nil {
			return nil, err
		}

		configResources = append(configResources, resources...)
	}
	return configResources, nil
}

func (i *FileImporter) loadResources(resources []Resource) ([]*Data, error) {
	configs := make([]*Data, 0)

	for _, resource := range resources {
		fileResource, canConvert := resource.(*FileResource)
		if !canConvert {
			continue
		}

		loader := resource.Loader()
		source, err := loader.LoadSource(fileResource.Location(), fileResource.File())

		if err != nil {
			return nil, err
		}

		configs = append(configs, NewData(source))
	}

	return configs, nil
}
