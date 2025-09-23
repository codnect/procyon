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

package config

import (
	"context"
	"fmt"
	stdio "io"

	"codnect.io/procyon/io"
	"gopkg.in/yaml.v3"
)

type PropertySourceLoader interface {
	Extensions() []string
	Load(ctx context.Context, name string, resource io.Resource) (PropertySource, error)
}

// YamlPropertySourceLoader struct is an implementation of the PropertySourceLoader interface for YAML contents.
type YamlPropertySourceLoader struct {
}

// NewYamlPropertySourceLoader function creates a new YamlSourceLoader.
func NewYamlPropertySourceLoader() *YamlPropertySourceLoader {
	return &YamlPropertySourceLoader{}
}

// Extensions method returns the file extensions supported by the YamlSourceLoader.
func (l *YamlPropertySourceLoader) Extensions() []string {
	return []string{"yaml", "yml"}
}

// Load method loads a config source from a resource.
func (l *YamlPropertySourceLoader) Load(ctx context.Context, name string, resource io.Resource) (PropertySource, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}

	if name == "" {
		return nil, fmt.Errorf("empty source name")
	}

	if resource == nil {
		return nil, fmt.Errorf("nil resource")
	}

	loaded := make(map[string]any)

	reader, err := resource.Reader()
	if err != nil {
		return nil, err
	}

	var data []byte
	data, err = stdio.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		return nil, err
	}

	return NewMapPropertySource(name, loaded), nil
}
