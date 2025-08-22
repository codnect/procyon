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

import "codnect.io/procyon/runtime/property"

// Data holds a property source for configuration data.
type Data struct {
	source property.Source
}

// NewData creates a new Data with the given property source.
func NewData(source property.Source) *Data {
	return &Data{
		source: source,
	}
}

// PropertySource returns the property source.
func (d *Data) PropertySource() property.Source {
	return d.source
}
