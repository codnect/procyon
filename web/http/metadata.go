package http

type Metadata map[any]any
type MetadataFunc func(metadata Metadata)

type metadataConsumes struct{}
type metadataProduces struct{}

var (
	metadataConsumesKey = &metadataConsumes{}
	metadataProducesKey = &metadataProduces{}
)

// ConsumesMetadata is used to specify the content types that the route accepts.
type ConsumesMetadata struct {
	contentTypes []string
}

// Consumes creates a new AcceptsMetadata with the provided content types.
func Consumes(contentTypes ...string) MetadataFunc {
	return func(metadata Metadata) {
		metadata[metadataConsumesKey] = ConsumesMetadata{
			contentTypes: contentTypes,
		}
	}
}

func ConsumesFromMetadata(metadata Metadata) ConsumesMetadata {
	if consumes, ok := metadata[metadataConsumesKey]; ok {
		return consumes.(ConsumesMetadata)
	}

	return ConsumesMetadata{}
}

// ContentTypes returns the content types that the route accepts.
func (m ConsumesMetadata) ContentTypes() []string {
	return m.contentTypes
}

// ProducesMetadata is used to specify the content types that the route produces.
type ProducesMetadata struct {
	contentTypes []string
}

// Produces creates a new ProducesMetadata with the provided content types.
func Produces(contentTypes ...string) MetadataFunc {
	return func(metadata Metadata) {
		metadata[metadataProducesKey] = ConsumesMetadata{
			contentTypes: contentTypes,
		}
	}
}

func ProducesFromMetadata(metadata Metadata) ProducesMetadata {
	if produces, ok := metadata[metadataProducesKey]; ok {
		return produces.(ProducesMetadata)
	}

	return ProducesMetadata{}
}

// ContentTypes returns the content types that the route produces.
func (m ProducesMetadata) ContentTypes() []string {
	return m.contentTypes
}
