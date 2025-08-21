// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package property

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

// SourceLoader interface provides methods for loading property sources.
type SourceLoader interface {
	SupportsFormat(format string) bool
	Load(name string, reader io.Reader) (Source, error)
}

// YamlSourceLoader struct is an implementation of the SourceLoader interface for YAML contents.
type YamlSourceLoader struct {
}

// NewYamlSourceLoader function creates a new YamlSourceLoader.
func NewYamlSourceLoader() *YamlSourceLoader {
	return &YamlSourceLoader{}
}

// SupportsFormat method returns the format supported by the YamlSourceLoader.
func (l *YamlSourceLoader) SupportsFormat(format string) bool {
	return strings.ToLower(format) == "yaml"
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
