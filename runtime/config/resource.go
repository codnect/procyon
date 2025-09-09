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
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"time"

	"codnect.io/procyon/runtime/property"
)

// Resource is an interface that represents a configuration resource.
type Resource interface {
	// Exists method checks if the resource exists.
	Exists() bool
	// Location method returns the location of the resource.
	Location() string
	// Profile method returns the profile of the resource.
	Profile() string
	// Reader method returns an io.ReadCloser for the resource.
	Reader() (io.ReadCloser, error)
	// PropertySourceLoader method returns the property source loader for the resource.
	PropertySourceLoader() property.SourceLoader
}

// FileResource struct represents a file-based configuration resource.
type FileResource struct {
	path    string
	name    string
	profile string
	loader  property.SourceLoader
	fsys    fs.FS
}

// newFileResource function creates a new FileResource with the given path, profile, and property source loader.
func newFileResource(fsys fs.FS, name, path, profile string, loader property.SourceLoader) *FileResource {
	return &FileResource{
		path:    path,
		name:    name,
		profile: profile,
		loader:  loader,
		fsys:    fsys,
	}

}

// Exists method checks if the file resource exists.
func (f *FileResource) Exists() bool {
	if _, err := fs.Stat(f.fsys, f.name); err != nil {
		return false
	}

	return true
}

// Location method returns the location of the file resource.
func (f *FileResource) Location() string {
	return f.path
}

// Profile method returns the profile of the file resource.
func (f *FileResource) Profile() string {
	return f.profile
}

// Reader method returns an io.ReadCloser for the file resource.
func (f *FileResource) Reader() (io.ReadCloser, error) {
	return f.fsys.Open(f.name)
}

// PropertySourceLoader method returns the property source loader for the file resource.
func (f *FileResource) PropertySourceLoader() property.SourceLoader {
	return f.loader
}

// URLResource struct represents a URL-based configuration resource.
type URLResource struct {
	url        *url.URL
	profile    string
	loader     property.SourceLoader
	httpClient *http.Client
}

// newURLResource function creates a new URLResource with the given URL, profile, and property source loader.
func newURLResource(url *url.URL, profile string, loader property.SourceLoader) *URLResource {
	return &URLResource{
		url:     url,
		profile: profile,
		loader:  loader,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Exists method checks if the URL resource exists by making a HEAD request.
func (u *URLResource) Exists() bool {
	resp, err := u.httpClient.Head(u.url.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	if resp.Body != nil {
		_ = resp.Body.Close()
	}

	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
}

// Location method returns the location of the URL resource.
func (u *URLResource) Location() string {
	return u.url.String()
}

// Profile method returns the profile of the URL resource.
func (u *URLResource) Profile() string {
	return u.profile
}

// Reader method returns an io.ReadCloser for the URL resource by making a GET request.
func (u *URLResource) Reader() (io.ReadCloser, error) {
	resp, err := u.httpClient.Get(u.url.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {

		if resp.Body != nil {
			_ = resp.Body.Close()
		}

		return nil, errors.New("resource does not exist")
	}

	return resp.Body, nil
}

// PropertySourceLoader method returns the property source loader for the URL resource.
func (u *URLResource) PropertySourceLoader() property.SourceLoader {
	return u.loader
}
