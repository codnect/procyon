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

	"codnect.io/procyon/runtime"
	"codnect.io/procyon/runtime/config"
)

// configPropertiesProcessor binds configuration properties to components
// that implement the config.Properties interface.
type configPropertiesProcessor struct {
	env runtime.Environment
}

// newConfigPropertiesProcessor creates a new configPropertiesProcessor
// using the given environment as the property source.
func newConfigPropertiesProcessor(env runtime.Environment) *configPropertiesProcessor {
	if env == nil {
		panic("nil environment")
	}

	return &configPropertiesProcessor{
		env: env,
	}
}

// ProcessAfterInit binds configuration properties to the given component
// if it implements the config.Properties interface.
func (c *configPropertiesProcessor) ProcessAfterInit(_ context.Context, instance any) (any, error) {
	if properties, ok := instance.(config.Properties); ok {
		binder := config.NewDefaultPropertyBinder(c.env.PropertySources())

		if err := binder.Bind(properties.Prefix(), properties); err != nil {
			return nil, err
		}
	}

	return instance, nil

}
