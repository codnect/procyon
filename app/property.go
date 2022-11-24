package app

import "strings"

const (
	NonOptionArgs = "nonOptionArgs"
)

type ArgumentsPropertySource struct {
	args *Arguments
}

func newArgumentPropertySources(args *Arguments) *ArgumentsPropertySource {
	return &ArgumentsPropertySource{
		args: args,
	}
}

func (s *ArgumentsPropertySource) Name() string {
	return "commandLineArgs"
}

func (s *ArgumentsPropertySource) Source() any {
	return s.args
}

func (s *ArgumentsPropertySource) ContainsProperty(name string) bool {
	if NonOptionArgs == name {
		return len(s.args.NonOptionArgs()) != 0
	}

	return s.args.ContainsOption(name)
}

func (s *ArgumentsPropertySource) Property(name string) (any, bool) {
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

func (s *ArgumentsPropertySource) PropertyOrDefault(name string, defaultValue any) any {
	val, ok := s.Property(name)
	if !ok {
		return defaultValue
	}

	return val
}

func (s *ArgumentsPropertySource) PropertyNames() []string {
	return s.args.OptionNames()
}

func (s *ArgumentsPropertySource) OptionValues(name string) []string {
	return s.args.OptionValues(name)
}

func (s *ArgumentsPropertySource) NonOptionArgs() []string {
	return s.args.NonOptionArgs()
}
