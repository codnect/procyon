package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/mock"
	"testing"
)

type applicationEventBroadcasterMock struct {
	mock.Mock
}

func (broadcaster *applicationEventBroadcasterMock) RegisterApplicationListener(listener context.ApplicationListener) {
	broadcaster.Called(listener)
}

func (broadcaster *applicationEventBroadcasterMock) UnregisterApplicationListener(listener context.ApplicationListener) {
	broadcaster.Called(listener)
}

func (broadcaster *applicationEventBroadcasterMock) RemoveAllApplicationListeners() {
	broadcaster.Called()
}

func (broadcaster *applicationEventBroadcasterMock) BroadcastEvent(context context.ApplicationContext, event context.ApplicationEvent) {
	broadcaster.Called(context, event)
}

func TestEventPublishRunListener(t *testing.T) {
	var app = &Application{
		listeners: make([]context.ApplicationListener, 0),
	}
	var appArgs ApplicationArguments
	eventPublishRunListener := NewEventPublishRunListener(app, appArgs)

	applicationEventBroadcasterMock := &applicationEventBroadcasterMock{}
	eventPublishRunListener.broadcaster = applicationEventBroadcasterMock

	var appContext context.ApplicationContext
	applicationEventBroadcasterMock.On("BroadcastEvent", appContext, mock.AnythingOfType("ApplicationStartingEvent")).Return()
	eventPublishRunListener.OnApplicationStarting()
	applicationEventBroadcasterMock.AssertExpectations(t)

	var environment core.ConfigurableEnvironment
	applicationEventBroadcasterMock.On("BroadcastEvent", appContext, mock.AnythingOfType("ApplicationEnvironmentPreparedEvent")).Return()
	eventPublishRunListener.OnApplicationEnvironmentPrepared(environment)
	applicationEventBroadcasterMock.AssertExpectations(t)

	var ctx context.ConfigurableApplicationContext
	applicationEventBroadcasterMock.On("BroadcastEvent", appContext, mock.AnythingOfType("ApplicationContextInitializedEvent")).Return()
	eventPublishRunListener.OnApplicationContextPrepared(ctx)
	applicationEventBroadcasterMock.AssertExpectations(t)

	applicationEventBroadcasterMock.On("BroadcastEvent", appContext, mock.AnythingOfType("ApplicationPreparedEvent")).Return()
	eventPublishRunListener.OnApplicationContextLoaded(ctx)
	applicationEventBroadcasterMock.AssertExpectations(t)

	applicationEventBroadcasterMock.On("BroadcastEvent", appContext, mock.AnythingOfType("ApplicationStartedEvent")).Return()
	eventPublishRunListener.OnApplicationStarted(ctx)
	applicationEventBroadcasterMock.AssertExpectations(t)

	applicationEventBroadcasterMock.On("BroadcastEvent", appContext, mock.AnythingOfType("ApplicationReadyEvent")).Return()
	eventPublishRunListener.OnApplicationRunning(ctx)
	applicationEventBroadcasterMock.AssertExpectations(t)

	var err error
	applicationEventBroadcasterMock.On("BroadcastEvent", appContext, mock.AnythingOfType("ApplicationFailedEvent")).Return()
	eventPublishRunListener.OnApplicationFailed(ctx, err)
	applicationEventBroadcasterMock.AssertExpectations(t)
}
