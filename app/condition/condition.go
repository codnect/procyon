package condition

type Condition interface {
	Matches(ctx Context) bool
}
