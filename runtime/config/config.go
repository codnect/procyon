package config

import (
	"codnect.io/procyon/runtime/prop"
)

type Data struct {
	source prop.Source
}

// NewData function creates a new Data.
func NewData(source prop.Source) *Data {
	if source == nil {
		panic("nil source")
	}

	return &Data{
		source,
	}
}

func (d *Data) PropertySource() prop.Source {
	return d.source
}
