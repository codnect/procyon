package env

const CmdlinePropertySourceName = "cmdlineArgs"
const NonOptionArgsPropertyName = "nonOptionArgs"

type CommandLinePropertySource interface {
	ContainsOption(name string) bool
	GetOptionValues(name string) []string
	GetNonOptionArgs() []string
}

type BaseCommandLinePropertySource struct {
	*BaseEnumerablePropertySource
}

func NewCommandLinePropertySource(source interface{}) *BaseCommandLinePropertySource {
	return &BaseCommandLinePropertySource{
		NewEnumerablePropertySourceWithSource(CmdlinePropertySourceName, source),
	}
}

func NewCommandLinePropertySourceWithName(name string, source interface{}) *BaseCommandLinePropertySource {
	return &BaseCommandLinePropertySource{
		NewEnumerablePropertySourceWithSource(name, source),
	}
}
