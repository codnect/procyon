package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/mock"
	"testing"
)

type applicationRunListenerMock struct {
	mock.Mock
}

func (listener *applicationRunListenerMock) OnApplicationStarting() {
	listener.Called()
}

func (listener *applicationRunListenerMock) OnApplicationEnvironmentPrepared(environment core.ConfigurableEnvironment) {
	listener.Called(environment)
}

func (listener *applicationRunListenerMock) OnApplicationContextPrepared(ctx context.ConfigurableApplicationContext) {
	listener.Called(ctx)
}

func (listener *applicationRunListenerMock) OnApplicationContextLoaded(ctx context.ConfigurableApplicationContext) {
	listener.Called(ctx)
}

func (listener *applicationRunListenerMock) OnApplicationStarted(ctx context.ConfigurableApplicationContext) {
	listener.Called(ctx)
}

func (listener *applicationRunListenerMock) OnApplicationRunning(ctx context.ConfigurableApplicationContext) {
	listener.Called(ctx)
}

func (listener *applicationRunListenerMock) OnApplicationFailed(ctx context.ConfigurableApplicationContext, err error) {
	listener.Called(ctx, err)
}

func TestApplicationRunListeners(t *testing.T) {
	applicationRunListenerMock := new(applicationRunListenerMock)
	listeners := make([]ApplicationRunListener, 0)
	listeners = append(listeners, applicationRunListenerMock)
	applicationRunListeners := NewApplicationRunListeners(listeners)

	applicationRunListenerMock.On("OnApplicationStarting").Return()
	applicationRunListeners.OnApplicationStarting()
	applicationRunListenerMock.AssertExpectations(t)

	var environment core.ConfigurableEnvironment
	applicationRunListenerMock.On("OnApplicationEnvironmentPrepared", environment).Return()
	applicationRunListeners.OnApplicationEnvironmentPrepared(environment)
	applicationRunListenerMock.AssertExpectations(t)

	var ctx context.ConfigurableApplicationContext
	applicationRunListenerMock.On("OnApplicationContextPrepared", ctx).Return()
	applicationRunListeners.OnApplicationContextPrepared(ctx)
	applicationRunListenerMock.AssertExpectations(t)

	applicationRunListenerMock.On("OnApplicationContextLoaded", ctx).Return()
	applicationRunListeners.OnApplicationContextLoaded(ctx)
	applicationRunListenerMock.AssertExpectations(t)

	applicationRunListenerMock.On("OnApplicationStarted", ctx).Return()
	applicationRunListeners.OnApplicationStarted(ctx)
	applicationRunListenerMock.AssertExpectations(t)

	applicationRunListenerMock.On("OnApplicationRunning", ctx).Return()
	applicationRunListeners.OnApplicationRunning(ctx)
	applicationRunListenerMock.AssertExpectations(t)

	var err error
	applicationRunListenerMock.On("OnApplicationFailed", ctx, err).Return()
	applicationRunListeners.OnApplicationFailed(ctx, err)
	applicationRunListenerMock.AssertExpectations(t)
}
