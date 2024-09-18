package condition

// OnObjectCondition struct represents a condition that checks if an object exists.
type OnObjectCondition struct {
	name string // The name of the object to check.
}

// OnObject function creates a new OnObjectCondition.
func OnObject(name string) *OnObjectCondition {
	return &OnObjectCondition{
		name: name,
	}
}

// MatchesCondition method checks if the object exists.
// It returns true if the object exists, false otherwise.
func (c *OnObjectCondition) MatchesCondition(ctx Context) bool {
	container := ctx.Container()
	if container == nil {
		return false
	}

	return container.Definitions().Contains(c.name) || container.Singletons().Contains(c.name)
}

// OnMissingObjectCondition struct represents a condition that checks if an object does not exist.
type OnMissingObjectCondition struct {
	name string // The name of the object to check.
}

// OnMissingObject function creates a new OnMissingObjectCondition.
func OnMissingObject(name string) *OnMissingObjectCondition {
	return &OnMissingObjectCondition{
		name: name,
	}
}

// MatchesCondition method checks if the object does not exist.
// It returns true if the object does not exist, false otherwise.
func (c *OnMissingObjectCondition) MatchesCondition(ctx Context) bool {
	container := ctx.Container()
	if container == nil {
		return false
	}

	return !container.Definitions().Contains(c.name) && !container.Singletons().Contains(c.name)
}
