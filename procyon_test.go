package procyon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplication_CreateApplicationAndContextId(t *testing.T) {
	procyonApp := NewProcyonApplication()
	appId, contextId := procyonApp.createApplicationAndContextId()
	assert.NotNil(t, appId)
	assert.NotNil(t, contextId)
}
