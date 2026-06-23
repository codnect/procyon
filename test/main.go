package main

import (
	"os"

	"codnect.io/procyon"
	"codnect.io/procyon/component"
	"codnect.io/procyon/runtime"
)

func main() {
	component.Register(newRunner)

	app := procyon.New()
	if err := app.Run(os.Args...); err != nil {
		os.Exit(1)
	}
}

type Runner struct {
}

func newRunner() Runner {
	return Runner{}
}

func (r Runner) Run(ctx runtime.Context, args *runtime.Args) error {
	//TODO implement me
	panic("implement me")
}
