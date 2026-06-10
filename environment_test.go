// Copyright 2026 Codnect
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

package procyon

import (
	"errors"
	stdio "io"
	"io/fs"
	"testing"

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
	"codnect.io/procyon/runtime/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type FakeFile struct {
	contents string
	offset   int
	readErr  error
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
	if f.readErr != nil {
		return 0, f.readErr
	}
	if f.offset >= len(f.contents) {
		return 0, stdio.EOF
	}
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

func TestEnvironment_ActiveProfiles(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(env *Environment)

		wantProfiles []string
		wantPanic    error
	}{
		{
			name: "valid profiles",
			preCondition: func(env *Environment) {
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"procyon.profiles.active": "dev",
				})

				env.PropertySources().PushBack(propSource)
			},
			wantProfiles: []string{"dev"},
		},
		{
			name: "invalid profiles",
			preCondition: func(env *Environment) {
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"procyon.profiles.active": " ",
				})

				env.PropertySources().PushBack(propSource)
			},
			wantPanic: errors.New("invalid profile: empty or blank profile"),
		},
		{
			name: "programmatically set profiles take precedence over property source",
			preCondition: func(env *Environment) {
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"procyon.profiles.active": "dev",
				})
				env.PropertySources().PushBack(propSource)

				if err := env.SetActiveProfiles("prod"); err != nil {
					panic(err)
				}
			},
			wantProfiles: []string{"prod"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()

			if tc.preCondition != nil {
				tc.preCondition(env)
			}

			// when
			if tc.wantPanic != nil {
				require.PanicsWithError(t, tc.wantPanic.Error(), func() {
					env.ActiveProfiles()
				})
				return
			}

			profiles := env.ActiveProfiles()

			// then
			assert.ElementsMatch(t, tc.wantProfiles, profiles)
		})
	}
}

func TestEnvironment_DefaultProfiles(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(env *Environment)

		wantProfiles []string
		wantPanic    error
	}{
		{
			name: "valid profiles",
			preCondition: func(env *Environment) {
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"procyon.profiles.default": "dev",
				})

				env.PropertySources().PushBack(propSource)
			},
			wantProfiles: []string{"dev"},
		},
		{
			name: "invalid profiles",
			preCondition: func(env *Environment) {
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"procyon.profiles.default": " ",
				})

				env.PropertySources().PushBack(propSource)
			},
			wantPanic: errors.New("invalid profile: empty or blank profile"),
		},
		{
			name: "programmatically set profiles take precedence over property source",
			preCondition: func(env *Environment) {
				propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
					"procyon.profiles.default": "dev",
				})
				env.PropertySources().PushBack(propSource)

				if err := env.SetDefaultProfiles("prod"); err != nil {
					panic(err)
				}
			},
			wantProfiles: []string{"prod"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()

			if tc.preCondition != nil {
				tc.preCondition(env)
			}

			// when
			if tc.wantPanic != nil {
				require.PanicsWithError(t, tc.wantPanic.Error(), func() {
					env.DefaultProfiles()
				})
				return
			}

			profiles := env.DefaultProfiles()

			// then
			assert.ElementsMatch(t, tc.wantProfiles, profiles)
		})
	}
}

func TestEnvironment_IsProfileActive(t *testing.T) {
	testCases := []struct {
		name           string
		activeProfiles string
		profileToCheck string
		wantResult     bool
	}{
		{
			name:           "Profile is active",
			activeProfiles: "dev, test",
			profileToCheck: "dev",
			wantResult:     true,
		},
		{
			name:           "Profile is not active",
			activeProfiles: "dev",
			profileToCheck: "test",
			wantResult:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()

			propSource := config.NewMapPropertySource("anyMapSource", map[string]any{
				"procyon.profiles.active": tc.activeProfiles,
			})

			env.PropertySources().PushBack(propSource)

			// when
			isActive := env.IsProfileActive(tc.profileToCheck)

			// then
			assert.Equal(t, tc.wantResult, isActive)
		})
	}
}

func TestEnvironment_SetActiveProfiles(t *testing.T) {
	testCases := []struct {
		name           string
		preCondition   func(env *Environment)
		activeProfiles []string

		wantProfiles []string
		wantErr      error
	}{
		{
			name:           "empty profile",
			activeProfiles: []string{""},
			wantErr:        errors.New("invalid profile: empty or blank profile"),
		},
		{
			name:           "blank profile",
			activeProfiles: []string{" "},
			wantErr:        errors.New("invalid profile: empty or blank profile"),
		},
		{
			name:           "Set single active profile",
			activeProfiles: []string{"dev"},
			wantProfiles:   []string{"dev"},
		},
		{
			name:           "Set multiple active profiles",
			activeProfiles: []string{"dev", "test"},
			wantProfiles:   []string{"dev", "test"},
		},
		{
			name: "override existing active profiles",
			preCondition: func(env *Environment) {
				err := env.SetActiveProfiles("dev", "test")
				require.NoError(t, err)
			},
			activeProfiles: []string{"prod"},
			wantProfiles:   []string{"prod"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()

			if tc.preCondition != nil {
				tc.preCondition(env)
			}

			// when
			err := env.SetActiveProfiles(tc.activeProfiles...)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			profiles := env.ActiveProfiles()
			assert.ElementsMatch(t, tc.wantProfiles, profiles)
		})
	}
}

func TestEnvironment_AddActiveProfiles(t *testing.T) {
	testCases := []struct {
		name            string
		initialProfiles []string
		profilesToAdd   []string

		wantProfiles []string
		wantErr      error
	}{
		{
			name:          "empty profile",
			profilesToAdd: []string{""},
			wantErr:       errors.New("invalid profile: empty or blank profile"),
		},
		{
			name:          "blank profile",
			profilesToAdd: []string{" "},
			wantErr:       errors.New("invalid profile: empty or blank profile"),
		},
		{
			name:            "Add single active profile",
			initialProfiles: []string{"dev"},
			profilesToAdd:   []string{"test"},
			wantProfiles:    []string{"dev", "test"},
		},
		{
			name:            "Add multiple active profiles",
			initialProfiles: []string{"dev"},
			profilesToAdd:   []string{"test", "secure"},
			wantProfiles:    []string{"dev", "test", "secure"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()
			err := env.SetActiveProfiles(tc.initialProfiles...)
			assert.NoError(t, err)

			// when
			err = env.AddActiveProfiles(tc.profilesToAdd...)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			profiles := env.ActiveProfiles()
			assert.ElementsMatch(t, tc.wantProfiles, profiles)
		})
	}
}

func TestEnvironment_SetDefaultProfiles(t *testing.T) {
	testCases := []struct {
		name            string
		preCondition    func(env *Environment)
		defaultProfiles []string

		wantProfiles []string
		wantErr      error
	}{
		{
			name:            "empty profile",
			defaultProfiles: []string{""},
			wantErr:         errors.New("invalid profile: empty or blank profile"),
		},
		{
			name:            "blank profile",
			defaultProfiles: []string{" "},
			wantErr:         errors.New("invalid profile: empty or blank profile"),
		},
		{
			name:            "Set single default profile",
			defaultProfiles: []string{"dev"},
			wantProfiles:    []string{"dev"},
		},
		{
			name:            "Set multiple default profiles",
			defaultProfiles: []string{"dev", "test"},
			wantProfiles:    []string{"dev", "test"},
		},
		{
			name: "override existing DEFAULT profiles",
			preCondition: func(env *Environment) {
				err := env.setDefaultProfiles("dev", "test")
				require.NoError(t, err)
			},
			defaultProfiles: []string{"prod"},
			wantProfiles:    []string{"prod"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			env := NewEnvironment()

			// when
			err := env.SetDefaultProfiles(tc.defaultProfiles...)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			profiles := env.DefaultProfiles()
			assert.ElementsMatch(t, tc.wantProfiles, profiles)
		})
	}
}

func TestEnvironment_PropertySources(t *testing.T) {
	// given
	env := NewEnvironment()

	// when
	propertySources := env.PropertySources()

	// then
	assert.NotNil(t, propertySources)
	assert.Equal(t, env.propertySources, propertySources)
}

func TestEnvironment_PropertyResolver(t *testing.T) {
	// given
	env := NewEnvironment()

	// when
	propertyResolver := env.PropertyResolver()

	// then
	assert.NotNil(t, propertyResolver)
	assert.Equal(t, env.propertyResolver, propertyResolver)
}

func TestConfigEnvCustomizer_CustomizeEnvironment(t *testing.T) {
	testCases := []struct {
		name         string
		preCondition func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application)
		env          runtime.Environment
		app          runtime.Application

		wantProperties map[string]any
		wantErr        error
	}{
		{
			name:    "nil environment",
			env:     nil,
			wantErr: errors.New("nil environment"),
		},
		{
			name:    "nil application",
			env:     NewEnvironment(),
			app:     nil,
			wantErr: errors.New("nil application"),
		},
		{
			name: "property source loader load error",
			env:  NewEnvironment(),
			app:  New(),
			preCondition: func(customizer *configEnvCustomizer, _ runtime.Environment, _ runtime.Application) {
				def, err := component.MakeDefinition(func() config.PropertySourceLoader {
					return nil
				})
				require.NoError(t, err)

				comp := component.Create(def)
				customizer.propSourceLoaders = []*component.Component{comp}
			},
			wantErr: errors.New("customize environment: load property source loader: load component \"propertySourceLoader\": not found"),
		},
		{
			name: "load config error",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				resourceResolver := &AnyMockResourceResolver{}
				resourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(nil, errors.New("resolve config file error"))

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(resourceResolver)
			},
			wantErr: errors.New("customize environment: config location resources/: load config for default profile: resolve config file error"),
		},
		{
			name: "procyon.profiles.default not allowed in custom profile config",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(noResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    default: another",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				err := env.SetDefaultProfiles("custom")
				require.NoError(t, err)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: config file resources/procyon-custom.yaml: procyon.profiles.default not allowed in profile-specific config"),
		},
		{
			name: "invalid default profile in default config file",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    default: \",test\"",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: invalid profile: empty or blank profile"),
		},
		{
			name: "invalid default profile from environment variable",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				envPropSource := config.NewMapPropertySource("envProps", map[string]any{
					"procyon.profiles.default": ",test",
				})
				env.PropertySources().PushBack(envPropSource)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: invalid profile: empty or blank profile"),
		},
		{
			name: "invalid default profile type in default config file",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    default: 123",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: invalid profile: procyon.profiles.default must be a string, got int"),
		},
		{
			name: "invalid active profile in default config file",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    active: \",test\"",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: invalid profile: empty or blank profile"),
		},
		{
			name: "invalid active profile type in default config file",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    active: 123",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: invalid profile: procyon.profiles.active must be a string, got int"),
		},
		{
			name: "invalid active profile from environment variable",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				envPropSource := config.NewMapPropertySource("envProps", map[string]any{
					"procyon.profiles.active": ",test",
				})
				env.PropertySources().PushBack(envPropSource)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: invalid profile: empty or blank profile"),
		},
		{
			name: "procyon.profiles.active not allowed in custom profile config",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(noResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    active: another",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				err := env.SetDefaultProfiles("custom")
				require.NoError(t, err)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantErr: errors.New("customize environment: config file resources/procyon-custom.yaml: procyon.profiles.active not allowed in profile-specific config"),
		},
		{
			name: "default profile: custom profile config overrides base config (programmatically set)",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "app:\n  name: demo\n  version: 1.0.0",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "app:\n  name: demo-custom",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				err := env.SetDefaultProfiles("custom")
				require.NoError(t, err)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantProperties: map[string]any{
				"app.name":    "demo-custom",
				"app.version": "1.0.0",
			},
		},
		{
			name: "default profile: custom profile config overrides base config (from environment variable)",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "app:\n  name: demo\n  version: 1.0.0",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "app:\n  name: demo-custom",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				envPropSource := config.NewMapPropertySource("envProps", map[string]any{
					"procyon.profiles.default": "custom",
				})
				env.PropertySources().PushBack(envPropSource)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantProperties: map[string]any{
				"app.name":    "demo-custom",
				"app.version": "1.0.0",
			},
		},
		{
			name: "default profile: custom profile config overrides base config (from base config)",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    default: custom\napp:\n  name: demo\n  version: 1.0.0",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "app:\n  name: demo-custom",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantProperties: map[string]any{
				"app.name":    "demo-custom",
				"app.version": "1.0.0",
			},
		},
		{
			name: "active profile: custom profile config overrides base config (programmatically set)",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "app:\n  name: demo\n  version: 1.0.0",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "app:\n  name: demo-custom",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				err := env.SetActiveProfiles("custom")
				require.NoError(t, err)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantProperties: map[string]any{
				"app.name":    "demo-custom",
				"app.version": "1.0.0",
			},
		},
		{
			name: "active profile: custom profile config overrides base config (from environment variable)",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "app:\n  name: demo\n  version: 1.0.0",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "app:\n  name: demo-custom",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				envPropSource := config.NewMapPropertySource("envProps", map[string]any{
					"procyon.profiles.active": "custom",
				})
				env.PropertySources().PushBack(envPropSource)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantProperties: map[string]any{
				"app.name":    "demo-custom",
				"app.version": "1.0.0",
			},
		},
		{
			name: "active profile: custom profile config overrides base config (from base config)",
			env:  NewEnvironment(),
			app:  &AnyMockApplication{},
			preCondition: func(customizer *configEnvCustomizer, env runtime.Environment, app runtime.Application) {
				noResource := &AnyMockResource{}
				noResource.On("Exists").Return(false)

				// resources/procyon.yaml
				defaultPropFile := &FakeFile{
					contents: "procyon:\n  profiles:\n    active: custom\napp:\n  name: demo\n  version: 1.0.0",
				}

				defaultPropResource := &AnyMockResource{}
				defaultPropResource.On("Exists").Return(true)
				defaultPropResource.On("Reader").Return(defaultPropFile, nil)

				propResourceResolver := &AnyMockResourceResolver{}
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yaml").
					Return(defaultPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon.yml").
					Return(noResource, nil)

				// resources/procyon-custom.yaml
				customPropFile := &FakeFile{
					contents: "app:\n  name: demo-custom",
				}

				customPropResource := &AnyMockResource{}
				customPropResource.On("Exists").Return(true)
				customPropResource.On("Reader").Return(customPropFile, nil)

				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yaml").
					Return(customPropResource, nil)
				propResourceResolver.On("Resolve", mock.AnythingOfType("context.backgroundCtx"), "resources/procyon-custom.yml").
					Return(noResource, nil)

				mockApp := app.(*AnyMockApplication)
				mockApp.On("ResourceResolver").Return(propResourceResolver)
			},
			wantProperties: map[string]any{
				"app.name":    "demo-custom",
				"app.version": "1.0.0",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			customizer := newConfigEnvCustomizer()

			if tc.preCondition != nil {
				tc.preCondition(customizer, tc.env, tc.app)
			}

			// when
			err := customizer.CustomizeEnvironment(tc.env, tc.app)

			// then
			if tc.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NoError(t, err)

			for wantKey, wantValue := range tc.wantProperties {
				value, exists := tc.env.PropertyResolver().Lookup(wantKey)
				assert.True(t, exists)
				assert.Equal(t, wantValue, value)
			}
		})
	}
}
