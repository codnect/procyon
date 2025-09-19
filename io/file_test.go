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
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"path/filepath"
	"testing"
	"time"
)

type FakeFile struct {
	contents string
	offset   int
	fileInfo fs.FileInfo
}

func (f *FakeFile) Reset() *FakeFile {
	f.offset = 0
	return f
}

func (f *FakeFile) Stat() (fs.FileInfo, error) {
	return f.fileInfo, nil
}

func (f *FakeFile) Read(p []byte) (int, error) {
	if f.offset >= len(f.contents) {
		return 0, io.EOF
	}
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

type AnyFileInfo struct {
	mock.Mock
}

func (i *AnyFileInfo) Name() string {
	return i.Called().String(0)
}

func (i *AnyFileInfo) Size() int64 {
	return int64(i.Called().Int(0))
}

func (i *AnyFileInfo) Mode() fs.FileMode {
	return i.Called().Get(0).(fs.FileMode)
}

func (i *AnyFileInfo) ModTime() time.Time {
	return i.Called().Get(0).(time.Time)
}
func (i *AnyFileInfo) IsDir() bool {
	return i.Called().Bool(0)
}

func (i *AnyFileInfo) Sys() any {
	return i.Called().Get(0)
}

type AnyDirFs struct {
	dir string
	mock.Mock
}

func (a *AnyDirFs) Open(name string) (fs.File, error) {
	result := a.Called(fmt.Sprintf("%s%s", a.dir, name))
	file := result.Get(0)
	if file == nil {
		return nil, result.Error(1)
	}

	return file.(fs.File), result.Error(1)
}

func (a *AnyDirFs) Stat(name string) (fs.FileInfo, error) {
	result := a.Called(fmt.Sprintf("%s%s", a.dir, name))
	fileInfo := result.Get(0)
	if fileInfo == nil {
		return nil, result.Error(1)
	}

	return fileInfo.(fs.FileInfo), result.Error(1)
}

func TestFileResource_Name(t *testing.T) {
	// given
	fileResource := NewFileResource("anyPath/anyName.anyExt")

	// when
	location := fileResource.Name()

	// then
	require.Equal(t, "anyName.anyExt", location)
}

func TestFileResource_Location(t *testing.T) {
	// given
	fileResource := NewFileResource("anyPath/anyName.anyExt")

	// when
	location := fileResource.Location()

	// then
	require.Equal(t, "anyPath/anyName.anyExt", location)
}

func TestFileResource_Exists(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(anyDirFs *AnyDirFs)
		path         string

		wantErr    error
		wantResult bool
	}{
		{
			name: "resource exists",
			preCondition: func(anyDirFs *AnyDirFs) {
				fileInfo := &AnyFileInfo{}
				anyDirFs.On("Stat", fmt.Sprintf("resource/procyon.yaml")).
					Return(fileInfo, nil).Once()
			},
			path:       "resource/procyon.yaml",
			wantResult: true,
		},
		{
			name: "resource does not exist",
			preCondition: func(anyDirFs *AnyDirFs) {
				anyDirFs.On("Stat", fmt.Sprintf("resource/procyon.yaml")).
					Return(nil, errors.New("resource does not exist")).Once()
			},
			path:       "resource/procyon.yaml",
			wantErr:    errors.New("resource does not exist"),
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			dir, _ := filepath.Split(tc.path)
			anyDirFs := &AnyDirFs{
				dir: dir,
			}

			if tc.preCondition != nil {
				tc.preCondition(anyDirFs)
			}

			fileResource := NewFileResource(tc.path)
			fileResource.fileSystem = anyDirFs

			// when
			exists := fileResource.Exists()

			// then
			require.Equal(t, tc.wantResult, exists)
		})
	}
}

func TestFileResource_Reader(t *testing.T) {
	fakeFile := &FakeFile{}

	testCases := []struct {
		name         string
		preCondition func(anyDirFs *AnyDirFs)
		path         string

		wantErr    error
		wantReader io.ReadCloser
	}{
		{
			name: "reader error",
			preCondition: func(anyDirFs *AnyDirFs) {
				anyDirFs.On("Open", fmt.Sprintf("resource/procyon.yaml")).
					Return(nil, errors.New("reader error")).Once()
			},
			path:    "resource/procyon.yaml",
			wantErr: errors.New("reader error"),
		},
		{
			name: "valid file resource",
			preCondition: func(anyDirFs *AnyDirFs) {
				fakeFile = &FakeFile{}
				anyDirFs.On("Open", fmt.Sprintf("resource/procyon.yaml")).
					Return(fakeFile, nil).Once()
			},
			path:       "resource/procyon.yaml",
			wantReader: fakeFile,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			dir, _ := filepath.Split(tc.path)
			anyDirFs := &AnyDirFs{
				dir: dir,
			}

			if tc.preCondition != nil {
				tc.preCondition(anyDirFs)
			}

			fileResource := NewFileResource(tc.path)
			fileResource.fileSystem = anyDirFs

			// when
			reader, err := fileResource.Reader()

			// then
			if tc.wantErr != nil {
				require.Nil(t, reader)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.Equal(t, tc.wantReader, reader)
		})
	}
}

func TestFileResource_File(t *testing.T) {
	fakeFile := &FakeFile{}

	testCases := []struct {
		name         string
		preCondition func(anyDirFs *AnyDirFs)
		path         string

		wantErr  error
		wantFile io.ReadCloser
	}{
		{
			name: "open error",
			preCondition: func(anyDirFs *AnyDirFs) {
				anyDirFs.On("Open", fmt.Sprintf("resource/procyon.yaml")).
					Return(nil, errors.New("open error")).Once()
			},
			path:    "resource/procyon.yaml",
			wantErr: errors.New("open error"),
		},
		{
			name: "valid file resource",
			preCondition: func(anyDirFs *AnyDirFs) {
				fakeFile = &FakeFile{}
				anyDirFs.On("Open", fmt.Sprintf("resource/procyon.yaml")).
					Return(fakeFile, nil).Once()
			},
			path:     "resource/procyon.yaml",
			wantFile: fakeFile,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			dir, _ := filepath.Split(tc.path)
			anyDirFs := &AnyDirFs{
				dir: dir,
			}

			if tc.preCondition != nil {
				tc.preCondition(anyDirFs)
			}

			fileResource := NewFileResource(tc.path)
			fileResource.fileSystem = anyDirFs

			// when
			file, err := fileResource.File()

			// then
			if tc.wantErr != nil {
				require.Nil(t, file)

				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.Equal(t, tc.wantFile, file)
		})
	}
}
