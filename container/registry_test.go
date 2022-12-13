package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefinitionRegistry_AddReturnsErrorIfDefinitionIsNil(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
	err := registry.Add(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "container: definition should not be nil", err.Error())
}

func TestDefinitionRegistry_AddReturnNilIfDefinitionIsAddedSuccessfully(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")
}

func TestDefinitionRegistry_AddReturnsErrorIfDefinitionIsDuplicated(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
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
	registry := NewDefinitionRegistry(nil)

	err := registry.Remove("anyType")

	assert.NotNil(t, err)
	assert.Equal(t, "container: no found definition with name anyType", err.Error())
}

func TestDefinitionRegistry_RemoveDeletesDefinitionFromRegistry(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
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
	registry := NewDefinitionRegistry(nil)
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	assert.True(t, registry.Contains("anyType"))
}

func TestDefinitionRegistry_ContainsReturnsFalseIfDefinitionDoesNotExistInRegistry(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
	assert.False(t, registry.Contains("anyType"))
}

func TestDefinitionRegistry_FindReturnsDefinitionIfItExistsInRegistry(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
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
	registry := NewDefinitionRegistry(nil)

	result, ok := registry.Find("anyType")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestDefinitionRegistry_DefinitionsReturnsRegisteredDefinitions(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
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
	registry := NewDefinitionRegistry(nil)
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
	registry := NewDefinitionRegistry(nil)
	names := registry.DefinitionNamesByType(nil)
	assert.NotNil(t, names)
	assert.Len(t, names, 0)
}

func TestDefinitionRegistry_DefinitionNamesByTypeReturnsRegisteredDefinitionNamesBasedOnType(t *testing.T) {
	registry := NewDefinitionRegistry(nil)
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
	registry := NewDefinitionRegistry(nil)
	def, err := MakeDefinition(AnyConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	err = registry.Add(def)

	assert.Nil(t, err)
	assert.Contains(t, registry.definitions, "anyType")

	count := registry.Count()
	assert.Equal(t, 1, count)
}
