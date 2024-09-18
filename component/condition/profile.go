package condition

import (
	"codnect.io/procyon/component/filter"
	"codnect.io/procyon/runtime"
)

// OnProfileCondition struct represents a condition that checks if a specific profile is active.
type OnProfileCondition struct {
	profiles []string // The profiles to check.
}

// OnProfile function creates a new OnProfileCondition.
func OnProfile(profiles ...string) *OnProfileCondition {
	return &OnProfileCondition{
		profiles: profiles,
	}
}

// MatchesCondition method checks if the profiles are active.
// It retrieves the runtime environment from the container and checks if each profile in the list is active.
// If all profiles are active, it returns true. If any profile is not active, it returns false.
func (c *OnProfileCondition) MatchesCondition(ctx Context) bool {
	container := ctx.Container()
	if container == nil {
		return false
	}

	result, err := container.GetObject(ctx, filter.ByTypeOf[runtime.Environment]())
	if err != nil {
		return false
	}

	environment := result.(runtime.Environment)

	for _, profile := range c.profiles {
		if !environment.IsProfileActive(profile) {
			return false
		}
	}

	return true
}
