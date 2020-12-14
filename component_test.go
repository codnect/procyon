package procyon

import (
	"github.com/procyon-projects/goo"
	web "github.com/procyon-projects/procyon-web"
	"github.com/stretchr/testify/assert"
	"testing"
)

type NonControllerStruct struct {
}

func newNonControllerStruct() NonControllerStruct {
	return NonControllerStruct{}
}

type ControllerStruct struct {
}

func newControllerStruct() ControllerStruct {
	return ControllerStruct{}
}

func (controller ControllerStruct) RegisterHandlers(registry web.HandlerRegistry) {

}

func TestControllerComponentProcessor_SupportsComponent(t *testing.T) {
	processor := newControllerComponentProcessor()
	assert.True(t, processor.SupportsComponent(goo.GetType(newControllerStruct)))
	assert.False(t, processor.SupportsComponent(goo.GetType(newNonControllerStruct)))
}

func TestControllerComponentProcessor_ProcessComponent(t *testing.T) {
	processor := newControllerComponentProcessor()
	assert.Nil(t, processor.ProcessComponent(goo.GetType(newControllerStruct)))
	assert.NotNil(t, processor.ProcessComponent(goo.GetType(newNonControllerStruct)))
}
