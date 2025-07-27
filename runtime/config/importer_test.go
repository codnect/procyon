package config

import (
	"codnect.io/procyon/runtime/prop"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewImporter(t *testing.T) {
	// given

	// when
	importer := NewImporter(nil, nil)

	// then
	assert.NotNil(t, importer)
}

func TestImporter_Import(t *testing.T) {
	testCases := []struct {
		name string

		loaders []Loader

		ctx      context.Context
		location string
		profiles []string

		resolveResources []Resource
		resolveErr       error

		wantErr error
	}{
		{
			name:     "no loaders",
			ctx:      context.Background(),
			location: "resources",
			profiles: []string{},
			loaders:  []Loader{},
			resolveResources: []Resource{
				newFileResource("anyFile", &FakeFile{}, prop.NewYamlSourceLoader()),
			},
			wantErr: errors.New("no loaders"),
		},
		{
			name:             "resolve error",
			ctx:              context.Background(),
			location:         "resources",
			profiles:         []string{},
			loaders:          []Loader{},
			resolveResources: nil,
			resolveErr:       errors.New("resolve error"),
			wantErr:          errors.New("resolve error"),
		},
		{
			name:     "multiple loaders",
			ctx:      context.Background(),
			location: "resources",
			profiles: []string{},
			loaders: []Loader{
				&AnyLoader{
					loadable: true,
					data:     &Data{},
				},
				&AnyLoader{
					loadable: true,
					data:     &Data{},
				},
			},
			resolveResources: []Resource{
				newFileResource("anyFile", &FakeFile{}, prop.NewYamlSourceLoader()),
			},
			wantErr: errors.New("multiple loaders"),
		},
		{
			name:     "import successfully",
			ctx:      context.Background(),
			location: "resources",
			profiles: []string{},
			loaders: []Loader{
				&AnyLoader{
					loadable: true,
					data:     &Data{},
				},
			},
			resolveResources: []Resource{
				newFileResource("anyFile", &FakeFile{}, prop.NewYamlSourceLoader()),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockResolver := &MockResourceResolver{}

			importer := NewImporter([]ResourceResolver{mockResolver}, tc.loaders)

			mockResolver.On("ResolveResources", tc.ctx, tc.location, tc.profiles).
				Return(tc.resolveResources, tc.resolveErr)

			// when
			data, err := importer.Import(tc.ctx, tc.location, tc.profiles)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, data)
		})
	}
}
