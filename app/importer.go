package app

type configImporter struct {
	resolvers []*configResourceResolver
}

func newConfigImporter(resolvers []*configResourceResolver) *configImporter {
	if len(resolvers) == 0 {
		panic("app: resolvers cannot be nil or empty")
	}

	return &configImporter{
		resolvers: resolvers,
	}
}

func (i *configImporter) Load(profiles []string, location string) ([]*configData, error) {
	resources, err := i.resolveResources(profiles, location)
	if err != nil {
		return nil, err
	}

	return i.loadResources(resources)
}

func (i *configImporter) resolveResources(profiles []string, location string) ([]*configResource, error) {
	configResources := make([]*configResource, 0)

	for _, resolver := range i.resolvers {
		if len(profiles) == 0 {
			resources, err := resolver.Resolve(location)

			if err != nil {
				return nil, err
			}

			configResources = append(configResources, resources...)
		} else {
			resources, err := resolver.ResolveProfiles(profiles, location)

			if err != nil {
				return nil, err
			}

			configResources = append(configResources, resources...)
		}
	}

	return configResources, nil
}

func (i *configImporter) loadResources(resources []*configResource) ([]*configData, error) {
	configs := make([]*configData, 0)

	for _, resource := range resources {
		loader := resource.PropertySourceLoader()
		source, err := loader.LoadSource(resource.ResourcePath(), resource.File())

		if err != nil {
			return nil, err
		}

		configs = append(configs, newConfigData(source))
	}

	return configs, nil
}
