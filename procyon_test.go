package procyon

import (
	"testing"
)

func TestProcyonApplication(t *testing.T) {
	app := NewProcyonApplication()
	app.SetApplicationRunListeners()
	app.Run()
	//assert.Equal(t, true, true)
}
