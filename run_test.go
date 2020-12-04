package procyon

import (
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestApplicationRunListener struct {
	counter int
}

func newTestApplicationRunListener() *TestApplicationRunListener {
	return &TestApplicationRunListener{}
}

func (listener *TestApplicationRunListener) OnApplicationStarting() {
	listener.counter++
}

func (listener *TestApplicationRunListener) OnApplicationEnvironmentPrepared(environment core.ConfigurableEnvironment) {
	listener.counter++
}

func (listener *TestApplicationRunListener) OnApplicationContextPrepared(ctx context.ConfigurableApplicationContext) {
	listener.counter++
}

func (listener *TestApplicationRunListener) OnApplicationContextLoaded(ctx context.ConfigurableApplicationContext) {
	listener.counter++
}

func (listener *TestApplicationRunListener) OnApplicationStarted(ctx context.ConfigurableApplicationContext) {
	listener.counter++
}

func (listener *TestApplicationRunListener) OnApplicationRunning(ctx context.ConfigurableApplicationContext) {
	listener.counter++
}

func (listener *TestApplicationRunListener) OnApplicationFailed(ctx context.ConfigurableApplicationContext, err error) {
	listener.counter++
}

func TestApplicationRunListeners(t *testing.T) {
	runListeners := make([]ApplicationRunListener, 0)
	testRunListener := newTestApplicationRunListener()
	runListeners = append(runListeners, testRunListener)
	applicationRunListeners := NewApplicationRunListeners(runListeners)

	counter := 1
	applicationRunListeners.OnApplicationStarting()
	assert.Equal(t, counter, testRunListener.counter)
	counter++

	applicationRunListeners.OnApplicationEnvironmentPrepared(nil)
	assert.Equal(t, counter, testRunListener.counter)
	counter++

	applicationRunListeners.OnApplicationContextPrepared(nil)
	assert.Equal(t, counter, testRunListener.counter)
	counter++

	applicationRunListeners.OnApplicationContextLoaded(nil)
	assert.Equal(t, counter, testRunListener.counter)
	counter++

	applicationRunListeners.OnApplicationStarted(nil)
	assert.Equal(t, counter, testRunListener.counter)
	counter++

	applicationRunListeners.OnApplicationRunning(nil)
	assert.Equal(t, counter, testRunListener.counter)
	counter++

	applicationRunListeners.OnApplicationFailed(nil, nil)
	assert.Equal(t, counter, testRunListener.counter)
}
