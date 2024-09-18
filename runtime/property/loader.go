package property

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

// SourceLoader interface provides methods for loading property sources.
type SourceLoader interface {
	FileExtensions() []string
	Load(name string, reader io.Reader) (Source, error)
}

// YamlSourceLoader struct is an implementation of the SourceLoader interface for YAML files.
type YamlSourceLoader struct {
}

// NewYamlSourceLoader function creates a new YamlSourceLoader.
func NewYamlSourceLoader() *YamlSourceLoader {
	return &YamlSourceLoader{}
}

// FileExtensions method returns the file extensions that this loader can handle.
// In this case, it returns "yml" and "yaml".
func (l *YamlSourceLoader) FileExtensions() []string {
	return []string{"yml", "yaml"}
}

// Load method loads a property source from a reader.
func (l *YamlSourceLoader) Load(name string, reader io.Reader) (Source, error) {
	loaded := make(map[string]interface{})

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		return nil, err
	}

	return NewMapSource(name, loaded), nil
}
