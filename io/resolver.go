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

package io

import (
	"context"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	// FileScheme represents the "file" scheme.
	FileScheme = "file"
)

// ResourceResolver is an interface for resolving resources from location strings.
type ResourceResolver interface {
	// Resolve method resolves a resource from the given location string.
	Resolve(ctx context.Context, location string) (Resource, error)
}

// DefaultResourceResolver is the default implementation of the ResourceResolver interface.
type DefaultResourceResolver struct {
}

// NewDefaultResourceResolver creates a new instance of DefaultResourceResolver.
func NewDefaultResourceResolver() *DefaultResourceResolver {
	return &DefaultResourceResolver{}
}

// Resolve method resolves a resource from the given location string.
func (r *DefaultResourceResolver) Resolve(_ context.Context, location string) (Resource, error) {
	resourceUrl, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	scheme := strings.ToLower(resourceUrl.Scheme)
	switch scheme {
	case "", FileScheme:
		fileLoc := normalizeLocation(resourceUrl)
		return NewFileResource(fileLoc), nil
	default:
		return NewURLResource(resourceUrl), nil
	}
}

// normalizeLocation normalizes the file location from the URL and raw location string.
func normalizeLocation(url *url.URL) string {
	fileLoc := url.Path
	if fileLoc == "" {
		fileLoc = url.Opaque
	}

	if url.Host != "" {
		fileLoc = url.Host + fileLoc
	}

	return filepath.Clean(fileLoc)
}
