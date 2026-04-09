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
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
	"codnect.io/procyon/runtime/config"
)

const (
	// DefaultProfilesProp is the property key for specifying default profiles.
	DefaultProfilesProp = "procyon.profiles.default"
	// ActiveProfilesProp is the property key for specifying active profiles.
	ActiveProfilesProp = "procyon.profiles.active"

	// DefaultConfigLocation is the default location for configuration files.
	DefaultConfigLocation = "resources/"
)

// Environment represents the default implementation of the Environment interface.
type Environment struct {
	activeProfiles  map[string]struct{}
	defaultProfiles map[string]struct{}

	propertySources     *config.PropertySources
	propertyResolver    config.PropertyResolver
	activeProfilesOnce  sync.Once
	defaultProfilesOnce sync.Once
	resolverOnce        sync.Once
	mu                  sync.RWMutex
}

// NewEnvironment creates a new Environment.
func NewEnvironment() *Environment {
	return &Environment{
		activeProfiles: map[string]struct{}{},
		defaultProfiles: map[string]struct{}{
			"default": {},
		},
		activeProfilesOnce:  sync.Once{},
		defaultProfilesOnce: sync.Once{},
		propertySources:     config.NewPropertySources(),
		mu:                  sync.RWMutex{},
	}
}

// doGetActiveProfiles retrieves the active profiles from the property resolver.
// IMPORTANT: Must be called WITHOUT holding e.mu to avoid deadlock.
func (e *Environment) doGetActiveProfiles() {
	e.activeProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Lookup(ActiveProfilesProp)

		if ok {
			activeProfiles := strings.Split(strings.TrimSpace(propertyValue.(string)), ",")
			err := e.setActiveProfilesWithLock(activeProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

// ActiveProfiles returns the active profiles.
func (e *Environment) ActiveProfiles() []string {
	e.doGetActiveProfiles()

	e.mu.RLock()
	defer e.mu.RUnlock()

	activeProfiles := make([]string, 0, len(e.activeProfiles))
	for profile := range e.activeProfiles {
		activeProfiles = append(activeProfiles, profile)
	}

	return activeProfiles
}

// doGetDefaultProfiles retrieves the default profiles from the property resolver.
// IMPORTANT: Must be called WITHOUT holding e.mu to avoid deadlock.
func (e *Environment) doGetDefaultProfiles() {
	e.defaultProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Lookup(DefaultProfilesProp)

		if ok {
			defaultProfiles := strings.Split(strings.TrimSpace(propertyValue.(string)), ",")
			err := e.setDefaultProfilesWithLock(defaultProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

// DefaultProfiles returns the default profiles.
func (e *Environment) DefaultProfiles() []string {
	e.doGetDefaultProfiles()

	e.mu.RLock()
	defer e.mu.RUnlock()

	profiles := make([]string, 0, len(e.defaultProfiles))
	for profile := range e.defaultProfiles {
		profiles = append(profiles, profile)
	}

	return profiles
}

// IsProfileActive checks if a profile is active.
func (e *Environment) IsProfileActive(profile string) bool {
	e.doGetActiveProfiles()

	e.mu.RLock()
	defer e.mu.RUnlock()

	_, ok := e.activeProfiles[profile]
	return ok
}

// setActiveProfilesWithLock clears and sets the active profiles.
func (e *Environment) setActiveProfilesWithLock(profiles ...string) error {
	for _, profile := range profiles {
		if err := validateProfile(profile); err != nil {
			return err
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Clear existing active profiles.
	for profile := range e.activeProfiles {
		delete(e.activeProfiles, profile)
	}

	for _, profile := range profiles {
		e.activeProfiles[profile] = struct{}{}
	}

	return nil
}

// SetActiveProfiles sets the active profiles.
func (e *Environment) SetActiveProfiles(profiles ...string) error {
	e.doGetActiveProfiles()
	return e.setActiveProfilesWithLock(profiles...)
}

// AddActiveProfiles adds active profiles.
func (e *Environment) AddActiveProfiles(profiles ...string) error {
	for _, profile := range profiles {
		if err := validateProfile(profile); err != nil {
			return err
		}
	}

	e.doGetActiveProfiles()

	e.mu.Lock()
	defer e.mu.Unlock()

	for _, profile := range profiles {
		e.activeProfiles[profile] = struct{}{}
	}

	return nil
}

// setDefaultProfilesWithLock clears and sets the default profiles.
func (e *Environment) setDefaultProfilesWithLock(profiles ...string) error {
	for _, profile := range profiles {
		if err := validateProfile(profile); err != nil {
			return err
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Clear existing default profiles.
	for profile := range e.defaultProfiles {
		delete(e.defaultProfiles, profile)
	}

	for _, profile := range profiles {
		e.defaultProfiles[profile] = struct{}{}
	}

	return nil
}

// SetDefaultProfiles sets the default profiles.
func (e *Environment) SetDefaultProfiles(profiles ...string) error {
	e.doGetDefaultProfiles()
	return e.setDefaultProfilesWithLock(profiles...)
}

// PropertySources returns the property sources.
func (e *Environment) PropertySources() *config.PropertySources {
	return e.propertySources
}

// PropertyResolver returns the property resolver.
func (e *Environment) PropertyResolver() config.PropertyResolver {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.propertyResolver == nil {
		e.propertyResolver = config.NewDefaultPropertyResolver(e.propertySources)
	}

	return e.propertyResolver
}

// validateProfile validates a profile name
func validateProfile(profile string) error {
	if strings.TrimSpace(profile) == "" {
		return fmt.Errorf("'%s' is a invalid profile", profile)
	}

	return nil
}

// configEnvCustomizer customizes the environment by loading configuration files.
// It tracks previously loaded config sources to prevent duplicate loading.
type configEnvCustomizer struct {
	loadedConfig []string
}

func newConfigEnvCustomizer() *configEnvCustomizer {
	return &configEnvCustomizer{
		loadedConfig: make([]string, 0),
	}
}

// CustomizeEnvironment loads configuration files and applies them to the environment.
// It first loads the base (profile-independent) config, resolves the default profiles,
// then loads profile-specific config files on top of the base configuration.
func (c *configEnvCustomizer) CustomizeEnvironment(env runtime.Environment, app runtime.Application) error {
	clear(c.loadedConfig)

	resResolver := app.ResourceResolver()
	loaders, err := loadPropSourceLoaders()
	if err != nil {
		return err
	}

	dataLoader := config.NewStandardDataLoader(resResolver, loaders...)
	propSources := config.NewPropertySources()
	resPropResolver := config.NewDefaultPropertyResolver(propSources)

	// Load base configuration without any profile filtering.
	err = c.loadConfig(propSources, dataLoader)
	if err != nil {
		return err
	}

	// Resolve which profiles should be active based on environment and base config.
	defaultProfiles := c.getDefaultProfiles(env, resPropResolver)

	profiles := make([]string, 0)
	profiles = append(profiles, defaultProfiles...)

	// Load profile-specific configuration files on top of the base config.
	err = c.loadConfig(propSources, dataLoader, profiles...)
	if err != nil {
		return err
	}

	return nil
}

// getDefaultProfiles determines the default profiles to be used based on the environment and the loaded configuration.
func (c *configEnvCustomizer) getDefaultProfiles(env runtime.Environment, resPropResolver config.PropertyResolver) []string {
	envPropVal, envPropExists := env.PropertyResolver().Lookup(DefaultProfilesProp)
	if envPropExists {
		envPropProfiles := strings.Split(strings.TrimSpace(envPropVal.(string)), ",")
		return envPropProfiles
	}

	envProfiles := env.DefaultProfiles()
	if len(envProfiles) != 0 && !slices.Equal(envProfiles, []string{"default"}) {
		return envProfiles
	}

	resPropVal, resPropExists := resPropResolver.Lookup(DefaultProfilesProp)
	if resPropExists {
		resPropProfiles := strings.Split(strings.TrimSpace(resPropVal.(string)), ",")
		return resPropProfiles
	}

	return []string{"default"}
}

// loadConfig loads configuration data from the specified location and profiles, and
// adds the property sources to the environment.
func (c *configEnvCustomizer) loadConfig(propSources *config.PropertySources, dataLoader config.DataLoader, profiles ...string) error {
	cfgData, err := dataLoader.Load(context.Background(), DefaultConfigLocation, profiles...)
	if err != nil {
		return err
	}

	for _, data := range cfgData {
		if slices.Contains(c.loadedConfig, data.PropertySource().Origin()) {
			continue
		}

		if err = checkForInvalidProps(data); err != nil && len(profiles) != 0 {
			return err
		}

		c.loadedConfig = append(c.loadedConfig, data.PropertySource().Origin())
		propSources.PushBack(data.PropertySource())
	}

	return nil
}

// checkForInvalidProps checks if the configuration data contains properties that are not allowed in
// profile-specific config files.
func checkForInvalidProps(data config.Data) error {
	propSource := data.PropertySource()

	if _, ok := propSource.Value(DefaultProfilesProp); ok {
		return fmt.Errorf("%s property cannot be set in profile specific config file", DefaultProfilesProp)
	}

	if _, ok := propSource.Value(ActiveProfilesProp); ok {
		return fmt.Errorf("%s property cannot be set in profile specific config file", ActiveProfilesProp)
	}

	return nil
}

// loadPropSourceLoaders loads all registered PropertySourceLoader components and returns them as a slice.
func loadPropSourceLoaders() ([]config.PropertySourceLoader, error) {
	loaders := make([]config.PropertySourceLoader, 0)
	components := component.ListOf[config.PropertySourceLoader]()

	for _, comp := range components {
		loader, err := component.Load[config.PropertySourceLoader](comp.Definition().Name())
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, loader)
	}

	return loaders, nil
}
