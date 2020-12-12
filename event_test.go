package procyon

import (
	"errors"
	context "github.com/procyon-projects/procyon-context"
	web "github.com/procyon-projects/procyon-web"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testProcyonApplicationEvent(t *testing.T,
	event ProcyonApplicationEvent,
	eventId context.ApplicationEventId,
	parentEventId context.ApplicationEventId,
	source interface{},
	application *ProcyonApplication,
	args ApplicationArguments) {

	assert.Equal(t, eventId, event.GetEventId())
	assert.Equal(t, parentEventId, event.GetParentEventId())
	assert.Equal(t, source, event.GetSource())
	assert.NotEqual(t, int64(0), event.GetTimestamp())
	assert.Equal(t, application, event.GetProcyonApplication())
	assert.Equal(t, args, event.GetArgs())
}

func TestApplicationStartingEvent(t *testing.T) {
	var application = NewProcyonApplication()
	var appArgs = getApplicationArguments(nil)
	event := NewApplicationStarting(application, appArgs)

	testProcyonApplicationEvent(t, event, ApplicationStartingEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, appArgs, event.GetArgs())
}

func TestApplicationEnvironmentPreparedEvent(t *testing.T) {
	var application = NewProcyonApplication()
	var appArgs = getApplicationArguments(nil)
	var environment = web.NewStandardWebEnvironment()
	event := NewApplicationEnvironmentPreparedEvent(application, appArgs, environment)

	testProcyonApplicationEvent(t, event, ApplicationEnvironmentPreparedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, environment, event.GetEnvironment())
}

func TestApplicationContextInitializedEvent(t *testing.T) {
	var application = NewProcyonApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewProcyonServerApplicationContext("app-id", "context-id")
	event := NewApplicationContextInitializedEvent(application, appArgs, ctx)

	testProcyonApplicationEvent(t, event, ApplicationContextInitializedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationPreparedEvent(t *testing.T) {
	var application = NewProcyonApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewProcyonServerApplicationContext("app-id", "context-id")
	event := NewApplicationPreparedEvent(application, appArgs, ctx)

	testProcyonApplicationEvent(t, event, ApplicationPreparedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationStartedEvent(t *testing.T) {
	var application = NewProcyonApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewProcyonServerApplicationContext("app-id", "context-id")
	event := NewApplicationStartedEvent(application, appArgs, ctx)

	testProcyonApplicationEvent(t, event, ApplicationStartedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationReadyEvent(t *testing.T) {
	var application = NewProcyonApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewProcyonServerApplicationContext("app-id", "context-id")
	event := NewApplicationReadyEvent(application, appArgs, ctx)

	testProcyonApplicationEvent(t, event, ApplicationReadyEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
}

func TestApplicationFailedEvent(t *testing.T) {
	var application = NewProcyonApplication()
	var appArgs = getApplicationArguments(nil)
	var ctx = web.NewProcyonServerApplicationContext("app-id", "context-id")
	var err = errors.New("test error")
	event := NewApplicationFailedEvent(application, appArgs, ctx, err)

	testProcyonApplicationEvent(t, event, ApplicationFailedEventId(), ApplicationEventId(), application, application, appArgs)
	assert.Equal(t, ctx, event.GetApplicationContext())
	assert.Equal(t, err, event.GetError())
}
