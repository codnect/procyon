package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type AnySliceType struct {
	anyField int
}

func AnySliceConstructFunction(t []*DependencyType) *AnySliceType {
	return &AnySliceType{}
}

func TestContainer_Start(t *testing.T) {
	c := New()
	err := c.Start()
	assert.Nil(t, err)
}

func TestContainer_DefinitionRegistry(t *testing.T) {
	c := New()
	definitionRegistry := c.DefinitionRegistry()
	assert.NotNil(t, definitionRegistry)
}

func TestContainer_InstanceRegistry(t *testing.T) {
	c := New()
	instanceRegistry := c.InstanceRegistry()
	assert.NotNil(t, instanceRegistry)
}

func TestContainer_Hooks(t *testing.T) {
	c := New()
	hooks := c.Hooks()
	assert.NotNil(t, hooks)
}

func TestContainer_Get(t *testing.T) {
	c := New()
	def, err := MakeDefinition(AnyConstructFunction, OptionalAt(0))
	assert.Nil(t, err)
	assert.NotNil(t, def)

	c.DefinitionRegistry().Add(def)

	var instance any
	instance, err = c.Get("anyType")
	assert.Nil(t, err)
	assert.NotNil(t, instance)
}

func TestContainer_GetSliceDependency(t *testing.T) {
	c := New()
	def, err := MakeDefinition(AnySliceConstructFunction)
	assert.Nil(t, err)
	assert.NotNil(t, def)

	c.DefinitionRegistry().Add(def)
	c.InstanceRegistry().Add("test1", &DependencyType{})
	c.InstanceRegistry().Add("test2", &DependencyType{})

	var instance any
	instance, err = c.Get("anySliceType")
	assert.Nil(t, err)
	assert.NotNil(t, instance)
}
