package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstanceRegistry_Add(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")
}

func TestInstanceRegistry_AddReturnsErrorIfInstanceIsDuplicated(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
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
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	result, ok := registry.Find("anyInstanceName")
	assert.True(t, ok)
	assert.Equal(t, instance, result)
}

func TestInstanceRegistry_FindReturnsNilIfInstanceIsNotFound(t *testing.T) {
	registry := NewInstanceRegistry()
	result, ok := registry.Find("anyInstanceName")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestInstanceRegistry_ContainsReturnsTrueIfInstanceExistsInRegistry(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	ok := registry.Contains("anyInstanceName")
	assert.True(t, ok)
}

func TestInstanceRegistry_ContainsReturnsFalseIfInstanceIsNotFoundInRegistry(t *testing.T) {
	registry := NewInstanceRegistry()
	ok := registry.Contains("anyInstanceName")
	assert.False(t, ok)
}

func TestInstanceRegistry_InstanceNames(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
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
	registry := NewInstanceRegistry()
	result, err := registry.FindByType(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "container: requiredType cannot be nil", err.Error())
	assert.Nil(t, result)
}

func TestInstanceRegistry_FindByTypeReturnsPointerInstanceIfRequiredTypeIsPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(TypeOf[*AnyType]())
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestInstanceRegistry_FindByTypeReturnsNonPointerInstanceIfRequiredTypeIsNotPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(TypeOf[AnyType]())
	assert.Nil(t, err)
	assert.Equal(t, *instance, result)
}

func TestInstanceRegistry_FindByTypeReturnsInstanceIfRequiredTypeIsInterface(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	var result any
	result, err = registry.FindByType(TypeOf[fmt.Stringer]())
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestInstanceRegistry_FindByTypeReturnsErrorIfMultipleInstancesExistForRequiredType(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
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
	result, err = registry.FindByType(TypeOf[*AnyType]())
	assert.NotNil(t, err)
	assert.Equal(t, "container: instances cannot be distinguished for required type *AnyType", err.Error())
	assert.Nil(t, result)
}

func TestInstanceRegistry_FindAllByTypeReturnsPointerInstancesIfRequiredTypeIsPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(TypeOf[*AnyType]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{instance, anotherInstance}, result)
}

func TestInstanceRegistry_FindAllByTypeReturnsNonPointerInstancesIfRequiredTypeIsNotPointer(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(TypeOf[AnyType]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{*instance, *anotherInstance}, result)
}

func TestInstanceRegistry_FindAllByTypeReturnsInstancesIfRequiredTypeIsInterface(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
	err := registry.Add("anyInstanceName", instance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anyInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anyInstanceName")

	anotherInstance := &AnyType{}
	err = registry.Add("anotherInstanceName", anotherInstance)

	assert.Nil(t, err)
	assert.Contains(t, registry.instances, "anotherInstanceName")
	assert.Contains(t, registry.typesOfInstances, "anotherInstanceName")

	result := registry.FindAllByType(TypeOf[fmt.Stringer]())
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, []any{instance, anotherInstance}, result)
}

func TestInstanceRegistry_CountReturnsNumberOfRegisteredInstances(t *testing.T) {
	instance := &AnyType{}
	registry := NewInstanceRegistry()
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
	registry := NewInstanceRegistry()
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
	registry := NewInstanceRegistry()

	result, err := registry.OrElseGet("anyInstanceName", func() (any, error) {
		return instance, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, instance, result)
}

func TestDefinitionRegistry_AddReturnsErrorIfDefinitionIsNil(t *testing.T) {
	registry := NewDefinitionRegistry()
	err := registry.Add(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "container: definition should not be nil", err.Error())
}

func TestDefinitionRegistry_AddReturnNilIfDefinitionIsAddedSuccessfully(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")
}

func TestDefinitionRegistry_AddReturnsErrorIfDefinitionIsDuplicated(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	err = registry.Add(def)

	assert.NotNil(t, err)
	assert.Equal(t, "container: definition with name anyType already exists", err.Error())
}

func TestDefinitionRegistry_RemoveReturnsErrorIfDefinitionIsNotFound(t *testing.T) {
	registry := NewDefinitionRegistry()

	err := registry.Remove("anyType")

	assert.NotNil(t, err)
	assert.Equal(t, "container: no found definition with name anyType", err.Error())
}

func TestDefinitionRegistry_RemoveDeletesDefinitionFromRegistry(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	err = registry.Remove("anyType")

	assert.Nil(t, err)
	assert.NotContains(t, registry.definitions, "anyType")
}

func TestDefinitionRegistry_ContainsReturnsTrueIfDefinitionExistsInRegistry(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	assert.True(t, registry.Contains("anyType"))
}

func TestDefinitionRegistry_ContainsReturnsFalseIfDefinitionDoesNotExistInRegistry(t *testing.T) {
	registry := NewDefinitionRegistry()
	assert.False(t, registry.Contains("anyType"))
}

func TestDefinitionRegistry_FindReturnsDefinitionIfItExistsInRegistry(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	result, ok := registry.Find("anyType")
	assert.True(t, ok)
	assert.Equal(t, def, result)
}

func TestDefinitionRegistry_FindReturnsNilIfDefinitionDoesNotExistInRegistry(t *testing.T) {
	registry := NewDefinitionRegistry()

	result, ok := registry.Find("anyType")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestDefinitionRegistry_DefinitionsReturnsRegisteredDefinitions(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	defs := registry.Definitions()
	assert.NotNil(t, defs)
	assert.Equal(t, []*Definition{def}, defs)
}

func TestDefinitionRegistry_DefinitionNamesReturnsRegisteredDefinitionNames(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	names := registry.DefinitionNames()
	assert.NotNil(t, names)
	assert.Equal(t, []string{"anyType"}, names)
}

func TestDefinitionRegistry_DefinitionNamesByTypeReturnsEmptyNamesIfRequiredTypeIsNil(t *testing.T) {
	registry := NewDefinitionRegistry()
	names := registry.DefinitionNamesByType(nil)
	assert.NotNil(t, names)
	assert.Len(t, names, 0)
}

func TestDefinitionRegistry_DefinitionNamesByTypeReturnsRegisteredDefinitionNamesBasedOnType(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	names := registry.DefinitionNamesByType(TypeOf[*AnyType]())
	assert.NotNil(t, names)
	assert.Equal(t, []string{"anyType"}, names)

	names = registry.DefinitionNamesByType(TypeOf[AnyType]())
	assert.NotNil(t, names)
	assert.Equal(t, []string{"anyType"}, names)

	names = registry.DefinitionNamesByType(TypeOf[fmt.Stringer]())
	assert.NotNil(t, names)
	assert.Equal(t, []string{"anyType"}, names)
}

func TestDefinitionRegistry_CountReturnsNumberOfDefinitions(t *testing.T) {
	registry := NewDefinitionRegistry()
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	count := registry.Count()
	assert.Equal(t, 1, count)
}
