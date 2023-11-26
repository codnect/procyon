package app

import (
	"codnect.io/procyon/app/availability"
	"codnect.io/procyon/app/condition"
	"codnect.io/procyon/component"
)

func init() {
	// app
	component.Register(newStartupListener, component.Name("procyonStartupListener"))
	component.Register(newEnvironmentCustomizer, component.Name("procyonEnvironmentCustomizer"))
	component.Register(NewDefaultLifecycleProcessor, component.Name("lifecycleProcessor")).
		ConditionalOn(condition.OnMissing("lifecycleProcessor"))

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
