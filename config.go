package procyon

import (
	"github.com/Rollcomp/procyon-configure"
)

func init() {
	/* Application Run Listeners */
	Register(NewEventPublishRunListener)
	/* Configuration Properties */
	Register(
		configure.NewServerConfiguration,
		configure.NewDataSourceConfiguration,
	)
}
