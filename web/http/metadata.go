package http

const (
	// MetadataKeyAccepts is the key for the Accepts metadata.
	MetadataKeyAccepts = "Accepts"
	// MetadataKeyProduces is the key for the Produces metadata.
	MetadataKeyProduces = "Produces"
)

// Metadata represents a metadata that is used to provide additional information to the route.
type Metadata interface {
	MetadataKey() string
}

// AcceptsMetadata is used to specify the content types that the route accepts.
type AcceptsMetadata struct {
	contentTypes []string
}

// Accepts creates a new AcceptsMetadata with the provided content types.
func Accepts(contentTypes ...string) AcceptsMetadata {
	metadata := AcceptsMetadata{
		contentTypes: make([]string, len(contentTypes)),
	}
	copy(metadata.contentTypes, contentTypes)
	return metadata
}

// MetadataKey returns the key of the metadata.
func (m AcceptsMetadata) MetadataKey() string {
	return MetadataKeyAccepts
}

// ContentTypes returns the content types that the route accepts.
func (m AcceptsMetadata) ContentTypes() []string {
	return m.contentTypes
}

// ProducesMetadata is used to specify the content types that the route produces.
type ProducesMetadata struct {
	contentTypes []string
}

// Produces creates a new ProducesMetadata with the provided content types.
func Produces(contentTypes ...string) ProducesMetadata {
	metadata := ProducesMetadata{
		contentTypes: make([]string, len(contentTypes)),
	}
	copy(metadata.contentTypes, contentTypes)
	return metadata
}

// MetadataKey returns the key of the metadata.
func (m ProducesMetadata) MetadataKey() string {
	return MetadataKeyProduces
}

// ContentTypes returns the content types that the route produces.
func (m ProducesMetadata) ContentTypes() []string {
	return m.contentTypes
}
