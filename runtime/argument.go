// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"errors"
	"flag"
	"os"
	"strings"
)

const (
	// NonOptionArgs represents the non-option arguments.
	NonOptionArgs = "nonOptionArgs"
)

// Args struct represents the command line arguments passed to the application.
type Args struct {
	optArgs     map[string][]string
	nonOptsArgs []string
}

// newArgs function creates a new Args struct.
func newArgs() *Args {
	return &Args{
		optArgs:     make(map[string][]string),
		nonOptsArgs: make([]string, 0),
	}
}

// OptionNames method returns the names of the option arguments.
func (a *Args) OptionNames() []string {
	optNames := make([]string, 0)

	for name := range a.optArgs {
		optNames = append(optNames, name)
	}

	return optNames
}

// ContainsOption method checks whether the option argument with the given name exists.
func (a *Args) ContainsOption(name string) bool {
	return a.optArgs[name] != nil
}

// OptionValues method returns the values of the option argument with the given name.
func (a *Args) OptionValues(name string) []string {
	return a.optArgs[name]
}

// NonOptionArgs method returns the non-option arguments.
func (a *Args) NonOptionArgs() []string {
	return a.nonOptsArgs
}

// addOptionArgs method adds a new option argument to the arguments.
func (a *Args) addOptionArgs(name string, value string) {
	if a.optArgs[name] == nil {
		a.optArgs[name] = make([]string, 0)
	}

	a.optArgs[name] = append(a.optArgs[name], value)
}

// addNonOptionArgs method adds a non-option argument to the arguments.
func (a *Args) addNonOptionArgs(value string) {
	a.nonOptsArgs = append(a.nonOptsArgs, value)
}

// mergeArguments function merges the command line arguments with the given arguments.
func mergeArguments(args ...string) []string {
	mergedArgs := make([]string, 0)
	copy(mergedArgs, os.Args)
	mergedArgs = append(mergedArgs, args...)
	return mergedArgs
}

// ParseArgs function parses the given and the command line arguments and returns an Args.
func ParseArgs(args []string) (*Args, error) {
	mergedArgs := mergeArguments(args...)
	cmdLineArgs := newArgs()
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

// ArgsPropertySource struct represents a source of arguments.
type ArgsPropertySource struct {
	args *Args
}

// NewArgsPropertySource function creates a new ArgsPropertySource with the given arguments.
func NewArgsPropertySource(args *Args) *ArgsPropertySource {
	return &ArgsPropertySource{
		args: args,
	}
}

// Name method returns the name of the source.
func (s *ArgsPropertySource) Name() string {
	return "commandLineArgs"
}

// Underlying returns the underlying source object.
func (s *ArgsPropertySource) Underlying() any {
	return s.args
}

// ContainsProperty method checks whether the argument with the given name exists.
func (s *ArgsPropertySource) ContainsProperty(name string) bool {
	if NonOptionArgs == name {
		return len(s.args.NonOptionArgs()) != 0
	}

	return s.args.ContainsOption(name)
}

// Property method returns the value of the argument with the given name.
func (s *ArgsPropertySource) Property(name string) (any, bool) {
	if NonOptionArgs == name {
		nonOptValues := s.args.NonOptionArgs()

		if nonOptValues != nil {
			return strings.Join(nonOptValues, ","), true
		}

		return nil, false
	}

	optValues := s.args.OptionValues(name)

	if optValues != nil {
		return strings.Join(optValues, ","), true
	}

	return nil, false
}

// PropertyOrDefault returns the value of the given argument name from the source.
// If the argument does not exist, it returns the default value.
func (s *ArgsPropertySource) PropertyOrDefault(name string, defaultValue any) any {
	val, ok := s.Property(name)
	if !ok {
		return defaultValue
	}

	return val
}

// PropertyNames method returns the names of the arguments.
func (s *ArgsPropertySource) PropertyNames() []string {
	return s.args.OptionNames()
}

// OptionValues method returns the values of the option argument with the given name.
func (s *ArgsPropertySource) OptionValues(name string) []string {
	return s.args.OptionValues(name)
}

// NonOptionArgs method returns the non-option arguments.
func (s *ArgsPropertySource) NonOptionArgs() []string {
	return s.args.NonOptionArgs()
}
