package procyon

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/condition"
	"codnect.io/procyon-core/event"
	"codnect.io/procyon-core/runtime"
)

func init() {
	component.Register(newEnvironmentConfigurer, component.Named("procyonEnvironmentConfigurer"))
	component.Register(runtime.NewEventMulticaster, component.Named("procyonRuntimeEventMulticaster")).
		ConditionalOn(condition.OnMissingType[event.Multicaster]())
}
