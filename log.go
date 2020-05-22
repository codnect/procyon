package procyon

import (
	core "github.com/Rollcomp/procyon-core"
	"log"
)

var (
	startupLogger = StartupLogger{}
)

type StartupLogger struct {
}

func (logger StartupLogger) LogStarting() {
	log.Println("Starting...")
	log.Println("Running with Procyon, Procyon " + Version)
}

func (logger StartupLogger) LogStarted(watch *core.TaskWatch) {
	lastTime := float32(watch.GetTotalTime()) / 1e9
	log.Printf("Started in %.2f second(s)\n", lastTime)
}
