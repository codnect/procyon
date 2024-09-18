package container

import (
	"codnect.io/procyon/component/filter"
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"
)

type Container interface {
	GetObject(ctx context.Context, filters ...filter.Filter) (any, error)
	ListObjects(ctx context.Context, filters ...filter.Filter) []any
	ContainsObject(name string) bool
	IsSingleton(name string) bool
	IsPrototype(name string) bool
	Definitions() DefinitionRegistry
	Singletons() SingletonRegistry
	Scopes() ScopeRegistry
	AddObjectProcessor(processor ObjectProcessor) error
	ObjectProcessorCount() int
}

type defaultContainer struct {
	definitions *objectDefinitionRegistry
	singletons  *singletonObjectRegistry
	scopes      *simpleScopeRegistry

	processors       []ObjectProcessor
	typesOfProcessor map[string]struct{}
	postProcessorMu  sync.RWMutex

	running bool
	mu      sync.RWMutex
}

// New function creates a new default container.
func New() Container {
	return &defaultContainer{
		definitions:      newObjectDefinitionRegistry(),
		singletons:       newSingletonObjectRegistry(),
		scopes:           newSimpleScopeRegistry(),
		processors:       make([]ObjectProcessor, 0),
		typesOfProcessor: map[string]struct{}{},
	}
}

// GetObject method gets an object from the container that matches the provided filters.
func (c *defaultContainer) GetObject(ctx context.Context, filters ...filter.Filter) (any, error) {
	if len(filters) == 0 {
		return nil, ErrNoFilterProvided
	}

	ctx = withObjectCreationState(ctx)

	candidate, err := c.Singletons().Find(filters...)
	if err == nil {
		return candidate, nil
	} else if !errors.Is(err, ErrObjectNotFound) {
		return nil, err
	}

	var definition *Definition
	definition, err = c.Definitions().Find(filters...)

	if err != nil {
		return nil, err
	}

	objectName := definition.Name()

	if definition.IsSingleton() {
		var object any
		object, err = c.Singletons().OrElseCreate(objectName, func(ctx context.Context) (any, error) {

			if log.IsDebugEnabled() {
				log.D(ctx, "Creating singleton object of type '{}'", definition.Type().String())
			}

			return c.createObject(ctx, definition, nil)
		})

		return object, err
	} else if definition.IsPrototype() {
		prototypeHolder := objectCreationStateFromContext(ctx)
		err = prototypeHolder.putToPreparation(objectName)

		if err != nil {
			return nil, err
		}

		defer prototypeHolder.removeFromPreparation(objectName)
		return c.createObject(ctx, definition, nil)
	}

	if strings.TrimSpace(definition.Scope()) == "" {
		return nil, fmt.Errorf("no scope name for required type %s", definition.Type().Name())
	}

	var scope Scope
	scope, err = c.scopes.Find(definition.Scope())

	if err != nil {
		return nil, err
	}

	return scope.GetObject(ctx, objectName, func(ctx context.Context) (any, error) {
		scopeHolder := objectCreationStateFromContext(ctx)
		err = scopeHolder.putToPreparation(objectName)

		if err != nil {
			return nil, err
		}

		defer scopeHolder.removeFromPreparation(objectName)
		return c.createObject(ctx, definition, nil)
	})
}

// ListObjects method lists all objects in the container that match the provided filters.
func (c *defaultContainer) ListObjects(ctx context.Context, filters ...filter.Filter) []any {
	objectList := make([]any, 0)
	singletonNames := c.singletons.Names()
	objectList = append(objectList, c.singletons.List(filters...)...)

	for _, definition := range c.definitions.List(filters...) {
		if (definition.IsSingleton() && !slices.Contains(singletonNames, definition.Name())) || !definition.IsSingleton() {
			object, err := c.GetObject(ctx, filter.ByName(definition.Name()))

			if err != nil {
				continue
			}

			objectList = append(objectList, object)
		}
	}

	return objectList
}

// ContainsObject method checks if an object exists in the container.
func (c *defaultContainer) ContainsObject(name string) bool {
	return c.singletons.Contains(name)
}

// IsSingleton method checks if an object is a singleton.
func (c *defaultContainer) IsSingleton(name string) bool {
	definition, ok := c.definitions.FindFirst(filter.ByName(name))
	return ok && definition.IsSingleton()
}

// IsPrototype method checks if an object is a prototype.
func (c *defaultContainer) IsPrototype(name string) bool {
	definition, ok := c.definitions.FindFirst(filter.ByName(name))
	return ok && definition.IsPrototype()
}

// Definitions method returns the definition registry of the container.
func (c *defaultContainer) Definitions() DefinitionRegistry {
	return c.definitions
}

// Singletons method returns the singleton registry of the container.
func (c *defaultContainer) Singletons() SingletonRegistry {
	return c.singletons
}

// Scopes method returns the scope registry of the container.
func (c *defaultContainer) Scopes() ScopeRegistry {
	return c.scopes
}

// AddObjectProcessor method adds an object processor to the container.
func (c *defaultContainer) AddObjectProcessor(processor ObjectProcessor) error {
	if processor == nil {
		return errors.New("nil processor")
	}

	defer c.postProcessorMu.Unlock()
	c.postProcessorMu.Lock()

	typ := reflect.TypeOf(processor)
	typeName := fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())

	if _, ok := c.typesOfProcessor[typeName]; ok {
		return fmt.Errorf("processor '%s' is already registered", typeName)
	}

	c.typesOfProcessor[typeName] = struct{}{}
	c.processors = append(c.processors, processor)
	return nil
}

// ObjectProcessorCount method returns the number of object processors in the container.
func (c *defaultContainer) ObjectProcessorCount() int {
	defer c.postProcessorMu.Unlock()
	c.postProcessorMu.Lock()
	return len(c.processors)
}

// createObject method creates an object based on a definition and arguments.
func (c *defaultContainer) createObject(ctx context.Context, definition *Definition, args []any) (object any, err error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	if definition == nil {
		return nil, errors.New("nil definition")
	}

	objectConstructor := definition.Constructor()
	argsCount := len(objectConstructor.Arguments())

	if argsCount != 0 && len(args) == 0 {
		var resolvedArguments []any
		resolvedArguments, err = c.resolveArguments(ctx, objectConstructor.Arguments())

		if err != nil {
			return nil, err
		}

		var results []any
		results, err = objectConstructor.Invoke(resolvedArguments...)

		if err != nil {
			return nil, err
		}

		resultType := reflect.TypeOf(results[0])
		resultValue := reflect.ValueOf(results[0])
		if (resultType.Kind() == reflect.Pointer || resultType.Kind() == reflect.Interface) && resultValue.IsZero() {
			return nil, fmt.Errorf("Constructor function '%s' returns nil", objectConstructor.Name())
		}

		object = results[0]
	} else if (argsCount == 0 && len(args) == 0) || (len(args) != 0 && argsCount == len(args)) {
		var results []any
		results, err = objectConstructor.Invoke(args...)

		if err != nil {
			return nil, err
		}

		resultType := reflect.TypeOf(results[0])
		resultValue := reflect.ValueOf(results[0])

		if (resultType.Kind() == reflect.Pointer || resultType.Kind() == reflect.Interface) && resultValue.IsZero() {
			return nil, fmt.Errorf("Constructor function '%s' returns nil", objectConstructor.Name())
		}

		object = results[0]
	} else {
		return nil, fmt.Errorf("the number of provided arguments is wrong for definition '%s'", definition.Name())
	}

	return c.initialize(ctx, object)
}

// resolveArguments method resolves the arguments for a constructor.
func (c *defaultContainer) resolveArguments(ctx context.Context, args []ConstructorArgument) ([]any, error) {
	arguments := make([]any, 0)

	for _, arg := range args {

		if arg.Type().Kind() == reflect.Slice {
			sliceType := arg.Type()
			sliceVal := reflect.MakeSlice(sliceType, 0, 0)

			objectList := c.ListObjects(ctx, filter.ByType(sliceType.Elem()))
			for _, object := range objectList {
				sliceVal = reflect.Append(sliceVal, reflect.ValueOf(object))
			}

			arguments = append(arguments, sliceVal.Interface())
			continue
		}

		var (
			object any
			err    error
		)

		/*
			resolvableInstance, exists := m.getResolvableInstance(arg.Type())
			if exists {
				arguments = append(arguments, resolvableInstance)
				continue
			}
		*/

		if arg.Name() != "" {
			object, err = c.GetObject(ctx, filter.ByName(arg.Name()))
		} else {
			object, err = c.GetObject(ctx, filter.ByType(arg.Type()))
		}

		if err != nil {
			argKind := arg.Type().Kind()

			if errors.Is(err, ErrObjectNotFound) && argKind != reflect.Pointer && argKind != reflect.Interface {
				var val reflect.Value
				val = reflect.New(arg.Type())

				object = val.Elem()
				arguments = append(arguments, object)
				continue
			}

			if !arg.IsOptional() && err != nil {
				return nil, err
			} else if arg.IsOptional() && err != nil {
				arguments = append(arguments, nil)
			}
		} else {
			arguments = append(arguments, object)
		}
	}

	return arguments, nil
}

// initialize method initializes an object.
func (c *defaultContainer) initialize(ctx context.Context, object any) (any, error) {
	result, err := c.applyProcessorsBeforeInit(ctx, object)
	if err != nil {
		return nil, err
	}

	if initialization, ok := object.(Initialization); ok {
		err = initialization.DoInit(ctx)

		if err != nil {
			return nil, err
		}
	}

	result, err = c.applyProcessorsAfterInit(ctx, object)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// applyProcessorsBeforeInit method applies the object processors before initialization.
func (c *defaultContainer) applyProcessorsBeforeInit(ctx context.Context, object any) (any, error) {
	for _, processor := range c.processors {
		result, err := processor.ProcessBeforeInit(ctx, object)

		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, fmt.Errorf("'%s' returns nil object from ProcessBeforeInit", reflect.TypeOf(processor).Name())
		}

		object = result
	}

	return object, nil
}

// applyProcessorsAfterInit method applies the object processors after initialization.
func (c *defaultContainer) applyProcessorsAfterInit(ctx context.Context, object any) (any, error) {
	for _, processor := range c.processors {
		result, err := processor.ProcessAfterInit(ctx, object)

		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, fmt.Errorf("'%s' returns nil object from ProcessAfterInit", reflect.TypeOf(processor).Name())
		}

		object = result
	}

	return object, nil
}

// loadObjectProcessors method loads the object processors.
func (c *defaultContainer) loadObjectProcessors(ctx context.Context) error {

	/*
		postProcessors := c.Definitions().List(filter.ByTypeOf[component.ObjectProcessor]())

		checker := component.newProcessorChecker(c, c.ObjectProcessorCount()+len(postProcessors)+1)
		_ = c.AddObjectProcessor(checker)

		for _, processorDefinition := range postProcessors {
			processor, err := c.GetObject(ctx, filter.ByName(processorDefinition.Name()))

			if err != nil {
				return err
			}

			err = c.AddObjectProcessor(processor.(component.ObjectProcessor))

			if err != nil {
				return err
			}
		}*/

	return nil
}
