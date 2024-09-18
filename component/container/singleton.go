package container

import (
	"codnect.io/procyon/component/filter"
	"context"
	"reflect"
	"sync"
)

// SingletonRegistry interface provides methods for managing and accessing singleton objects
type SingletonRegistry interface {
	// Register adds a new singleton object to the registry
	Register(name string, object any) error
	// Remove removes a singleton object from the registry
	Remove(name string) error
	// Find returns the singleton object that matches the provided filters
	Find(filters ...filter.Filter) (any, error)
	//	FindFirst returns the first singleton object that matches the provided filters
	FindFirst(filters ...filter.Filter) (any, bool)
	// List returns a list of singleton objects that match the provided filters
	List(filters ...filter.Filter) []any
	// OrElseCreate returns the singleton object with the provided name if it exists, otherwise
	// it creates a new object using the provided provider function
	OrElseCreate(name string, provider ObjectProviderFunc) (any, error)
	// Contains checks if a singleton object with the provided name exists
	Contains(name string) bool
	// Names returns the names of all singleton objects
	Names() []string
	// Count returns the number of singleton objects
	Count() int
}

// singletonObjectRegistry struct represents a registry for singleton objects.
type singletonObjectRegistry struct {
	// singletonObjects is a map that stores singleton objects by their names
	singletonObjects map[string]any
	// typesOfSingletonObjects is a map that stores the types of singleton objects
	typesOfSingletonObjects map[string]reflect.Type
	objectCreationContext   context.Context
	// muSingletonObjects is a mutex that protects the singletonObjects map
	muSingletonObjects sync.RWMutex
}

// newSingletonObjectRegistry creates a new singleton object registry.
func newSingletonObjectRegistry() *singletonObjectRegistry {
	return &singletonObjectRegistry{
		singletonObjects:        make(map[string]any),
		typesOfSingletonObjects: make(map[string]reflect.Type),
		objectCreationContext:   withObjectCreationState(context.Background()),
	}
}

// Register adds a new singleton object to the registry
func (r *singletonObjectRegistry) Register(name string, object any) error {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	if _, exists := r.singletonObjects[name]; exists {
		return ErrObjectAlreadyExists
	}

	r.singletonObjects[name] = object
	r.typesOfSingletonObjects[name] = reflect.TypeOf(object)
	return nil
}

// Remove removes a singleton object from the registry
func (r *singletonObjectRegistry) Remove(name string) error {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	if _, exists := r.singletonObjects[name]; !exists {
		return ErrObjectNotFound
	}

	delete(r.singletonObjects, name)
	delete(r.typesOfSingletonObjects, name)
	return nil
}

// Find returns the singleton object that matches the provided filters
func (r *singletonObjectRegistry) Find(filters ...filter.Filter) (any, error) {
	if len(filters) == 0 {
		return nil, ErrNoFilterProvided
	}

	objectList := r.List(filters...)

	if len(objectList) > 1 {
		return nil, ErrMultipleObjectsFound
	} else if len(objectList) == 0 {
		return nil, ErrObjectNotFound
	}

	return objectList[0], nil
}

// FindFirst returns the first singleton object that matches the provided filters
func (r *singletonObjectRegistry) FindFirst(filters ...filter.Filter) (any, bool) {
	objectList := r.List(filters...)

	if len(objectList) == 0 {
		return nil, false
	}

	return objectList[0], true
}

// List returns a list of singleton objects that match the provided filters
func (r *singletonObjectRegistry) List(filters ...filter.Filter) []any {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	filterOpts := filter.Of(filters...)
	objectList := make([]any, 0)

	for objectName, objectType := range r.typesOfSingletonObjects {

		if filterOpts.Name != "" && filterOpts.Name != objectName {
			continue
		}

		if filterOpts.Type == nil {
			objectList = append(objectList, r.singletonObjects[objectName])
			continue
		}

		if convertibleTo(objectType, filterOpts.Type) {
			objectList = append(objectList, r.singletonObjects[objectName])
		}
	}

	return objectList
}

// OrElseCreate returns the singleton object with the provided name if it exists, otherwise
// it creates a new object using the provided provider function
func (r *singletonObjectRegistry) OrElseCreate(name string, provider ObjectProviderFunc) (any, error) {
	object, err := r.Find(filter.ByName(name))

	if err == nil {
		return object, nil
	}

	state := objectCreationStateFromContext(r.objectCreationContext)
	err = state.putToPreparation(name)

	if err != nil {
		return nil, err
	}

	defer state.removeFromPreparation(name)
	object, err = provider(r.objectCreationContext)

	if err != nil {
		return nil, err
	}

	r.singletonObjects[name] = object
	r.typesOfSingletonObjects[name] = reflect.TypeOf(object)

	return object, nil
}

// Contains checks if a singleton object with the provided name exists
func (r *singletonObjectRegistry) Contains(name string) bool {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	_, exists := r.singletonObjects[name]
	return exists
}

// Names returns the names of all singleton objects
func (r *singletonObjectRegistry) Names() []string {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	names := make([]string, 0)
	for name := range r.singletonObjects {
		names = append(names, name)
	}

	return names
}

// Count returns the number of singleton objects
func (r *singletonObjectRegistry) Count() int {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	return len(r.singletonObjects)
}
