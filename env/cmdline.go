package env

import (
	"strings"
)

const CmdlinePropertySourceName = "cmdlineArgs"
const NonOptionArgsPropertyName = "nonOptionArgs"

type CommandLinePropertySource interface {
	ContainsOption(name string) bool
	GetOptionValues(name string) []string
	GetNonOptionArgs() []string
}

type BaseCommandLinePropertySource struct {
	CommandLinePropertySource
	*BaseEnumerablePropertySource
}

func NewCommandLinePropertySource(source interface{}) *BaseCommandLinePropertySource {
	return &BaseCommandLinePropertySource{
		BaseEnumerablePropertySource: NewEnumerablePropertySourceWithSource(CmdlinePropertySourceName, source),
	}
}

func NewCommandLinePropertySourceWithName(name string, source interface{}) *BaseCommandLinePropertySource {
	return &BaseCommandLinePropertySource{
		BaseEnumerablePropertySource: NewEnumerablePropertySourceWithSource(name, source),
	}
}

func (source *BaseCommandLinePropertySource) ContainsProperty(name string) bool {
	if NonOptionArgsPropertyName == name {
		return len(source.GetNonOptionArgs()) != 0
	}
	return source.ContainsOption(name)
}

func (source *BaseCommandLinePropertySource) GetProperty(name string) interface{} {
	if NonOptionArgsPropertyName == name {
		nonOptValues := source.GetNonOptionArgs()
		if nonOptValues != nil {
			return strings.Join(nonOptValues, ",")
		}
		return nil
	}
	optValues := source.GetOptionValues(name)
	if optValues != nil {
		return strings.Join(optValues, ",")
	}
	return nil
}

type SimpleCommandLinePropertySource struct {
	*BaseCommandLinePropertySource
}

func NewSimpleCommandLinePropertySource(args []string) SimpleCommandLinePropertySource {
	parser := NewCommandLineArgsParser()
	return SimpleCommandLinePropertySource{
		BaseCommandLinePropertySource: NewCommandLinePropertySource(parser.Parse(args)),
	}
}

func SimpleCommandLinePropertySourceWithName(name string, args []string) SimpleCommandLinePropertySource {
	parser := NewCommandLineArgsParser()
	return SimpleCommandLinePropertySource{
		BaseCommandLinePropertySource: NewCommandLinePropertySourceWithName(name, parser.Parse(args)),
	}
}

func (source SimpleCommandLinePropertySource) ContainsOption(name string) bool {
	cmdLineArgs := source.Source.(CommandLineArgs)
	return cmdLineArgs.containsOption(name)
}

func (source SimpleCommandLinePropertySource) GetOptionValues(name string) []string {
	cmdLineArgs := source.Source.(CommandLineArgs)
	return cmdLineArgs.getOptionValues(name)
}

func (source SimpleCommandLinePropertySource) GetNonOptionArgs() []string {
	cmdLineArgs := source.Source.(CommandLineArgs)
	return cmdLineArgs.getNonOptionArgs()
}

func (source SimpleCommandLinePropertySource) GetPropertyNames() []string {
	cmdLineArgs := source.Source.(CommandLineArgs)
	return cmdLineArgs.getOptionNames()
}
