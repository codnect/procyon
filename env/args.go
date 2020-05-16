package env

import "flag"

type CommandLineArgs struct {
	optionArgs    map[string][]string
	nonOptionArgs []string
}

func NewCommandLineArgs() *CommandLineArgs {
	return &CommandLineArgs{
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

func (args *CommandLineArgs) getOptionNames() []string {
	return nil
}

func (args *CommandLineArgs) containsOption(name string) bool {
	return args.optionArgs[name] != nil
}

func (args *CommandLineArgs) getOptionValues(name string) []string {
	return args.optionArgs[name]
}

func (args *CommandLineArgs) addNonOptionArgs(value string) {
	args.nonOptionArgs = append(args.nonOptionArgs, value)
}

func (args *CommandLineArgs) getNonOptionArgs() []string {
	return args.nonOptionArgs
}

type CommandLineArgsParser interface {
	Parse(args []string) *CommandLineArgs
}

type SimpleCommandLineArgsParser struct {
}

func NewCommandLineArgsParser() CommandLineArgsParser {
	return &SimpleCommandLineArgsParser{}
}

func (parser *SimpleCommandLineArgsParser) Parse(args []string) *CommandLineArgs {
	appArgumentFlagSet := flag.NewFlagSet("ProcyonApplicationArguments", flag.ContinueOnError)
	_ = appArgumentFlagSet.Parse(args)
	//TODO will be implemented
	return nil
}
