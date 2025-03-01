package metadata

type Key string

// Metadata represents a metadata that is used to provide additional information
type Metadata interface {
	MetadataKey() Key
}
