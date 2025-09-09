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
	"codnect.io/procyon/runtime/property"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"syscall"
	"testing"
	"time"
)

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

type FakeFile struct {
	contents string
	offset   int
	fileInfo fs.FileInfo
}

// Reset prepares a FakeFile for reuse.
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

func TestFileResource_Exists(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(anyDirFs *AnyDirFs)
		resourceName string
		path         string
		profile      string
		loader       property.SourceLoader

		wantErr    error
		wantResult bool
	}{
		{
			name: "resource exists",
			preCondition: func(anyDirFs *AnyDirFs) {
				fileInfo := &AnyFileInfo{}
				anyDirFs.On("Stat", fmt.Sprintf("anyPath/anyName")).Return(fileInfo, nil).Once()
			},
			resourceName: "anyName",
			path:         "anyPath/",
			profile:      "anyProfile",
			loader:       property.NewYamlSourceLoader(),
			wantResult:   true,
		},
		{
			name: "resource does not exist",
			preCondition: func(anyDirFs *AnyDirFs) {
				anyDirFs.On("Stat", fmt.Sprintf("anyPath/anyName")).
					Return(nil, errors.New("resource does not exist")).Once()
			},
			resourceName: "anyName",
			path:         "anyPath/",
			profile:      "anyProfile",
			wantErr:      errors.New("resource does not exist"),
			wantResult:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			anyDirFs := &AnyDirFs{
				dir: tc.path,
			}

			if tc.preCondition != nil {
				tc.preCondition(anyDirFs)
			}

			res := newFileResource(anyDirFs, tc.resourceName, tc.path, tc.profile, tc.loader)

			// when
			exists := res.Exists()

			// then
			require.Equal(t, tc.wantResult, exists)
		})
	}
}

func TestFileResource_Location(t *testing.T) {
	// given
	res := newFileResource(&AnyDirFs{}, "anyName", "anyPath", "anyProfile", property.NewYamlSourceLoader())

	// when
	location := res.Location()

	// then
	require.Equal(t, "anyPath", location)
}

func TestFileResource_Profile(t *testing.T) {
	// given
	res := newFileResource(&AnyDirFs{}, "anyName", "anyPath", "anyProfile", property.NewYamlSourceLoader())

	// when
	profile := res.Profile()

	// then
	require.Equal(t, "anyProfile", profile)
}

func TestFileResource_Reader(t *testing.T) {
	fakeFile := &FakeFile{}

	testCases := []struct {
		name         string
		preCondition func(anyDirFs *AnyDirFs)
		resourceName string
		path         string
		profile      string
		loader       property.SourceLoader

		wantErr    error
		wantReader io.ReadCloser
	}{
		{
			name: "reader error",
			preCondition: func(anyDirFs *AnyDirFs) {
				anyDirFs.On("Open", fmt.Sprintf("anyPath/anyName")).
					Return(nil, errors.New("reader error")).Once()
			},
			resourceName: "anyName",
			path:         "anyPath/",
			profile:      "anyProfile",
			wantErr:      errors.New("reader error"),
		},
		{
			name: "valid file resource",
			preCondition: func(anyDirFs *AnyDirFs) {
				fakeFile = &FakeFile{}
				anyDirFs.On("Open", fmt.Sprintf("anyPath/anyName")).
					Return(fakeFile, nil).Once()
			},
			resourceName: "anyName",
			path:         "anyPath/",
			profile:      "anyProfile",
			wantReader:   fakeFile,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			anyDirFs := &AnyDirFs{
				dir: tc.path,
			}

			if tc.preCondition != nil {
				tc.preCondition(anyDirFs)
			}

			res := newFileResource(anyDirFs, tc.resourceName, tc.path, tc.profile, tc.loader)

			// when
			reader, err := res.Reader()

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

func TestFileResource_PropertySourceLoader(t *testing.T) {
	// given
	loader := property.NewYamlSourceLoader()
	res := newFileResource(&AnyDirFs{}, "anyName", "anyPath", "anyProfile", loader)

	// when
	result := res.PropertySourceLoader()

	// then
	require.Equal(t, loader, result)
}

func TestURLResource_Exists(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/anyValidResource" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	validResourceUrl, _ := url.Parse(fmt.Sprintf("%s/%s", testServer.URL, "anyValidResource"))
	nonValidResourceUrl, _ := url.Parse(fmt.Sprintf("%s/%s", testServer.URL, "anyNonValidResource"))

	testCases := []struct {
		name    string
		url     *url.URL
		profile string
		loader  property.SourceLoader

		wantErr    error
		wantResult bool
	}{
		{
			name:       "resource exists",
			url:        validResourceUrl,
			profile:    "anyProfile",
			loader:     property.NewYamlSourceLoader(),
			wantResult: true,
		},
		{
			name:       "resource does not exist",
			url:        nonValidResourceUrl,
			profile:    "anyProfile",
			loader:     property.NewYamlSourceLoader(),
			wantErr:    errors.New("resource does not exist"),
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			res := newURLResource(tc.url, tc.profile, tc.loader)

			// when
			exists := res.Exists()

			// then
			require.Equal(t, tc.wantResult, exists)
		})
	}
}

func TestURLResource_Location(t *testing.T) {
	// given
	loader := property.NewYamlSourceLoader()
	resourceUrl, _ := url.Parse("http://localhost:8080/anyValidResource")
	res := newURLResource(resourceUrl, "anyProfile", loader)

	// when
	location := res.Location()

	// then
	require.Equal(t, resourceUrl.String(), location)
}

func TestURLResource_Profile(t *testing.T) {
	// given
	loader := property.NewYamlSourceLoader()
	resourceUrl, _ := url.Parse("http://localhost:8080/anyValidResource")
	res := newURLResource(resourceUrl, "anyProfile", loader)

	// when
	profile := res.Profile()

	// then
	require.Equal(t, "anyProfile", profile)
}

func TestURLResource_Reader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/anyValidResource" {
			w.Write([]byte("anyValidData"))
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	validResourceUrl, _ := url.Parse(fmt.Sprintf("%s/%s", testServer.URL, "anyValidResource"))
	nonValidResourceUrl, _ := url.Parse(fmt.Sprintf("%s/%s", testServer.URL, "anyNonValidResource"))
	unknownResourceUrl, _ := url.Parse("http://localhost:8080/anyUnknownResource")

	testCases := []struct {
		name    string
		url     *url.URL
		profile string
		loader  property.SourceLoader

		wantErr  error
		wantData []byte
	}{
		{
			name:    "resource does not exists",
			url:     nonValidResourceUrl,
			profile: "anyProfile",
			loader:  property.NewYamlSourceLoader(),
			wantErr: errors.New("resource does not exist"),
		},
		{
			name:     "valid url resource",
			url:      validResourceUrl,
			profile:  "anyProfile",
			loader:   property.NewYamlSourceLoader(),
			wantData: []byte("anyValidData"),
		},
		{
			name:    "unknown url resource",
			url:     unknownResourceUrl,
			profile: "anyProfile",
			loader:  property.NewYamlSourceLoader(),
			wantErr: &url.Error{
				Op:  "Get",
				URL: "http://localhost:8080/anyUnknownResource",
				Err: &net.OpError{
					Op:     "dial",
					Net:    "tcp",
					Source: nil,
					Addr: &net.TCPAddr{
						IP:   net.ParseIP("127.0.0.1"),
						Port: 8080,
					},
					Err: &os.SyscallError{
						Syscall: "connect",
						Err:     syscall.ECONNREFUSED,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			res := newURLResource(tc.url, tc.profile, tc.loader)

			// when
			reader, err := res.Reader()

			// then
			if tc.wantErr != nil {
				require.Nil(t, reader)
				require.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}

			data, _ := io.ReadAll(reader)
			require.Equal(t, tc.wantData, data)
		})
	}
}

func TestURLResource_PropertySourceLoader(t *testing.T) {
	// given
	loader := property.NewYamlSourceLoader()
	resourceUrl, _ := url.Parse("http://localhost:8080/anyValidResource")
	res := newURLResource(resourceUrl, "anyProfile", loader)

	// when
	result := res.PropertySourceLoader()

	// then
	require.Equal(t, loader, result)
}
