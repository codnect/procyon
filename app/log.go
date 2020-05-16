package app

import (
	"log"
	"procyon/util"
)

type AppStartupLogger struct {
}

func NewAppStartupLogger() AppStartupLogger {
	return AppStartupLogger{}
}

func (logger AppStartupLogger) LogStarting() {
	log.Println("Starting...")
	log.Println("Running with Procyon, Procyon " + Version)
}

func (logger AppStartupLogger) LogStarted(watch *util.TaskWatch) {
	lastTime := float32(watch.GetTotalTimeNanoSeconds()) / 1e9
	log.Printf("Started in %.2f second(s)\n", lastTime)
}
