package app

import (
	"log"
	"procyon/util"
)

type StartupLogger struct {
}

func NewStartupLogger() StartupLogger {
	return StartupLogger{}
}

func (logger StartupLogger) LogStarting() {
	log.Println("Starting...")
	log.Println("Running with Procyon, Procyon " + Version)
}

func (logger StartupLogger) LogStarted(watch *util.TaskWatch) {
	lastTime := float32(watch.GetTotalTime()) / 1e9
	log.Printf("Started in %.2f second(s)\n", lastTime)
}
