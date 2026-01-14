package httpx

// ConsumesMetadata is used to specify the content types that the route accepts.
type ConsumesMetadata struct {
	contentTypes []string
}

// Consumes creates a new AcceptsMetadata with the provided content types.
func Consumes(contentTypes ...string) ConsumesMetadata {
	return ConsumesMetadata{
		contentTypes: contentTypes,
	}
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
func Produces(contentTypes ...string) ProducesMetadata {
	return ProducesMetadata{
		contentTypes: contentTypes,
	}
}

// ContentTypes returns the content types that the route produces.
func (m ProducesMetadata) ContentTypes() []string {
	return m.contentTypes
}
