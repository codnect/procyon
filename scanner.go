package procyon

import (
	"fmt"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
)

type componentScanner struct {
}

func newComponentScanner() componentScanner {
	return componentScanner{}
}

func (scanner componentScanner) scan(contextId string, logger context.Logger) (int, error) {
	processors, err := scanner.getProcessorInstances()
	if err != nil {
		return -1, nil
	}
	var componentCount = 0
	result := core.VisitComponentTypes(func(componentName string, componentType *core.Type) error {
		logger.T(contextId, fmt.Sprintf("Component : %s", componentName))
		err := scanner.checkComponent(componentType, processors)
		if err != nil {
			return err
		}
		componentCount++
		return nil
	})
	return componentCount, result
}

func (scanner componentScanner) checkComponent(componentType *core.Type, processors []interface{}) (err error) {
	for _, processorInstance := range processors {
		if processor, ok := processorInstance.(core.ComponentProcessor); ok {
			if processor.SupportsComponent(componentType) {
				err = processor.ProcessComponent(componentType)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (scanner componentScanner) getProcessorInstances() ([]interface{}, error) {
	var instances []interface{}
	result := core.VisitComponentProcessors(func(processorName string, processorType *core.Type) error {
		instance, err := peas.CreateInstance(processorType, []interface{}{})
		if err != nil {
			return err
		}
		instances = append(instances, instance)
		return nil
	})
	return instances, result
}