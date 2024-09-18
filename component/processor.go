package component

/*
type processorChecker struct {
	container      *Container
	processorCount int
}

func newProcessorChecker(container *Container, processorCount int) processorChecker {
	return processorChecker{
		container:      container,
		processorCount: processorCount,
	}
}

func (c processorChecker) ProcessBeforeInit(ctx context.Context, object any) (any, error) {
	return object, nil
}

func (c processorChecker) ProcessAfterInit(ctx context.Context, object any) (any, error) {
	/*if _, ok := object.(ObjectProcessor); !ok && c.container.ObjectProcessorCount() < c.processorCount {
		typ := reflect.TypeOf(object)
		typeName := fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
		log.I(ctx, "Component '{}' is not eligible for ObjectProcessor", typeName)
	}

	return object, nil
}
*/
