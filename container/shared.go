package container

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/reflector"
	"sync"
)

type SharedInstances struct {
	instances              map[string]any
	instancesInPreparation map[string]struct{}
	typesOfInstances       map[string]reflector.Type
	muInstances            sync.RWMutex
}

func NewSharedInstances() *SharedInstances {
	return &SharedInstances{
		instances:              make(map[string]any),
		instancesInPreparation: make(map[string]struct{}),
		typesOfInstances:       make(map[string]reflector.Type),
		muInstances:            sync.RWMutex{},
	}
}

func (s *SharedInstances) Add(name string, instance any) error {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	if _, exists := s.instances[name]; exists {
		return fmt.Errorf("container: instance with name %s already exists", name)
	}

	s.instances[name] = instance
	s.typesOfInstances[name] = reflector.TypeOfAny(instance)
	return nil
}

func (s *SharedInstances) Find(name string) (any, bool) {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	if instance, exists := s.instances[name]; exists {
		return instance, true
	}

	return nil, false
}

func (s *SharedInstances) Contains(name string) bool {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	_, exists := s.instances[name]
	return exists
}

func (s *SharedInstances) InstanceNames() []string {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	names := make([]string, 0)
	for name, _ := range s.instances {
		names = append(names, name)
	}

	return names
}

func (s *SharedInstances) FindByType(requiredType reflector.Type) (any, error) {
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

func (s *SharedInstances) FindAllByType(requiredType reflector.Type) []any {
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

func (s *SharedInstances) OrElseGet(name string, supplier func() (any, error)) (any, error) {
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

func (s *SharedInstances) Count() int {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	return len(s.instances)
}

func (s *SharedInstances) putToPreparation(name string) error {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()

	if _, ok := s.instancesInPreparation[name]; ok {
		return fmt.Errorf("container: instance with name %s is currently in preparation, maybe it has got circular dependency cycle", name)
	}

	s.instancesInPreparation[name] = struct{}{}
	return nil
}

func (s *SharedInstances) removeFromPreparation(name string) {
	defer s.muInstances.Unlock()
	s.muInstances.Lock()
	delete(s.instancesInPreparation, name)
}
