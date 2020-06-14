package procyon

import (
	"fmt"
	core "github.com/procyon-projects/procyon-core"
)

type StartupLogger struct {
	logger core.Logger
}

func NewStartupLogger(logger core.Logger) StartupLogger {
	return StartupLogger{
		logger,
	}
}

func (startupLogger StartupLogger) LogStarting(appId string, contextId string) {
	startupLogger.logger.Info("Starting...")
	startupLogger.logger.Info("Application Id : ", appId)
	startupLogger.logger.Info("Application Context Id : ", contextId)
	startupLogger.logger.Info("Running with Procyon, Procyon " + Version)
}

func (startupLogger StartupLogger) LogStarted(watch *core.TaskWatch) {
	lastTime := float32(watch.GetTotalTime()) / 1e9
	formattedText := fmt.Sprintf("Started in %.2f second(s)\n", lastTime)
	startupLogger.logger.Info(formattedText)
}
