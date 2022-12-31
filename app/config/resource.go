package config

import (
	"github.com/procyon-projects/procyon/app/env/property"
	"io/fs"
	"path/filepath"
	"strings"
)

type Resource interface {
	Name() string
	Location() string
	Profile() string
	Loader() property.SourceLoader
}

type FileResource struct {
	path   string
	file   fs.File
	loader property.SourceLoader
}

func NewFileResource(path string, file fs.File, loader property.SourceLoader) *FileResource {
	if strings.TrimSpace(path) == "" {
		panic("app: path cannot be empty or blank")
	}

	if file == nil {
		panic("app: file cannot be nil")
	}

	if loader == nil {
		panic("app: loader cannot be nil")
	}

	return &FileResource{
		path:   path,
		file:   file,
		loader: loader,
	}
}

func (r *FileResource) File() fs.File {
	return r.file
}

func (r *FileResource) Location() string {
	return r.path
}

func (r *FileResource) Name() string {
	return filepath.Base(r.Location())
}

func (r *FileResource) Profile() string {
	fileName := filepath.Base(r.Location())
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	nameParts := strings.Split(fileName, "-")
	if len(nameParts) == 1 {
		return ""
	}

	return nameParts[len(nameParts)-1]
}

func (r *FileResource) Loader() property.SourceLoader {
	return r.loader
}
