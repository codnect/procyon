package cfg

import (
	"codnect.io/procyon/runtime/property"
	"io/fs"
)

type Resource interface {
	Location() string
	Profile() string
}

type FileResource struct {
	path   string
	file   fs.File
	loader property.SourceLoader
}

func newFileResource(path string, file fs.File, loader property.SourceLoader) *FileResource {
	return &FileResource{
		path:   path,
		file:   file,
		loader: loader,
	}
}

func (f *FileResource) Location() string {
	//TODO implement me
	panic("implement me")
}

func (f *FileResource) Profile() string {
	//TODO implement me
	panic("implement me")
}

func (f *FileResource) File() *fs.File {
	//TODO implement me
	panic("implement me")
}
