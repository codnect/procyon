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
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// FileResource struct represents a file-based resource.
type FileResource struct {
	fileSystem fs.FS
	name       string
	path       string
}

// NewFileResource function creates a new FileResource with the given file path.
func NewFileResource(path string) *FileResource {
	dir, fileName := filepath.Split(path)
	fileSystem := os.DirFS(dir)
	return &FileResource{
		fileSystem: fileSystem,
		name:       fileName,
		path:       path,
	}
}

// Name method returns the name of the file resource.
func (r *FileResource) Name() string {
	return r.name
}

// Location method returns the location of the file resource.
func (r *FileResource) Location() string {
	return r.path
}

// Exists method checks if the file resource exists.
func (r *FileResource) Exists() bool {
	if _, err := fs.Stat(r.fileSystem, r.name); err != nil {
		return false
	}

	return true
}

// Reader method opens the file resource and returns an io.ReadCloser.
func (r *FileResource) Reader() (io.ReadCloser, error) {
	return r.File()
}

// File method opens the file resource and returns a fs.File.
func (r *FileResource) File() (fs.File, error) {
	return r.fileSystem.Open(r.name)
}
