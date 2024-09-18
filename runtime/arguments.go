package runtime

import (
	"errors"
	"flag"
	"os"
	"strings"
)

// Arguments struct represents the command line arguments passed to the application.
type Arguments struct {
	optArgs     map[string][]string
	nonOptsArgs []string
}

// newArguments function creates a new Arguments struct.
func newArguments() *Arguments {
	return &Arguments{
		optArgs:     make(map[string][]string),
		nonOptsArgs: make([]string, 0),
	}
}

// OptionNames method returns the names of the option arguments.
func (a *Arguments) OptionNames() []string {
	optNames := make([]string, 0)

	for name, _ := range a.optArgs {
		optNames = append(optNames, name)
	}

	return optNames
}

// ContainsOption method checks whether the option argument with the given name exists.
func (a *Arguments) ContainsOption(name string) bool {
	return a.optArgs[name] != nil
}

// OptionValues method returns the values of the option argument with the given name.
func (a *Arguments) OptionValues(name string) []string {
	return a.optArgs[name]
}

// NonOptionArgs method returns the non-option arguments.
func (a *Arguments) NonOptionArgs() []string {
	return a.nonOptsArgs
}

// addOptionArgs method adds a new option argument to the arguments.
func (a *Arguments) addOptionArgs(name string, value string) {
	if a.optArgs[name] == nil {
		a.optArgs[name] = make([]string, 0)
	}

	a.optArgs[name] = append(a.optArgs[name], value)
}

// addNonOptionArgs method adds a non-option argument to the arguments.
func (a *Arguments) addNonOptionArgs(value string) {
	a.nonOptsArgs = append(a.nonOptsArgs, value)
}

// mergeArguments function merges the command line arguments with the given arguments.
func mergeArguments(args ...string) []string {
	mergedArgs := make([]string, 0)
	copy(mergedArgs, os.Args)
	mergedArgs = append(mergedArgs, args...)
	return mergedArgs
}

// ParseArguments function parses the given and the command line arguments and returns an Arguments.
func ParseArguments(args []string) (*Arguments, error) {
	mergedArgs := mergeArguments(args...)
	cmdLineArgs := newArguments()
	appArgumentFlagSet := flag.NewFlagSet("ApplicationArguments", flag.ContinueOnError)

	err := appArgumentFlagSet.Parse(mergedArgs)

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
				return cmdLineArgs, errors.New("invalid argument syntax : " + arg)
			}

			cmdLineArgs.addOptionArgs(optionName, optionValue)
		} else {
			cmdLineArgs.addNonOptionArgs(arg)
		}

	}

	return cmdLineArgs, nil
}
