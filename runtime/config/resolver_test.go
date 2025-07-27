package config

import (
	"codnect.io/procyon/runtime/prop"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type MockResourceResolver struct {
	mock.Mock
}

func (m *MockResourceResolver) ResolveResources(ctx context.Context, location string, profiles []string) ([]Resource, error) {
	result := m.Called(ctx, location, profiles)

	resources := result.Get(0)
	err := result.Error(1)

	if resources == nil {
		return nil, err
	}

	return resources.([]Resource), err
}

func TestNewFileResourceResolver(t *testing.T) {
	yamlSourceLoader := prop.NewYamlSourceLoader()

	testCases := []struct {
		name    string
		loaders []prop.SourceLoader

		wantPanic error
	}{
		{
			name:      "no loaders",
			loaders:   []prop.SourceLoader{},
			wantPanic: fmt.Errorf("no loaders"),
		},
		{
			name:    "with valid loaders",
			loaders: []prop.SourceLoader{yamlSourceLoader},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewFileResourceResolver(tc.loaders)
				})
				return
			}

			resolver := NewFileResourceResolver(tc.loaders)

			// then
			require.NotNil(t, resolver)
		})
	}
}

func TestFileResourceResolver_Resolve(t *testing.T) {
	osFile := &os.File{}
	yamlSourceLoader := prop.NewYamlSourceLoader()

	testCases := []struct {
		name     string
		loaders  []prop.SourceLoader
		ctx      context.Context
		location string
		profiles []string

		osStat map[string]error
		osOpen map[string]*os.File

		wantErr       error
		wantResources []Resource
	}{
		{
			name:    "nil context",
			loaders: []prop.SourceLoader{yamlSourceLoader},
			ctx:     nil,
			wantErr: errors.New("nil context"),
		},
		{
			name:    "empty location",
			loaders: []prop.SourceLoader{yamlSourceLoader},
			ctx:     context.Background(),
			wantErr: errors.New("empty location"),
		},
		{
			name:     "file error",
			loaders:  []prop.SourceLoader{yamlSourceLoader},
			ctx:      context.Background(),
			location: "resources",
			osStat: map[string]error{
				"resources/procyon.yml":  errors.New("no file"),
				"resources/procyon.yaml": nil,
			},
			osOpen: map[string]*os.File{
				"resources/procyon.yaml": nil,
			},
			wantErr: errors.New("open error"),
		},
		{
			name:     "no profiles",
			loaders:  []prop.SourceLoader{yamlSourceLoader},
			ctx:      context.Background(),
			location: "resources",
			osStat: map[string]error{
				"resources/procyon.yml":  errors.New("no file"),
				"resources/procyon.yaml": nil,
			},
			osOpen: map[string]*os.File{
				"resources/procyon.yaml": osFile,
			},
			wantResources: []Resource{
				newFileResource("resources/procyon.yaml", osFile, yamlSourceLoader),
			},
		},
		{
			name:     "with profiles",
			loaders:  []prop.SourceLoader{yamlSourceLoader},
			ctx:      context.Background(),
			location: "resources",
			profiles: []string{"dev"},
			osStat: map[string]error{
				"resources/procyon-dev.yml":  errors.New("no file"),
				"resources/procyon-dev.yaml": nil,
			},
			osOpen: map[string]*os.File{
				"resources/procyon-dev.yaml": osFile,
			},
			wantResources: []Resource{
				newFileResource("resources/procyon-dev.yaml", osFile, yamlSourceLoader),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockFileSys := &MockFileSystem{}

			resolver := NewFileResourceResolver(tc.loaders)
			resolver.fileSystem = mockFileSys

			// os.Stat
			for filePath, err := range tc.osStat {
				mockFileSys.On("Stat", filePath).Return(nil, err)
			}

			// os.Open
			for filePath, file := range tc.osOpen {
				if file != nil {
					mockFileSys.On("Open", filePath).Return(file, nil)
				} else {
					mockFileSys.On("Open", filePath).Return(nil, errors.New("open error"))
				}
			}

			// when
			resources, err := resolver.ResolveResources(tc.ctx, tc.location, tc.profiles)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			require.Len(t, resources, len(tc.wantResources))

			for index, wantResource := range tc.wantResources {
				assert.Equal(t, wantResource.Name(), resources[index].Name())
				assert.Equal(t, wantResource.Location(), resources[index].Location())
				assert.Equal(t, wantResource.Profile(), resources[index].Profile())
				assert.Equal(t, wantResource.SourceLoader(), resources[index].SourceLoader())
			}
		})
	}
}
