package app

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type Arguments struct {
	optArgs     map[string][]string
	nonOptsArgs []string
}

func newArguments() *Arguments {
	return &Arguments{
		optArgs:     make(map[string][]string),
		nonOptsArgs: make([]string, 0),
	}
}

func (a *Arguments) addOptionArgs(name string, value string) {
	if a.optArgs[name] == nil {
		a.optArgs[name] = make([]string, 0)
	}

	a.optArgs[name] = append(a.optArgs[name], value)
}

func (a *Arguments) OptionNames() []string {
	optNames := make([]string, 0)

	for name, _ := range a.optArgs {
		optNames = append(optNames, name)
	}

	return optNames
}

func (a *Arguments) ContainsOption(name string) bool {
	return a.optArgs[name] != nil
}

func (a *Arguments) OptionValues(name string) []string {
	return a.optArgs[name]
}

func (a *Arguments) addNonOptionArgs(value string) {
	a.nonOptsArgs = append(a.nonOptsArgs, value)
}

func (a *Arguments) NonOptionArgs() []string {
	return a.nonOptsArgs
}

func mergeArguments(args ...string) []string {
	mergedArgs := make([]string, 0)
	copy(mergedArgs, os.Args)
	mergedArgs = append(mergedArgs, args...)
	return mergedArgs
}

func parseArguments(args []string) (*Arguments, error) {
	cmdLineArgs := newArguments()
	appArgumentFlagSet := flag.NewFlagSet("ApplicationArguments", flag.ContinueOnError)

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
