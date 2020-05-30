package procyon

import (
	configure "github.com/Rollcomp/procyon-configure"
	context "github.com/Rollcomp/procyon-context"
	core "github.com/Rollcomp/procyon-core"
)

func init() {
	/* Configuration Properties */
	core.Register(
		configure.NewServerConfiguration,
		configure.NewDataSourceConfiguration,
	)
	/* Application Run Listeners */
	core.Register(NewEventPublishRunListener)
	/* Application Listeners */
	core.Register(NewBootstrapListener)
	/* Pea Processors */
	core.Register(context.NewConfigurationPropertiesBindingProcessor)
}
