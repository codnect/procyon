package procyon

import (
	"flag"
	"testing"
)

func init() {
	flag.Bool("nonoptionarg", false, "")
	flag.Bool("fork", false, "")
}

func TestProcyonApplication(t *testing.T) {
	app := NewProcyonApplication()
	app.Run()
	//assert.Equal(t, true, true)
}
