package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	Register(AnyConstructFunction, Optional[*DependencyType]())
	assert.Len(t, definitions, 1)
	assert.Contains(t, definitions, "anyType")
	definitions = map[string]*Definition{}
}

func TestRegister_PanicsIfDefinitionIsAlreadyRegisteredWithSameName(t *testing.T) {
	Register(AnyConstructFunction)
	assert.Len(t, definitions, 1)
	assert.Contains(t, definitions, "anyType")

	assert.Panics(t, func() {
		Register(AnyConstructFunction)
		assert.Len(t, definitions, 1)
		assert.Contains(t, definitions, "anyType")
	})

	definitions = map[string]*Definition{}
}

func TestRegister_PanicsIfThereIsProblemWithDefinition(t *testing.T) {
	assert.Panics(t, func() {
		Register(AnyConstructFunction, OptionalAt(2))
		assert.Len(t, definitions, 0)
		assert.NotContains(t, definitions, "anyType")
	})

	definitions = map[string]*Definition{}
}

func TestRegister_WithDifferentNames(t *testing.T) {
	Register(AnyConstructFunction)
	assert.Len(t, definitions, 1)
	assert.Contains(t, definitions, "anyType")

	Register(AnyConstructFunction, Name("anotherName"))
	assert.Len(t, definitions, 2)
	assert.Contains(t, definitions, "anyType")
	assert.Contains(t, definitions, "anotherName")

	definitions = map[string]*Definition{}
}

func TestCopyDefinitions_ReturnsCopyOfRegisteredDefinitions(t *testing.T) {
	Register(AnyConstructFunction)
	assert.Len(t, definitions, 1)
	assert.Contains(t, definitions, "anyType")

	copiedDefinitions := copyDefinitions()
	assert.Equal(t, definitions, copiedDefinitions)

	definitions = map[string]*Definition{}
}
