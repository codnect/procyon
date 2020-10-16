package procyon

import (
	"github.com/codnect/goo"
	context "github.com/procyon-projects/procyon-context"
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

type serviceComponentProcessor struct {
}

func newServiceComponentProcessor() serviceComponentProcessor {
	return serviceComponentProcessor{}
}

func (processor serviceComponentProcessor) SupportsComponent(typ goo.Type) bool {
	returnType := typ.(goo.Function).GetFunctionReturnTypes()[0]
	if returnType.(goo.Struct).Implements(goo.GetType((*context.Service)(nil)).(goo.Interface)) {
		return true
	}
	return false
}

func (processor serviceComponentProcessor) ProcessComponent(typ goo.Type) error {
	return nil
}

type repositoryComponentProcessor struct {
}

func newRepositoryComponentProcessor() repositoryComponentProcessor {
	return repositoryComponentProcessor{}
}

func (processor repositoryComponentProcessor) SupportsComponent(typ goo.Type) bool {
	returnType := typ.(goo.Function).GetFunctionReturnTypes()[0]
	if returnType.(goo.Struct).Implements(goo.GetType((*context.Repository)(nil)).(goo.Interface)) {
		return true
	}
	return false
}

func (processor repositoryComponentProcessor) ProcessComponent(typ goo.Type) error {
	return nil
}
