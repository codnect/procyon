package container

import (
	"codnect.io/reflector"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSharedInstances_Register(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")
}

func TestSharedInstances_AddReturnsErrorIfInstanceIsDuplicated(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	err = registry.Register("anyInstanceName", instance)

	assert.NotNil(t, err)
	assert.Equal(t, "container: instance with name anyInstanceName already exists", err.Error())
}

func TestSharedInstances_FindReturnsInstance(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	result, ok := registry.Find("anyInstanceName")
	assert.True(t, ok)
	assert.Equal(t, instance, result)
}

func TestSharedInstances_FindReturnsNilIfInstanceIsNotFound(t *testing.T) {
	registry := NewSharedInstances()
	result, ok := registry.Find("anyInstanceName")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestSharedInstances_ContainsReturnsTrueIfInstanceExistsInRegistry(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	ok := registry.Contains("anyInstanceName")
	assert.True(t, ok)
}

func TestSharedInstances_ContainsReturnsFalseIfInstanceIsNotFoundInRegistry(t *testing.T) {
	registry := NewSharedInstances()
	ok := registry.Contains("anyInstanceName")
	assert.False(t, ok)
}

func TestSharedInstances_InstanceNames(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	result := registry.InstanceNames()
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, result, []string{"anyInstanceName"})
}

func TestSharedInstances_FindByTypeReturnsErrorIfRequiredTypeIsNil(t *testing.T) {
	registry := NewSharedInstances()
	result, err := registry.FindByType(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "container: requiredType cannot be nil", err.Error())
	assert.Nil(t, result)
}

func TestSharedInstances_FindByTypeReturnsPointerInstanceIfRequiredTypeIsPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[*AnyType]())
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestSharedInstances_FindByTypeReturnsNonPointerInstanceIfRequiredTypeIsNotPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[AnyType]())
	assert.Nil(t, err)
	assert.Equal(t, *instance, result)
}

func TestSharedInstances_FindByTypeReturnsInstanceIfRequiredTypeIsInterface(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[fmt.Stringer]())
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestSharedInstances_FindByTypeReturnsErrorIfMultipleInstancesExistForRequiredType(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Register("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[*AnyType]())
	assert.NotNil(t, err)
	assert.Equal(t, "container: instances cannot be distinguished for required type *AnyType", err.Error())
	assert.Nil(t, result)
}

func TestSharedInstances_FindAllByTypeReturnsPointerInstancesIfRequiredTypeIsPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Register("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(reflector.TypeOf[*AnyType]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{instance, anotherInstance}, result)
}

func TestSharedInstances_FindAllByTypeReturnsNonPointerInstancesIfRequiredTypeIsNotPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Register("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(reflector.TypeOf[AnyType]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{*instance, *anotherInstance}, result)
}

func TestSharedInstances_FindAllByTypeReturnsInstancesIfRequiredTypeIsInterface(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Register("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(reflector.TypeOf[fmt.Stringer]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{instance, anotherInstance}, result)
}

func TestSharedInstances_CountReturnsNumberOfRegisteredInstances(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Register("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	count := registry.Count()
	assert.Equal(t, 2, count)
}

func TestSharedInstances_OrElseGetReturnsInstanceIfThereIsAnyInstanceWithSpecifiedName(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances().(*sharedInstances)
	err := registry.Register("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.OrElseGet("anyInstanceName", func() (any, error) {
		return &AnyType{}, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestSharedInstances_OrElseGetCreatesAndReturnsNewInstanceIfInstanceIsNotFound(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()

	result, err := registry.OrElseGet("anyInstanceName", func() (any, error) {
		return instance, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}
