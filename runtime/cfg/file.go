package cfg

import (
	"os"
)

type fileSystem interface {
	Stat(name string) (os.FileInfo, error)
	Open(name string) (*os.File, error)
}

type osFileSystem struct {
}

func newOsFileSystem() *osFileSystem {
	return &osFileSystem{}
}

func (o *osFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (o *osFileSystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}
