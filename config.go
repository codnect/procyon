package procyon

import (
	core "github.com/procyon-projects/procyon-core"
)

func init() {
	/* Default Component Processors */
	core.Register(
		newRepositoryComponentProcessor,
		newServiceComponentProcessor,
		newControllerComponentProcessor,
	)
	/* Application Run Listeners */
	core.Register(NewEventPublishRunListener)
	/* Application Listeners */
	core.Register(NewBootstrapListener)
}
