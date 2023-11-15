package web

type ServerStartStopLifecycle struct {
	server  Server
	running bool
}

func (l *ServerStartStopLifecycle) Start() error {
	err := l.server.Start()

	if err != nil {
		return err
	}

	l.running = true
	return nil
}

func (l *ServerStartStopLifecycle) Stop() error {
	err := l.server.Stop()

	if err != nil {
		return err
	}

	l.running = false
	return nil
}

func (l *ServerStartStopLifecycle) IsRunning() bool {
	return l.running
}
