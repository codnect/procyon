package config

import (
	"codnect.io/procyon/runtime/property"
)

type Config struct {
	source property.Source
}

// New function creates a new Config.
func New(source property.Source) *Config {
	if source == nil {
		panic("nil property source")
	}

	return &Config{
		source,
	}
}

func (d *Config) PropertySource() property.Source {
	return d.source
}
