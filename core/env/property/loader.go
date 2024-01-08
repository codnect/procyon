package property

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

type SourceLoader interface {
	FileExtensions() []string
	LoadSource(name string, reader io.Reader) (Source, error)
}

type YamlPropertySourceLoader struct {
}

func NewYamlPropertySourceLoader() *YamlPropertySourceLoader {
	return &YamlPropertySourceLoader{}
}

func (l *YamlPropertySourceLoader) FileExtensions() []string {
	return []string{"yml", "yaml"}
}

func (l *YamlPropertySourceLoader) LoadSource(name string, reader io.Reader) (Source, error) {
	loaded := make(map[string]interface{})

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		return nil, err
	}

	return NewMapPropertySource(name, loaded), nil
}
