package config

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Loader is an interface that represents a configuration loader.
// It provides methods to check if a resource is loadable and to load configurations from a resource.
type Loader interface {
	IsLoadable(resource Resource) bool
	LoadConfig(ctx context.Context, resource Resource) (*Config, error)
}

// FileLoader is a struct that represents a file loader.
type FileLoader struct {
}

// NewFileLoader function creates a new FileLoader.
func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// IsLoadable method checks if a resource is a file resource.
// It returns true if the resource is a file resource, false otherwise.
func (l *FileLoader) IsLoadable(resource Resource) bool {
	_, canConvert := resource.(*FileResource)
	return canConvert
}

// LoadConfig method loads configurations from a file resource.
// It returns a configuration and an error if the loading fails.
func (l *FileLoader) LoadConfig(ctx context.Context, resource Resource) (*Config, error) {
	if resource == nil {
		return nil, errors.New("nil context")
	}

	if resource == nil {
		return nil, errors.New("nil resource")
	}

	if fileResource, ok := resource.(*FileResource); ok {
		loader := fileResource.Loader()
		source, err := loader.Load(fileResource.Name(), fileResource.File())

		if err != nil {
			return nil, err
		}

		return New(source), err
	}

	return nil, fmt.Errorf("resource '%s' is not supported", reflect.TypeOf(resource).Name())
}
