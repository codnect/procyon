package procyon

import (
	configure "github.com/procyon-projects/procyon-configure"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
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
