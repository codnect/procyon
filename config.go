package procyon

import (
	configure "github.com/procyon-projects/procyon-configure"
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
}
