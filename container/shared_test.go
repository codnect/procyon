package container

import (
	"fmt"
	"github.com/procyon-projects/reflector"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstanceRegistry_Add(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")
}

func TestInstanceRegistry_AddReturnsErrorIfInstanceIsDuplicated(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	err = registry.Add("anyInstanceName", instance)

	assert.NotNil(t, err)
	assert.Equal(t, "container: instance with name anyInstanceName already exists", err.Error())
}

func TestInstanceRegistry_FindReturnsInstance(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	result, ok := registry.Find("anyInstanceName")
	assert.True(t, ok)
	assert.Equal(t, instance, result)
}

func TestInstanceRegistry_FindReturnsNilIfInstanceIsNotFound(t *testing.T) {
	registry := NewSharedInstances()
	result, ok := registry.Find("anyInstanceName")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestInstanceRegistry_ContainsReturnsTrueIfInstanceExistsInRegistry(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	ok := registry.Contains("anyInstanceName")
	assert.True(t, ok)
}

func TestInstanceRegistry_ContainsReturnsFalseIfInstanceIsNotFoundInRegistry(t *testing.T) {
	registry := NewSharedInstances()
	ok := registry.Contains("anyInstanceName")
	assert.False(t, ok)
}

func TestInstanceRegistry_InstanceNames(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	result := registry.InstanceNames()
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, result, []string{"anyInstanceName"})
}

func TestInstanceRegistry_FindByTypeReturnsErrorIfRequiredTypeIsNil(t *testing.T) {
	registry := NewSharedInstances()
	result, err := registry.FindByType(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "container: requiredType cannot be nil", err.Error())
	assert.Nil(t, result)
}

func TestInstanceRegistry_FindByTypeReturnsPointerInstanceIfRequiredTypeIsPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[*AnyType]())
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestInstanceRegistry_FindByTypeReturnsNonPointerInstanceIfRequiredTypeIsNotPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[AnyType]())
	assert.Nil(t, err)
	assert.Equal(t, *instance, result)
}

func TestInstanceRegistry_FindByTypeReturnsInstanceIfRequiredTypeIsInterface(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[fmt.Stringer]())
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestInstanceRegistry_FindByTypeReturnsErrorIfMultipleInstancesExistForRequiredType(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	var result any
	result, err = registry.FindByType(reflector.TypeOf[*AnyType]())
	assert.NotNil(t, err)
	assert.Equal(t, "container: instances cannot be distinguished for required type *AnyType", err.Error())
	assert.Nil(t, result)
}

func TestInstanceRegistry_FindAllByTypeReturnsPointerInstancesIfRequiredTypeIsPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(reflector.TypeOf[*AnyType]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{instance, anotherInstance}, result)
}

func TestInstanceRegistry_FindAllByTypeReturnsNonPointerInstancesIfRequiredTypeIsNotPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(reflector.TypeOf[AnyType]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{*instance, *anotherInstance}, result)
}

func TestInstanceRegistry_FindAllByTypeReturnsInstancesIfRequiredTypeIsInterface(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(reflector.TypeOf[fmt.Stringer]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{instance, anotherInstance}, result)
}

func TestInstanceRegistry_CountReturnsNumberOfRegisteredInstances(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	count := registry.Count()
	assert.Equal(t, 2, count)
}

func TestInstanceRegistry_OrElseGetReturnsInstanceIfThereIsAnyInstanceWithSpecifiedName(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()
	err := registry.Add("anyInstanceName", instance)

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

func TestInstanceRegistry_OrElseGetCreatesAndReturnsNewInstanceIfInstanceIsNotFound(t *testing.T) {
	instance := &AnyType{}
	registry := NewSharedInstances()

	result, err := registry.OrElseGet("anyInstanceName", func() (any, error) {
		return instance, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}
