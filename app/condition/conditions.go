package condition

type OnPropertyCondition struct {
}

func OnProperty(name string) *OnPropertyCondition {
	return nil
}

func (c *OnPropertyCondition) Prefix(prefix string) *OnPropertyCondition {
	return c
}

func (c *OnPropertyCondition) HavingValue(value string) *OnPropertyCondition {
	return c
}

func (c *OnPropertyCondition) MatchIfMissing(matchIfMissing bool) *OnPropertyCondition {
	return c
}

func (c *OnPropertyCondition) Matches(ctx Context) bool {
	return false
}

type OnMissingCondition struct {
}

func OnMissing(name string) *OnMissingCondition {
	return nil
}

func (c *OnMissingCondition) Matches(ctx Context) bool {
	return false
}

type OnTypeCondition struct {
}

func OnType[T any]() *OnTypeCondition {
	return nil
}

func (c *OnTypeCondition) Matches(ctx Context) bool {
	return false
}

type OnMissingTypeCondition struct {
}

func OnMissingType[T any]() *OnMissingTypeCondition {
	return nil
}

func (c *OnMissingTypeCondition) Matches(ctx Context) bool {
	return false
}
