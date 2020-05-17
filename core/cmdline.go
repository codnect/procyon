package core

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

type AbstractCommandLinePropertySource struct {
	CommandLinePropertySource
	AbstractEnumerablePropertySource
}

func NewCommandLinePropertySource(source interface{}) AbstractCommandLinePropertySource {
	cmdLinePropertySource := AbstractCommandLinePropertySource{
		AbstractEnumerablePropertySource: NewAbstractEnumerablePropertySourceWithSource(CmdlinePropertySourceName, source),
	}
	cmdLinePropertySource.EnumerablePropertySource = cmdLinePropertySource
	return cmdLinePropertySource
}

func NewCommandLinePropertySourceWithName(name string, source interface{}) AbstractCommandLinePropertySource {
	cmdLinePropertySource := AbstractCommandLinePropertySource{
		AbstractEnumerablePropertySource: NewAbstractEnumerablePropertySourceWithSource(name, source),
	}
	cmdLinePropertySource.EnumerablePropertySource = cmdLinePropertySource
	return cmdLinePropertySource
}

func (source AbstractCommandLinePropertySource) ContainsProperty(name string) bool {
	if NonOptionArgsPropertyName == name {
		return len(source.CommandLinePropertySource.GetNonOptionArgs()) != 0
	}
	return source.CommandLinePropertySource.ContainsOption(name)
}

func (source AbstractCommandLinePropertySource) GetProperty(name string) interface{} {
	if NonOptionArgsPropertyName == name {
		nonOptValues := source.CommandLinePropertySource.GetNonOptionArgs()
		if nonOptValues != nil {
			return strings.Join(nonOptValues, ",")
		}
		return nil
	}
	optValues := source.CommandLinePropertySource.GetOptionValues(name)
	if optValues != nil {
		return strings.Join(optValues, ",")
	}
	return nil
}

type SimpleCommandLinePropertySource struct {
	AbstractCommandLinePropertySource
}

func NewSimpleCommandLinePropertySource(args []string) SimpleCommandLinePropertySource {
	cmdLineArgs, err := NewCommandLineArgsParser().Parse(args)
	if err != nil {
		panic(err)
	}
	cmdlinePropertySource := SimpleCommandLinePropertySource{
		AbstractCommandLinePropertySource: NewCommandLinePropertySource(cmdLineArgs),
	}
	cmdlinePropertySource.AbstractCommandLinePropertySource.CommandLinePropertySource = cmdlinePropertySource
	return cmdlinePropertySource
}

func SimpleCommandLinePropertySourceWithName(name string, args []string) SimpleCommandLinePropertySource {
	cmdLineArgs, err := NewCommandLineArgsParser().Parse(args)
	if err != nil {
		panic(err)
	}
	cmdlinePropertySource := SimpleCommandLinePropertySource{
		AbstractCommandLinePropertySource: NewCommandLinePropertySourceWithName(name, cmdLineArgs),
	}
	cmdlinePropertySource.AbstractCommandLinePropertySource.CommandLinePropertySource = cmdlinePropertySource
	return cmdlinePropertySource
}

func (source SimpleCommandLinePropertySource) ContainsOption(name string) bool {
	cmdLineArgs := source.GetSource().(CommandLineArgs)
	return cmdLineArgs.containsOption(name)
}

func (source SimpleCommandLinePropertySource) GetOptionValues(name string) []string {
	cmdLineArgs := source.GetSource().(CommandLineArgs)
	return cmdLineArgs.getOptionValues(name)
}

func (source SimpleCommandLinePropertySource) GetNonOptionArgs() []string {
	cmdLineArgs := source.GetSource().(CommandLineArgs)
	return cmdLineArgs.getNonOptionArgs()
}

func (source SimpleCommandLinePropertySource) GetPropertyNames() []string {
	cmdLineArgs := source.GetSource().(CommandLineArgs)
	return cmdLineArgs.getOptionNames()
}
