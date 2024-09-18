package condition

import (
	"codnect.io/procyon/component/filter"
	"codnect.io/procyon/runtime"
)

// OnPropertyCondition struct represents a condition that checks if a specific property has a specific value
// HavingValue and MatchIfMissing allow further customizations.
type OnPropertyCondition struct {
	name           string // The name of the property to check.
	value          any    // The value of the property to check.
	matchIfMissing bool   // Determines if the condition should match if the property is missing.
}

// OnProperty function creates a new OnPropertyCondition.
func OnProperty(name string) *OnPropertyCondition {
	return &OnPropertyCondition{
		name: name,
	}
}

// HavingValue sets the value of the property to check.
func (c *OnPropertyCondition) HavingValue(value any) *OnPropertyCondition {
	c.value = value
	return c
}

// MatchIfMissing sets whether the condition should match if the property is missing.
func (c *OnPropertyCondition) MatchIfMissing(matchIfMissing bool) *OnPropertyCondition {
	c.matchIfMissing = matchIfMissing
	return c
}

// MatchesCondition method checks if the property has the specified value.
// If the property is missing, it returns the value of MatchIfMissing.
// If the property has the specified value, it returns true.
// If the property is a boolean type and the specified value is nil, it returns true if the value is not false.
func (c *OnPropertyCondition) MatchesCondition(ctx Context) bool {
	container := ctx.Container()
	if container == nil {
		return false
	}

	result, err := container.GetObject(ctx, filter.ByTypeOf[runtime.Environment]())
	if err != nil {
		return false
	}

	environment := result.(runtime.Environment)
	property, exists := environment.PropertyResolver().Property(c.name)

	if !exists {
		return c.matchIfMissing
	}

	if c.value == nil {
		switch val := property.(type) {
		case bool:
			return val != false
		case string:
			return val != "false"
		default:
			return true
		}
	}

	return property == c.value
}
