package procyon

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/runtime/env"
	"codnect.io/procyon-core/runtime/env/property"
	"fmt"
	"strings"
	"sync"
)

type Environment struct {
	activeProfiles  map[string]struct{}
	defaultProfiles map[string]struct{}

	sources             *property.Sources
	resolver            property.Resolver
	activeProfilesOnce  sync.Once
	defaultProfilesOnce sync.Once
	resolverOnce        sync.Once
	mu                  sync.RWMutex
}

func newEnvironment() *Environment {
	return &Environment{
		activeProfiles: map[string]struct{}{},
		defaultProfiles: map[string]struct{}{
			"default": {},
		},
		sources:             property.NewPropertySources(),
		activeProfilesOnce:  sync.Once{},
		defaultProfilesOnce: sync.Once{},
		mu:                  sync.RWMutex{},
	}
}

func (e *Environment) customize(container container.Container) error {
	customizers, err := getComponentsByType[env.Customizer](container)
	if err != nil {
		return err
	}

	for _, customizer := range customizers {
		err = customizer.CustomizeEnvironment(e)

		if err != nil {
			return err
		}
	}

	return err
}

func (e *Environment) validateProfile(profile string) error {
	if strings.TrimSpace(profile) == "" {
		return fmt.Errorf("env: `%s` is a invalid profile", profile)
	}

	return nil
}

func (e *Environment) doGetActiveProfiles() {
	e.activeProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Property("procyon.profiles.active")

		if ok {
			activeProfiles := strings.Split(strings.TrimSpace(propertyValue), ",")
			err := e.SetActiveProfiles(activeProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

func (e *Environment) ActiveProfiles() []string {
	e.doGetActiveProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	activeProfiles := make([]string, 0)
	for profile := range e.activeProfiles {
		activeProfiles = append(activeProfiles, profile)
	}

	return activeProfiles
}

func (e *Environment) doGetDefaultProfiles() {
	e.defaultProfilesOnce.Do(func() {
		propertyValue, ok := e.PropertyResolver().Property("procyon.profiles.default")

		if ok {
			defaultProfiles := strings.Split(strings.TrimSpace(propertyValue), ",")
			err := e.SetDefaultProfiles(defaultProfiles...)

			if err != nil {
				panic(err)
			}
		}
	})
}

func (e *Environment) DefaultProfiles() []string {
	e.doGetDefaultProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	profiles := make([]string, 0)
	for profile := range e.defaultProfiles {
		profiles = append(profiles, profile)
	}

	return profiles
}

func (e *Environment) IsProfileActive(profile string) bool {
	defer e.mu.Unlock()
	e.mu.Lock()

	if _, ok := e.activeProfiles[profile]; ok {
		return true
	}

	return false
}

func (e *Environment) clearActiveProfiles() {
	defer e.mu.Unlock()
	e.mu.Lock()

	for profile := range e.activeProfiles {
		delete(e.activeProfiles, profile)
	}
}

func (e *Environment) SetActiveProfiles(profiles ...string) error {
	e.clearActiveProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	for _, profile := range profiles {
		err := e.validateProfile(profile)
		if err != nil {
			return err
		}

		e.activeProfiles[profile] = struct{}{}
	}

	return nil
}
func (e *Environment) AddActiveProfile(profiles ...string) error {
	for _, profile := range profiles {
		err := e.validateProfile(profile)

		if err != nil {
			return err
		}
	}

	e.doGetActiveProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	for _, profile := range profiles {
		e.activeProfiles[profile] = struct{}{}
	}

	return nil
}

func (e *Environment) clearDefaultProfiles() {
	defer e.mu.Unlock()
	e.mu.Lock()

	for profile := range e.defaultProfiles {
		delete(e.defaultProfiles, profile)
	}
}

func (e *Environment) SetDefaultProfiles(profiles ...string) error {
	e.clearDefaultProfiles()

	defer e.mu.Unlock()
	e.mu.Lock()

	for _, profile := range profiles {
		err := e.validateProfile(profile)
		if err != nil {
			return err
		}

		e.defaultProfiles[profile] = struct{}{}
	}

	return nil
}

func (e *Environment) Merge(parent env.Environment) {
	parentPropertySources := parent.PropertySources().ToSlice()

	if len(parentPropertySources) != 0 {
		for _, propertySource := range parentPropertySources {
			if !e.sources.Contains(propertySource.Name()) {
				e.sources.AddLast(propertySource)
			}
		}
	}

	defer e.mu.Unlock()
	e.mu.Lock()

	parentActiveProfiles := parent.ActiveProfiles()
	if len(parentActiveProfiles) != 0 {
		for _, profile := range parentActiveProfiles {
			e.activeProfiles[profile] = struct{}{}
		}
	}

	parentDefaultProfiles := parent.DefaultProfiles()
	if len(parentDefaultProfiles) != 0 {
		for _, profile := range parentDefaultProfiles {
			e.defaultProfiles[profile] = struct{}{}
		}
	}
}

func (e *Environment) Variables() env.Variables {
	return nil
}

func (e *Environment) PropertySources() *property.Sources {
	return e.sources
}

func (e *Environment) PropertyResolver() property.Resolver {
	defer e.mu.Unlock()
	e.mu.Lock()

	if e.resolver == nil {
		e.resolver = property.NewSourcesResolver(e.sources)
	}

	return e.resolver
}
