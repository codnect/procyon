package app

import (
	"github.com/procyon-projects/procyon/app/availability"
	"github.com/procyon-projects/procyon/app/component"
	"github.com/procyon-projects/procyon/app/condition"
)

func init() {
	// app
	component.Register(newStartupListener, component.Name("procyonStartupListener"))

	// availability
	component.Register(availability.NewStateHolder, component.Name("availabilityStateHolder"))
	component.Register(availability.NewLivenessStateHealthChecker).
		ConditionalOn(condition.OnMissing("livenessStateHealthChecker")).
		ConditionalOn(condition.OnProperty("enabled").
			Prefix("procyon.health.check.livenessstate").
			HavingValue("true"),
		)
	component.Register(availability.NewReadinessStateHealthChecker).
		ConditionalOn(condition.OnMissing("readinessStateHealthChecker")).
		ConditionalOn(condition.OnProperty("enabled").
			Prefix("procyon.health.check.readiness").
			HavingValue("true"),
		)
}
