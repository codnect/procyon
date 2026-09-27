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

package condition

import (
	"context"
	"fmt"

	"go.codnect.io/procyon/component"
	"go.codnect.io/procyon/runtime"
)

// Property is a condition that matches when the specified environment
// properties exist and optionally have the expected value.
type Property struct {
	prefix   string
	names    []string
	value    string
	hasValue bool
}

// OnProperty creates a condition for the specified property names.
//
// When multiple names are specified, all properties must exist.
func OnProperty(names ...string) *Property {
	return &Property{
		names: names,
	}
}

// Prefix sets the prefix applied to all property names.
func (p *Property) Prefix(prefix string) *Property {
	p.prefix = prefix
	return p
}

// Equals requires the properties to have the specified value.
func (p *Property) Equals(value string) *Property {
	p.value = value
	p.hasValue = true
	return p
}

// Matches reports whether all configured properties satisfy the condition.
func (p *Property) Matches(ctx context.Context, container component.Container) bool {
	env, err := component.ResolveType[runtime.Environment](ctx, container)
	if err != nil {
		return false
	}

	for _, name := range p.names {
		key := propertyKey(p.prefix, name)

		value, ok := env.PropertyResolver().Lookup(key)
		if !ok {
			return false
		}

		if p.hasValue && fmt.Sprint(value) != p.value {
			return false
		}
	}

	return true
}

// MissingProperty is a condition that matches when the specified environment
// properties do not exist.
type MissingProperty struct {
	prefix string
	names  []string
}

// OnMissingProperty creates a condition for the specified missing property names.
//
// When multiple names are specified, all properties must be missing.
func OnMissingProperty(names ...string) *MissingProperty {
	return &MissingProperty{
		names: names,
	}
}

// Prefix sets the prefix applied to all property names.
func (p *MissingProperty) Prefix(prefix string) *MissingProperty {
	p.prefix = prefix
	return p
}

// Matches reports whether all configured properties are missing.
func (p *MissingProperty) Matches(ctx context.Context, container component.Container) bool {
	env, err := component.ResolveType[runtime.Environment](ctx, container)
	if err != nil {
		return false
	}

	for _, name := range p.names {
		key := propertyKey(p.prefix, name)

		if _, ok := env.PropertyResolver().Lookup(key); ok {
			return false
		}
	}

	return true
}

func propertyKey(prefix string, name string) string {
	if prefix == "" {
		return name
	}

	return prefix + "." + name
}
