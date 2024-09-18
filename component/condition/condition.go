package condition

type Condition interface {
	MatchesCondition(ctx Context) bool
}
