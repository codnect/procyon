package main

import (
	"codnect.io/procyon"
	"codnect.io/procyon/runtime/cfg"
	"codnect.io/procyon/runtime/property"
)

func main() {
	app := procyon.New()
	if err := app.Run(); err != nil {
		panic(err)
	}

	loader := property.NewYamlSourceLoader()
	loader.Load("", nil)

	var resource cfg.Resource
}
