package procyon

import (
	"errors"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	web "github.com/procyon-projects/procyon-web"
)

type controllerComponentProcessor struct {
}

func newControllerComponentProcessor() controllerComponentProcessor {
	return controllerComponentProcessor{}
}

func (processor controllerComponentProcessor) SupportsComponent(typ *core.Type) bool {
	returnType := core.GetFunctionFirstReturnType(typ)
	if returnType.Typ.Implements(core.GetType((*web.Controller)(nil)).Typ) {
		return true
	}
	return false
}

func (processor controllerComponentProcessor) ProcessComponent(typ *core.Type) error {
	inputTypes := core.GetFunctionInputTypes(typ)
	if inputTypes != nil && len(inputTypes) > 0 {
		return errors.New("the constructor of controller cannot take in parameters : " + typ.String())
	}
	return nil
}

type serviceComponentProcessor struct {
}

func newServiceComponentProcessor() serviceComponentProcessor {
	return serviceComponentProcessor{}
}

func (processor serviceComponentProcessor) SupportsComponent(typ *core.Type) bool {
	returnType := core.GetFunctionFirstReturnType(typ)
	if returnType.Typ.Implements(core.GetType((*context.Service)(nil)).Typ) {
		return true
	}
	return false
}

func (processor serviceComponentProcessor) ProcessComponent(typ *core.Type) error {
	inputTypes := core.GetFunctionInputTypes(typ)
	if inputTypes != nil && len(inputTypes) > 0 {
		return errors.New("the constructor of service cannot take in parameters : " + typ.String())
	}
	return nil
}

type repositoryComponentProcessor struct {
}

func newRepositoryComponentProcessor() repositoryComponentProcessor {
	return repositoryComponentProcessor{}
}

func (processor repositoryComponentProcessor) SupportsComponent(typ *core.Type) bool {
	returnType := core.GetFunctionFirstReturnType(typ)
	if returnType.Typ.Implements(core.GetType((*context.Repository)(nil)).Typ) {
		return true
	}
	return false
}

func (processor repositoryComponentProcessor) ProcessComponent(typ *core.Type) error {
	inputTypes := core.GetFunctionInputTypes(typ)
	if inputTypes != nil && len(inputTypes) > 0 {
		return errors.New("the constructor of repository cannot take in parameters : " + typ.String())
	}
	return nil
}
