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
	"errors"
)

// Loader is an interface that represents a configuration loader.
// It provides methods to check if a resource is loadable and to load configurations from a resource.
type Loader interface {
	// IsLoadable method checks if a resource is loadable by the loader.
	IsLoadable(resource Resource) bool
	// Load method loads configurations from a resource.
	Load(ctx context.Context, resource Resource) (*Data, error)
}

// DefaultLoader is a struct that represents a file loader.
type DefaultLoader struct {
}

// NewDefaultLoader function creates a new DefaultLoader.
func NewDefaultLoader() *DefaultLoader {
	return &DefaultLoader{}
}

// IsLoadable method checks if a resource is a file resource.
// It returns true if the resource is a file resource, false otherwise.
func (l *DefaultLoader) IsLoadable(resource Resource) bool {
	if resource == nil {
		return false
	}

	if resource.PropertySourceLoader() == nil {
		return false
	}

	return resource.Exists()
}

// Load method loads configurations from a file resource.
// It returns a configuration and an error if the loading fails.
func (l *DefaultLoader) Load(ctx context.Context, resource Resource) (*Data, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if resource == nil {
		return nil, errors.New("nil resource")
	}

	loader := resource.PropertySourceLoader()
	if loader == nil {
		return nil, errors.New("nil resource loader")
	}

	if !resource.Exists() {
		return nil, errors.New("no resource found")
	}

	reader, readerErr := resource.Reader()
	if readerErr != nil {
		return nil, readerErr
	}

	defer reader.Close()
	source, loadErr := loader.Load(resource.Location(), reader)
	if loadErr != nil {
		return nil, loadErr
	}

	return NewData(source), nil
}
