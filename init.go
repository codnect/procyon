package procyon

/*
import (
	"codnect.io/procyon"
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/condition"
	"codnect.io/procyon/app/availability"
)

func init() {
	// app
	component.Register(procyon.newStartupListener, component.Name("procyonStartupListener"))
	component.Register(procyon.newEnvironmentCustomizer, component.Name("procyonEnvironmentCustomizer"))
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
*/
