package config

import (
	"codnect.io/procyon/runtime/prop"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"testing"
)

type FakeFile struct {
	Content      []byte
	ReadPosition int
}

func (f *FakeFile) Stat() (fs.FileInfo, error) {
	return nil, nil
}

func (f *FakeFile) Read(bytes []byte) (int, error) {
	if f.ReadPosition >= len(f.Content) {
		return 0, io.EOF
	}
	n := copy(bytes, f.Content[f.ReadPosition:])
	f.ReadPosition += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

type AnyResource struct {
	mock.Mock
}

func (a *AnyResource) Name() string {
	result := a.Called()
	return result.String(0)
}

func (a *AnyResource) Location() string {
	result := a.Called()
	return result.String(0)
}

func (a *AnyResource) Profile() string {
	result := a.Called()
	return result.String(0)
}

func (a *AnyResource) SourceLoader() prop.SourceLoader {
	result := a.Called()
	return result.Get(0).(prop.SourceLoader)
}

func TestNewFileResource(t *testing.T) {
	fakeFile := &FakeFile{}
	yamlSourceLoader := prop.NewYamlSourceLoader()

	testCases := []struct {
		name   string
		path   string
		file   fs.File
		loader prop.SourceLoader

		wantPanic error
	}{
		{
			name:      "empty path",
			path:      "",
			wantPanic: fmt.Errorf("empty or blank path"),
		},
		{
			name:      "blank path",
			path:      " ",
			wantPanic: fmt.Errorf("empty or blank path"),
		},
		{
			name:      "nil file",
			path:      "anyFilePath",
			file:      nil,
			wantPanic: fmt.Errorf("nil file"),
		},
		{
			name:      "nil loader",
			path:      "anyFilePath",
			file:      fakeFile,
			loader:    nil,
			wantPanic: fmt.Errorf("nil loader"),
		},
		{
			name:   "valid file resource",
			path:   "anyFilePath",
			file:   fakeFile,
			loader: yamlSourceLoader,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					newFileResource(tc.path, tc.file, tc.loader)
				})
				return
			}

			fileResource := newFileResource(tc.path, tc.file, tc.loader)

			// then
			require.NotNil(t, fileResource)
		})
	}
}

func TestFileResource_Name(t *testing.T) {
	// given
	path := "./path/anyFileName.anyExt"
	fakeFile := &FakeFile{}
	loader := prop.NewYamlSourceLoader()

	fileResource := newFileResource(path, fakeFile, loader)

	// when
	name := fileResource.Name()

	// then
	assert.Equal(t, "anyFileName.anyExt", name)
}

func TestFileResource_File(t *testing.T) {
	// given
	path := "./path/anyFileName.anyExt"
	fakeFile := &FakeFile{}
	loader := prop.NewYamlSourceLoader()

	fileResource := newFileResource(path, fakeFile, loader)

	// when
	file := fileResource.File()

	// then
	assert.Equal(t, fakeFile, file)
}

func TestFileResource_Profile(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		file   fs.File
		loader prop.SourceLoader

		wantProfile string
	}{
		{
			name:        "no profile",
			path:        "./path/anyFileName.anyExt",
			file:        &FakeFile{},
			loader:      prop.NewYamlSourceLoader(),
			wantProfile: "",
		},
		{
			name:        "any profile",
			path:        "./path/anyFileName-anyProfile.anyExt",
			file:        &FakeFile{},
			loader:      prop.NewYamlSourceLoader(),
			wantProfile: "anyProfile",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			fileResource := newFileResource(tc.path, tc.file, tc.loader)

			// when
			profile := fileResource.Profile()

			// then
			assert.Equal(t, tc.wantProfile, profile)
		})
	}
}

func TestFileResource_Location(t *testing.T) {
	// given
	path := "./path/anyFileName-anyProfile.anyExt"
	fakeFile := &FakeFile{}
	loader := prop.NewYamlSourceLoader()

	fileResource := newFileResource(path, fakeFile, loader)

	// when
	location := fileResource.Location()

	// then
	assert.Equal(t, path, location)
}

func TestFileResource_SourceLoader(t *testing.T) {
	// given
	path := "./path/anyFileName-anyProfile.anyExt"
	fakeFile := &FakeFile{}
	yamlSourceLoader := prop.NewYamlSourceLoader()

	fileResource := newFileResource(path, fakeFile, yamlSourceLoader)

	// when
	loader := fileResource.SourceLoader()

	// then
	assert.Equal(t, yamlSourceLoader, loader)
}
