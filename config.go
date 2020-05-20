package procyon

import "procyon/event"

func init() {
	RegisterAppRunListener(event.NewPublishRunListener())
}
