package procyon

import (
	"fmt"
	core "github.com/procyon-projects/procyon-core"
)

var (
	startupLogger = StartupLogger{}
)

type StartupLogger struct {
}

func (logger StartupLogger) LogStarting() {
	core.Log.Info("Starting...")
	core.Log.Info("Application Id : ", core.GetApplicationId())
	core.Log.Info("Running with Procyon, Procyon " + Version)
}

func (logger StartupLogger) LogStarted(watch *core.TaskWatch) {
	lastTime := float32(watch.GetTotalTime()) / 1e9
	formattedText := fmt.Sprintf("Started in %.2f second(s)\n", lastTime)
	core.Log.Info(formattedText)
}
