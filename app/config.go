package app

import (
	"github.com/procyon-projects/procyon/app/env/property"
	"io/fs"
	"path/filepath"
	"strings"
)

type configData struct {
	source property.Source
}

func newConfigData(source property.Source) *configData {
	return &configData{
		source,
	}
}

func (d *configData) PropertySource() property.Source {
	return d.source
}

type configResource struct {
	resourcePath string
	file         fs.File
	loader       property.SourceLoader
}

func newConfigResource(resourcePath string, file fs.File, loader property.SourceLoader) *configResource {
	if strings.TrimSpace(resourcePath) == "" {
		panic("app: resourcePath cannot be empty or blank")
	}

	if file == nil {
		panic("app: file cannot be nil")
	}

	if loader == nil {
		panic("app: loader cannot be nil")
	}

	return &configResource{
		resourcePath: resourcePath,
		file:         file,
		loader:       loader,
	}
}

func (r *configResource) File() fs.File {
	return r.file
}

func (r *configResource) ResourcePath() string {
	return r.resourcePath
}

func (r *configResource) ResourceName() string {
	return filepath.Base(r.ResourcePath())
}

func (r *configResource) Profile() string {
	fileName := filepath.Base(r.ResourcePath())
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	nameParts := strings.Split(fileName, "-")
	if len(nameParts) == 1 {
		return ""
	}

	return nameParts[len(nameParts)-1]
}

func (r *configResource) PropertySourceLoader() property.SourceLoader {
	return r.loader
}
