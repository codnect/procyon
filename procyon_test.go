package procyon

import (
	"testing"
)

func TestProcyonApplication(t *testing.T) {
	//x := core.GetComponentTypes(core.GetType((*EventPublishRunListener)(nil)))
	//log.Print(x)
	NewProcyonApplication().Run()
}
