package procyon

import (
	"github.com/procyon-projects/procyon-core"
)

type ApplicationArguments interface {
	ContainsOption(name string) bool
	GetOptionNames() []string
	GetOptionValues(name string) []string
	GetSourceArgs() []string
	GetNonOptionArgs() []string
}

type DefaultApplicationArguments struct {
	source core.SimpleCommandLinePropertySource
	args   []string
}

func GetApplicationArguments(args []string) ApplicationArguments {
	return &DefaultApplicationArguments{
		args:   args,
		source: core.NewSimpleCommandLinePropertySource(args),
	}
}

func (arg DefaultApplicationArguments) ContainsOption(name string) bool {
	return arg.source.ContainsOption(name)
}

func (arg DefaultApplicationArguments) GetOptionNames() []string {
	return arg.source.GetPropertyNames()
}

func (arg DefaultApplicationArguments) GetOptionValues(name string) []string {
	return arg.source.GetOptionValues(name)
}

func (arg DefaultApplicationArguments) GetSourceArgs() []string {
	return arg.args
}

func (arg DefaultApplicationArguments) GetNonOptionArgs() []string {
	return arg.source.GetNonOptionArgs()
}
