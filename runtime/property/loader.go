package property

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

// SourceLoader interface provides methods for loading property sources.
type SourceLoader interface {
	FileExtensions() []string
	Load(name string, reader io.Reader) (Source, error)
}

// YamlSourceLoader struct is an implementation of the SourceLoader interface for YAML contents.
type YamlSourceLoader struct {
}

// NewYamlSourceLoader function creates a new YamlSourceLoader.
func NewYamlSourceLoader() *YamlSourceLoader {
	return &YamlSourceLoader{}
}

// FileExtensions method returns the file extensions supported by the YamlSourceLoader.
func (l *YamlSourceLoader) FileExtensions() []string {
	return []string{"yaml", "yml"}
}

// Load method loads a property source from a reader.
func (l *YamlSourceLoader) Load(name string, reader io.Reader) (Source, error) {
	if name == "" {
		return nil, fmt.Errorf("empty source name")
	}

	if reader == nil {
		return nil, fmt.Errorf("nil reader")
	}

	loaded := make(map[string]any)

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		return nil, err
	}

	return NewMapSource(name, loaded), nil
}
