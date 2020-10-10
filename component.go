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
	return nil
}

type peaComponentProcessor struct {
}

func newPeaComponentProcessor() peaComponentProcessor {
	return peaComponentProcessor{}
}

func (processor peaComponentProcessor) SupportsComponent(typ *core.Type) bool {
	return true
}

func (processor peaComponentProcessor) ProcessComponent(typ *core.Type) error {
	retType := core.GetFunctionFirstReturnType(typ)
	numField := core.GetNumField(retType)
	for index := 0; index < numField; index++ {
		structField := core.GetStructFieldByIndex(retType, index)
		_, hasInjectTag := structField.Tag.Lookup("inject")
		if hasInjectTag {
			if !core.IsExportedField(structField) {
				return errors.New("the tag of inject cannot be used on unexported fields : " + retType.String() + "->" + structField.Name)
			}
		}
	}
	return nil
}
