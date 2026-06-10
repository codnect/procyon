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
	"errors"
	"fmt"
	"maps"
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

var (
	// reservedDefaultProfiles contains the built-in default profile names.
	reservedDefaultProfiles = map[string]struct{}{"default": {}}
)

// Environment represents the default implementation of the Environment interface.
type Environment struct {
	activeProfiles  map[string]struct{}
	defaultProfiles map[string]struct{}

	propertySources  *config.PropertySources
	propertyResolver config.PropertyResolver

	mu sync.RWMutex
}

// NewEnvironment creates a new Environment.
func NewEnvironment() *Environment {
	propSources := config.NewPropertySources()
	propResolver := config.NewDefaultPropertyResolver(propSources)

	return &Environment{
		activeProfiles: map[string]struct{}{},
		defaultProfiles: map[string]struct{}{
			"default": {},
		},
		propertySources:  propSources,
		propertyResolver: propResolver,
	}
}

// doGetActiveProfiles resolves active profiles from the property source
// if they have not been set programmatically yet. Caller must hold e.mu.
func (e *Environment) doGetActiveProfiles() {
	if len(e.activeProfiles) != 0 {
		return
	}

	propertyValue, ok := e.PropertyResolver().Lookup(ActiveProfilesProp)
	if !ok {
		return
	}

	strPropVal, isString := propertyValue.(string)
	if !isString {
		return
	}

	profiles := strings.Split(strings.TrimSpace(strPropVal), ",")
	if err := e.setActiveProfiles(profiles...); err != nil {
		panic(err)
	}
}

// ActiveProfiles returns the active profiles.
func (e *Environment) ActiveProfiles() []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.doGetActiveProfiles()

	activeProfiles := make([]string, 0, len(e.activeProfiles))
	for profile := range e.activeProfiles {
		activeProfiles = append(activeProfiles, profile)
	}
	return activeProfiles
}

// doGetDefaultProfiles resolves default profiles from the property source
// if they are still at their reserved default. Caller must hold e.mu.
func (e *Environment) doGetDefaultProfiles() {
	if !maps.Equal(e.defaultProfiles, reservedDefaultProfiles) {
		return
	}

	propertyValue, ok := e.PropertyResolver().Lookup(DefaultProfilesProp)
	if !ok {
		return
	}

	strPropVal, isString := propertyValue.(string)
	if !isString {
		return
	}

	profiles := strings.Split(strings.TrimSpace(strPropVal), ",")
	if err := e.setDefaultProfiles(profiles...); err != nil {
		panic(err)
	}
}

// DefaultProfiles returns the default profiles.
func (e *Environment) DefaultProfiles() []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.doGetDefaultProfiles()

	profiles := make([]string, 0, len(e.defaultProfiles))
	for profile := range e.defaultProfiles {
		profiles = append(profiles, profile)
	}
	return profiles
}

// IsProfileActive checks if a profile is active.
func (e *Environment) IsProfileActive(profile string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.doGetActiveProfiles()

	_, ok := e.activeProfiles[profile]
	return ok
}

// SetActiveProfiles sets the active profiles.
func (e *Environment) SetActiveProfiles(profiles ...string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.setActiveProfiles(profiles...)
}

// AddActiveProfiles adds active profiles.
func (e *Environment) AddActiveProfiles(profiles ...string) error {
	if err := validateProfiles(profiles...); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.doGetActiveProfiles()

	for _, profile := range profiles {
		e.activeProfiles[profile] = struct{}{}
	}
	return nil
}

// SetDefaultProfiles sets the default profiles.
func (e *Environment) SetDefaultProfiles(profiles ...string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.setDefaultProfiles(profiles...)
}

// PropertySources returns the property sources.
func (e *Environment) PropertySources() *config.PropertySources {
	return e.propertySources
}

// PropertyResolver returns the property resolver.
func (e *Environment) PropertyResolver() config.PropertyResolver {
	return e.propertyResolver
}

// setActiveProfiles clears and sets the active profiles.
func (e *Environment) setActiveProfiles(profiles ...string) error {
	if err := validateProfiles(profiles...); err != nil {
		return err
	}

	for profile := range e.activeProfiles {
		delete(e.activeProfiles, profile)
	}

	for _, profile := range profiles {
		e.activeProfiles[profile] = struct{}{}
	}

	return nil
}

// setDefaultProfiles clears and sets the default profiles.
func (e *Environment) setDefaultProfiles(profiles ...string) error {
	if err := validateProfiles(profiles...); err != nil {
		return err
	}

	for profile := range e.defaultProfiles {
		delete(e.defaultProfiles, profile)
	}

	for _, profile := range profiles {
		e.defaultProfiles[profile] = struct{}{}
	}

	return nil
}

// validateProfiles validates all given profile names.
func validateProfiles(profiles ...string) error {
	for _, profile := range profiles {
		if err := validateProfile(profile); err != nil {
			return err
		}
	}

	return nil
}

// validateProfile validates a single profile name.
func validateProfile(profile string) error {
	if strings.TrimSpace(profile) == "" {
		return errors.New("invalid profile: empty or blank profile")
	}

	return nil
}

// configEnvCustomizer customizes the environment by loading configuration files.
// It tracks previously loaded config sources to prevent duplicate loading.
type configEnvCustomizer struct {
	propSourceLoaders        []*component.Component
	propSourceLoaderLoadFunc func(name string) (config.PropertySourceLoader, error)
}

func newConfigEnvCustomizer() *configEnvCustomizer {
	return &configEnvCustomizer{
		propSourceLoaders: component.ListOf[config.PropertySourceLoader](),
		propSourceLoaderLoadFunc: func(name string) (config.PropertySourceLoader, error) {
			return component.Load[config.PropertySourceLoader](name)
		},
	}
}

// CustomizeEnvironment loads configuration files and applies them to the environment.
// It first loads the base (profile-independent) config, resolves the default profiles,
// then loads profile-specific config files on top of the base configuration.
func (c *configEnvCustomizer) CustomizeEnvironment(env runtime.Environment, app runtime.Application) error {
	if env == nil {
		return errors.New("nil environment")
	}

	if app == nil {
		return errors.New("nil application")
	}

	if err := c.customizeEnvironment(env, app); err != nil {
		return fmt.Errorf("customize environment: %w", err)
	}
	return nil
}

func (c *configEnvCustomizer) customizeEnvironment(env runtime.Environment, app runtime.Application) error {
	loaders, err := c.loadPropSourceLoaders()
	if err != nil {
		return err
	}

	dataLoader := config.NewStandardDataLoader(app.ResourceResolver(), loaders...)
	propSources := config.NewPropertySources()
	resPropResolver := config.NewDefaultPropertyResolver(propSources)

	// Load base configuration without any profile filtering.
	err = c.loadConfig(propSources, dataLoader)
	if err != nil {
		return err
	}

	var (
		defaultProfiles []string
		activeProfiles  []string
	)

	// Resolve which profiles should be active based on environment and base config.
	defaultProfiles, err = c.getDefaultProfiles(env, resPropResolver)
	if err != nil {
		return err
	}

	activeProfiles, err = c.getActiveProfiles(env, resPropResolver)
	if err != nil {
		return err
	}

	profiles := make([]string, 0)
	if len(activeProfiles) > 0 {
		profiles = append(profiles, activeProfiles...)
	} else {
		profiles = append(profiles, defaultProfiles...)
	}

	// Load profile-specific configuration files on top of the base config.
	err = c.loadConfig(propSources, dataLoader, profiles...)
	if err != nil {
		return err
	}

	return c.applyToEnvironment(env, defaultProfiles, activeProfiles, propSources)
}

// loadPropSourceLoaders loads all registered PropertySourceLoader components and returns them as a slice.
func (c *configEnvCustomizer) loadPropSourceLoaders() ([]config.PropertySourceLoader, error) {
	loaders := make([]config.PropertySourceLoader, 0)

	for _, comp := range c.propSourceLoaders {
		loader, err := c.propSourceLoaderLoadFunc(comp.Definition().Name())
		if err != nil {
			return nil, fmt.Errorf("load property source loader: %w", err)
		}

		loaders = append(loaders, loader)
	}

	return loaders, nil
}

// getDefaultProfiles determines the default profiles to be used based on the environment and the loaded configuration.
func (c *configEnvCustomizer) getDefaultProfiles(env runtime.Environment, resPropResolver config.PropertyResolver) ([]string, error) {
	profiles, err := lookupProfilesProp(env.PropertyResolver(), DefaultProfilesProp)
	if err != nil {
		return nil, err
	} else if len(profiles) > 0 {
		return profiles, nil
	}

	profiles = env.DefaultProfiles()
	if len(profiles) != 0 && !slices.Equal(profiles, []string{"default"}) {
		return profiles, nil
	}

	profiles, err = lookupProfilesProp(resPropResolver, DefaultProfilesProp)
	if err != nil {
		return nil, err
	} else if len(profiles) > 0 {
		return profiles, nil
	}

	return []string{"default"}, nil
}

// getActiveProfiles determines the active profiles to be used based on the environment and the loaded configuration.
func (c *configEnvCustomizer) getActiveProfiles(env runtime.Environment, resPropResolver config.PropertyResolver) ([]string, error) {
	profiles, err := lookupProfilesProp(env.PropertyResolver(), ActiveProfilesProp)
	if err != nil {
		return nil, err
	} else if len(profiles) > 0 {
		return profiles, nil
	}

	profiles = env.ActiveProfiles()
	if len(profiles) != 0 {
		return profiles, nil
	}

	profiles, err = lookupProfilesProp(resPropResolver, ActiveProfilesProp)
	if err != nil {
		return nil, err
	} else if len(profiles) > 0 {
		return profiles, nil
	}

	return []string{}, nil
}

// loadConfig loads configuration data from the specified location and profiles, and
// adds the property sources to the environment.
func (c *configEnvCustomizer) loadConfig(propSources *config.PropertySources, dataLoader config.DataLoader, profiles ...string) error {
	cfgData, err := dataLoader.Load(context.Background(), DefaultConfigLocation, profiles...)
	if err != nil {
		return fmt.Errorf("config location %s: %w", DefaultConfigLocation, err)
	}

	for _, data := range cfgData {
		if err = checkProfileSpecificProps(data); err != nil && len(profiles) != 0 {
			return fmt.Errorf("config file %s: %w", data.PropertySource().Origin(), err)
		}

		propSources.PushFront(data.PropertySource())
	}

	return nil
}

// applyToEnvironment adds loaded property sources to the environment and applies the resolved default
// and active profiles.
func (c *configEnvCustomizer) applyToEnvironment(env runtime.Environment, defaultProfiles, activeProfiles []string, propSources *config.PropertySources) error {
	for _, propSource := range propSources.Slice() {
		env.PropertySources().PushBack(propSource)
	}

	// Profiles are already validated by getDefaultProfiles and getActiveProfiles;
	// these calls should never return an error.
	err := env.SetDefaultProfiles(defaultProfiles...)
	if err != nil {
		return err
	}

	err = env.SetActiveProfiles(activeProfiles...)
	if err != nil {
		return err
	}

	return nil
}

// checkProfileSpecificProps checks whether profile-specific configuration data contains properties
// that are only allowed in base configuration files.
func checkProfileSpecificProps(data config.Data) error {
	propSource := data.PropertySource()

	if _, ok := propSource.Value(DefaultProfilesProp); ok {
		return fmt.Errorf("%s not allowed in profile-specific config", DefaultProfilesProp)
	}

	if _, ok := propSource.Value(ActiveProfilesProp); ok {
		return fmt.Errorf("%s not allowed in profile-specific config", ActiveProfilesProp)
	}

	return nil
}

// lookupProfilesProp looks up a profile property and returns its parsed profile names.
func lookupProfilesProp(resolver config.PropertyResolver, key string) ([]string, error) {
	val, ok := resolver.Lookup(key)
	if !ok {
		return nil, nil
	}

	s, isString := val.(string)
	if !isString {
		return nil, fmt.Errorf("invalid profile: %s must be a string, got %T", key, val)
	}

	return splitProfiles(s)
}

// splitProfiles splits a comma-separated profile string into validated profile names.
func splitProfiles(s string) ([]string, error) {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if err := validateProfile(p); err != nil {
			return nil, err
		}

		out = append(out, p)
	}
	return out, nil
}
