package procyon

import (
	"errors"
	"time"
)

type TaskWatch struct {
	taskName             string
	startTimeNanoSeconds int
	totalTimeNanoSeconds int
}

func NewTaskWatch() *TaskWatch {
	return &TaskWatch{
		taskName: "[empty_task]",
	}
}

func NewTaskWatchWithName(taskName string) *TaskWatch {
	return &TaskWatch{
		taskName: taskName,
	}
}

func (watch *TaskWatch) Start() error {
	if watch.taskName != "" {
		return errors.New("TaskWatch is already running")
	}
	watch.startTimeNanoSeconds = time.Now().Nanosecond()
	return nil
}

func (watch *TaskWatch) Stop() error {
	if watch.taskName == "" {
		return errors.New("TaskWatch is not running")
	}
	watch.totalTimeNanoSeconds = time.Now().Nanosecond() - watch.startTimeNanoSeconds
	watch.taskName = ""
	return nil
}

func (watch *TaskWatch) IsRunning() bool {
	return watch.taskName != ""
}

func (watch *TaskWatch) GetTotalTimeNanoSeconds() int {
	return watch.totalTimeNanoSeconds
}
