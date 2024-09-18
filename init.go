package procyon

import (
	"codnect.io/procyon/component"
	"codnect.io/procyon/component/condition"
	"codnect.io/procyon/runtime"
	"codnect.io/procyon/runtime/config"
	"codnect.io/procyon/runtime/event"
	"codnect.io/procyon/runtime/property"
)

func init() {
	// core
	component.Register(newConfigContextConfigurer, component.WithName("procyonConfigContextConfigurer"))
	// runtime/event
	component.Register(event.NewSimpleMulticaster, component.WithName("procyonEventMulticaster"),
		component.WithCondition(condition.OnMissingType[event.Multicaster]()),
	)
	// runtime/config
	component.Register(config.NewDefaultResourceResolver, component.WithName("procyonDefaultConfigResourceResolver"))
	component.Register(config.NewFileLoader, component.WithName("procyonConfigFileLoader"))
	component.Register(config.NewImporter, component.WithName("procyonConfigImporter"))
	// runtime/property
	component.Register(property.NewYamlSourceLoader, component.WithName("procyonYamlPropertySourceLoader"))
	// runtime
	component.Register(runtime.NewServerProperties, component.WithPrototypeScope())
	component.Register(runtime.NewLifecycleProperties, component.WithSingletonScope())
}
