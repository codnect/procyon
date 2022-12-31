package container

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/reflector"
	"sync"
)

type SharedInstances interface {
	Register(name string, instance any) error
	Find(name string) (any, bool)
	Contains(name string) bool
	InstanceNames() []string
	FindByType(requiredType reflector.Type) (any, error)
	FindAllByType(requiredType reflector.Type) []any
	OrElseGet(name string, supplier func() (any, error)) (any, error)
	Count() int
}

type sharedInstances struct {
	instances              map[string]any
	instancesInPreparation map[string]struct{}
	typesOfInstances       map[string]reflector.Type
	muInstances            sync.RWMutex
}

func NewSharedInstances() SharedInstances {
	return &sharedInstances{
		instances:              make(map[string]any),
		instancesInPreparation: make(map[string]struct{}),
		typesOfInstances:       make(map[string]reflector.Type),
		muInstances:            sync.RWMutex{},
	}
}

func (s *sharedInstances) Register(name string, instance any) error {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	if _, exists := s.instances[name]; exists {
		return fmt.Errorf("container: instance with name %s already exists", name)
	}

	s.instances[name] = instance
	s.typesOfInstances[name] = reflector.TypeOfAny(instance)
	return nil
}

func (s *sharedInstances) Find(name string) (any, bool) {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	if instance, exists := s.instances[name]; exists {
		return instance, true
	}

	return nil, false
}

func (s *sharedInstances) Contains(name string) bool {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	_, exists := s.instances[name]
	return exists
}

func (s *sharedInstances) InstanceNames() []string {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	names := make([]string, 0)
	for name, _ := range s.instances {
		names = append(names, name)
	}

	return names
}

func (s *sharedInstances) FindByType(requiredType reflector.Type) (any, error) {
	if requiredType == nil {
		return nil, errors.New("container: requiredType cannot be nil")
	}

	instances := s.FindAllByType(requiredType)
	if len(instances) > 1 {
		return nil, fmt.Errorf("container: instances cannot be distinguished for required type %s", requiredType.Name())
	}

	if len(instances) == 0 {
		return nil, &notFoundError{
			ErrorString: fmt.Sprintf("container: not found any instance of type %s", requiredType.Name()),
		}
	}

	return instances[0], nil
}

func (s *sharedInstances) FindAllByType(requiredType reflector.Type) []any {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	instances := make([]any, 0)

	for name, typ := range s.typesOfInstances {

		if typ.CanConvert(requiredType) {
			instances = append(instances, s.instances[name])
		} else if reflector.IsPointer(typ) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
			ptrType := reflector.ToPointer(typ)

			if ptrType.Elem().CanConvert(requiredType) {
				val, err := ptrType.Elem().Value()

				if err == nil {
					instances = append(instances, val)
				}
			}
		}

	}

	return instances
}

func (s *sharedInstances) OrElseGet(name string, supplier func() (any, error)) (any, error) {
	instance, ok := s.Find(name)

	if ok {
		return instance, nil
	}

	err := s.putToPreparation(name)

	if err != nil {
		return nil, err
	}

	defer func() {
		s.removeFromPreparation(name)
	}()

	instance, err = supplier()

	if err != nil {
		return nil, err
	}

	s.instances[name] = instance
	s.typesOfInstances[name] = reflector.TypeOfAny(instance)

	return instance, nil
}

func (s *sharedInstances) Count() int {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	return len(s.instances)
}

func (s *sharedInstances) putToPreparation(name string) error {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	if _, ok := s.instancesInPreparation[name]; ok {
		return fmt.Errorf("container: instance with name %s is currently in preparation, maybe it has got circular dependency cycle", name)
	}

	s.instancesInPreparation[name] = struct{}{}
	return nil
}

func (s *sharedInstances) removeFromPreparation(name string) {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()
	delete(s.instancesInPreparation, name)
}
