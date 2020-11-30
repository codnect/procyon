package procyon

import (
	core "github.com/procyon-projects/procyon-core"
)

func init() {
	/* Default Component Processors */
	core.Register(
		newControllerComponentProcessor,
	)
	/* Application Run Listeners */
	core.Register(NewEventPublishRunListener)
	/* Application Listeners */
	core.Register(NewBootstrapListener)
}
