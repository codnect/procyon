package procyon

import (
	"github.com/codnect/goo"
	web "github.com/procyon-projects/procyon-web"
)

type controllerComponentProcessor struct {
}

func newControllerComponentProcessor() controllerComponentProcessor {
	return controllerComponentProcessor{}
}

func (processor controllerComponentProcessor) SupportsComponent(typ goo.Type) bool {
	returnType := typ.(goo.Function).GetFunctionReturnTypes()[0]
	if returnType.(goo.Struct).Implements(goo.GetType((*web.Controller)(nil)).(goo.Interface)) {
		return true
	}
	return false
}

func (processor controllerComponentProcessor) ProcessComponent(typ goo.Type) error {
	return nil
}
