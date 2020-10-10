package procyon

import (
	"fmt"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
)

type StartupLogger struct {
	logger context.Logger
}

func NewStartupLogger(logger context.Logger) StartupLogger {
	return StartupLogger{
		logger,
	}
}

func (startupLogger StartupLogger) LogStarting(appId string, contextId string) {
	startupLogger.logger.I(contextId, "Starting...")
	startupLogger.logger.I(contextId, "Application Id : ", appId)
	startupLogger.logger.I(contextId, "Application Context Id : ", contextId)
	startupLogger.logger.I(contextId, "Running with Procyon, Procyon "+Version)
}

func (startupLogger StartupLogger) LogStarted(contextId string, watch *core.TaskWatch) {
	lastTime := float32(watch.GetTotalTime()) / 1e9
	formattedText := fmt.Sprintf("Started in %.2f second(s)\n", lastTime)
	startupLogger.logger.I(contextId, formattedText)
}
