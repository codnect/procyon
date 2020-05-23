package procyon

func init() {
	RegisterAppRunListener(NewEventPublishRunListener())
}
