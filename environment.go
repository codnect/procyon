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
	"fmt"
	"strings"
	"sync"

	"codnect.io/procyon/runtime/config"
)

const (
	// DefaultProfilesProp is the property key for specifying default profiles.
	DefaultProfilesProp = "procyon.profiles.default"
	// ActiveProfilesProp is the property key for specifying active profiles.
	ActiveProfilesProp = "procyon.profiles.active"
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
