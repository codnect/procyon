package procyon

import (
	core "github.com/procyon-projects/procyon-core"
)

var (
	startupLogger = StartupLogger{}
)

type StartupLogger struct {
}

func (logger StartupLogger) LogStarting() {
	core.Logger.Info("Starting...")
	core.Logger.Info("Running with Procyon, Procyon " + Version)
}

func (logger StartupLogger) LogStarted(watch *core.TaskWatch) {
	lastTime := float32(watch.GetTotalTime()) / 1e9
	core.Logger.Infof("Started in %.2f second(s)\n", lastTime)
}
