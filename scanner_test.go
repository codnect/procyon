package procyon

import (
	"github.com/codnect/goo"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComponentScanner_scan(t *testing.T) {
	var contextIdArray [36]byte
	core.GenerateUUID(contextIdArray[:])
	contextId := context.ContextId(contextIdArray[:])

	loggerMock := loggerMock{}

	componentScanner := newComponentScanner()
	count, err := componentScanner.scan(contextId, loggerMock)
	assert.Equal(t, 10, count)
	assert.Nil(t, err)
}

func TestComponentScanner_checkComponent(t *testing.T) {
	componentScanner := newComponentScanner()
	instances, err := componentScanner.getProcessorInstances()
	assert.Nil(t, err)
	err = componentScanner.checkComponent(goo.GetType(newControllerStruct), instances)
	assert.Nil(t, err)
}

func TestComponentScanner_getProcessorInstances(t *testing.T) {
	componentScanner := newComponentScanner()
	instances, err := componentScanner.getProcessorInstances()
	assert.Equal(t, 1, len(instances))
	assert.Nil(t, err)
}
