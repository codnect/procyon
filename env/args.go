package env

import (
	"errors"
	"flag"
	"reflect"
	"strings"
)

type CommandLineArgs struct {
	optionArgs    map[string][]string
	nonOptionArgs []string
}

func NewCommandLineArgs() CommandLineArgs {
	return CommandLineArgs{
		optionArgs:    make(map[string][]string),
		nonOptionArgs: make([]string, 0),
	}
}

func (args *CommandLineArgs) addOptionArgs(name string, value string) {
	if args.optionArgs[name] == nil {
		args.optionArgs[name] = make([]string, 0)
	}
	args.optionArgs[name] = append(args.optionArgs[name], value)
}

func (args CommandLineArgs) getOptionNames() []string {
	argMapKeys := reflect.ValueOf(args.optionArgs).MapKeys()
	optionNames := make([]string, len(argMapKeys))
	for i := 0; i < len(argMapKeys); i++ {
		optionNames[i] = argMapKeys[i].String()
	}
	return optionNames
}

func (args CommandLineArgs) containsOption(name string) bool {
	return args.optionArgs[name] != nil
}

func (args CommandLineArgs) getOptionValues(name string) []string {
	return args.optionArgs[name]
}

func (args *CommandLineArgs) addNonOptionArgs(value string) {
	args.nonOptionArgs = append(args.nonOptionArgs, value)
}

func (args CommandLineArgs) getNonOptionArgs() []string {
	return args.nonOptionArgs
}

type CommandLineArgsParser interface {
	Parse(args []string) (CommandLineArgs, error)
}

type SimpleCommandLineArgsParser struct {
}

func NewCommandLineArgsParser() SimpleCommandLineArgsParser {
	return SimpleCommandLineArgsParser{}
}

func (parser SimpleCommandLineArgsParser) Parse(args []string) (CommandLineArgs, error) {
	cmdLineArgs := NewCommandLineArgs()

	appArgumentFlagSet := flag.NewFlagSet("ProcyonApplicationArguments", flag.ExitOnError)
	err := appArgumentFlagSet.Parse(args)
	if err != nil {
		return cmdLineArgs, err
	}

	for _, arg := range appArgumentFlagSet.Args() {
		if strings.HasPrefix(arg, "--") {
			optionText := arg[2:]
			indexOfEqualSign := strings.Index(optionText, "=")
			optionName := ""
			optionValue := ""
			if indexOfEqualSign > -1 {
				optionName = optionText[0:indexOfEqualSign]
				optionValue = optionText[indexOfEqualSign+1:]
			} else {
				optionName = optionText
			}
			optionName = strings.TrimSpace(optionName)
			optionValue = strings.TrimSpace(optionValue)
			if optionName == "" {
				return cmdLineArgs, errors.New("Invalid argument syntax : " + arg)
			}
			cmdLineArgs.addOptionArgs(optionName, optionValue)
		} else {
			cmdLineArgs.addNonOptionArgs(arg)
		}
	}

	return cmdLineArgs, nil
}
