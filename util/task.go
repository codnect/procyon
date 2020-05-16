package util

import (
	"errors"
	"time"
)

type TaskWatch struct {
	taskName  string
	startTime int64
	totalTime int64
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
	if watch.taskName != "" && watch.startTime != 0 {
		return errors.New("TaskWatch is already running")
	}
	watch.startTime = time.Now().Unix()
	return nil
}

func (watch *TaskWatch) Stop() error {
	if watch.taskName == "" {
		return errors.New("TaskWatch is not running")
	}
	watch.totalTime = time.Now().Unix() - watch.startTime
	watch.taskName = ""
	return nil
}

func (watch *TaskWatch) IsRunning() bool {
	return watch.taskName != ""
}

func (watch *TaskWatch) GetTotalTime() int64 {
	return watch.totalTime
}
