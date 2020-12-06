package procyon

import (
	"fmt"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestApplication_CreateApplicationAndContextId(t *testing.T) {
	procyonApp := NewProcyonApplication()
	appId, contextId := procyonApp.createApplicationAndContextId()
	assert.NotNil(t, appId)
	assert.NotNil(t, contextId)
}

func TestApplication_getAppRunListenerInstances(t *testing.T) {
	procyonApp := NewProcyonApplication()
	runListeners, err := procyonApp.getAppRunListenerInstances(nil)
	assert.NotNil(t, runListeners)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(runListeners.listeners))
}

type loggerMock struct {
	mock.Mock
}

func (l loggerMock) Trace(ctx interface{}, message interface{}) {
}

func (l loggerMock) Debug(ctx interface{}, message interface{}) {
}

func (l loggerMock) Info(ctx interface{}, message interface{}) {
	l.Called(ctx, message)
}

func (l loggerMock) Warning(ctx interface{}, message interface{}) {
}

func (l loggerMock) Error(ctx interface{}, message interface{}) {
}

func (l loggerMock) Fatal(ctx interface{}, message interface{}) {
}

func (l loggerMock) Panic(ctx interface{}, message interface{}) {
}

func (l loggerMock) Tracef(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Debugf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Infof(ctx interface{}, format string, args ...interface{}) {
	l.Called(ctx, format, args)
}

func (l loggerMock) Warningf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Errorf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Fatalf(ctx interface{}, format string, args ...interface{}) {
}

func (l loggerMock) Panicf(ctx interface{}, format string, args ...interface{}) {
}

func Test_logStarting(t *testing.T) {
	var appId context.ApplicationId = "app-id"
	var contextId context.ContextId = "context-id"
	loggerMock := loggerMock{}

	loggerMock.On("Info", contextId, "Starting...")

	var args = make([]interface{}, 0)
	args = append(args, appId)
	loggerMock.On("Infof", contextId, "Application Id : %s", args)

	args = make([]interface{}, 0)
	args = append(args, contextId)
	loggerMock.On("Infof", contextId, "Application Context Id : %s", args)

	loggerMock.On("Info", contextId, "Running with Procyon, Procyon "+Version)

	logStarting(loggerMock, appId, contextId)
	loggerMock.AssertExpectations(t)
}

func Test_logStarted(t *testing.T) {
	var contextId context.ContextId = "context-id"
	loggerMock := loggerMock{}

	taskWatch := core.NewTaskWatch()
	taskWatch.Start()
	time.Sleep(1000)
	taskWatch.Stop()

	lastTime := float32(taskWatch.GetTotalTime()) / 1e9
	formattedText := fmt.Sprintf("Started in %.2f second(s)\n", lastTime)
	loggerMock.On("Info", contextId, formattedText)

	logStarted(loggerMock, contextId, taskWatch)
	loggerMock.AssertExpectations(t)
}
