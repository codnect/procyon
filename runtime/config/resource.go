package config

import (
	"codnect.io/procyon/runtime/property"
	"io/fs"
	"path/filepath"
	"strings"
)

// Resource is an interface that represents a resource.
type Resource interface {
	Name() string
	Location() string
	Profile() string
	Loader() property.SourceLoader
}

// FileResource is a struct that represents a file resource.
type FileResource struct {
	path   string
	file   fs.File
	loader property.SourceLoader
}

// newFileResource function creates a new FileResource with the provided path, file, and loader.
func newFileResource(path string, file fs.File, loader property.SourceLoader) *FileResource {
	if strings.TrimSpace(path) == "" {
		panic("cannot create file resource with empty or blank path")
	}

	if file == nil {
		panic("nil file")
	}

	if loader == nil {
		panic("nil loader")
	}

	return &FileResource{
		path:   path,
		file:   file,
		loader: loader,
	}
}

// File method returns the file.
func (r *FileResource) File() fs.File {
	return r.file
}

// Location method returns the location of the file.
func (r *FileResource) Location() string {
	return r.path
}

// Name method returns the name of the file.
func (r *FileResource) Name() string {
	return filepath.Base(r.Location())
}

// Profile method returns the profile of the file.
func (r *FileResource) Profile() string {
	fileName := filepath.Base(r.Location())
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	nameParts := strings.Split(fileName, "-")
	if len(nameParts) == 1 {
		return ""
	}

	return nameParts[len(nameParts)-1]
}

// Loader method returns the loader of the file.
func (r *FileResource) Loader() property.SourceLoader {
	return r.loader
}
