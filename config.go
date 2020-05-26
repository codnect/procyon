package procyon

import (
	configure "github.com/Rollcomp/procyon-configure"
	core "github.com/Rollcomp/procyon-core"
)

func init() {
	/* Application Run Listeners */
	core.Register(NewEventPublishRunListener)
	/* Application Listeners */
	core.Register(NewBootstrapListener)
	/* Configuration Properties */
	core.Register(
		configure.NewServerConfiguration,
		configure.NewDataSourceConfiguration,
	)
}
