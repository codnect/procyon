package container

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"testing"
)

type AnyType struct {
}

type DependencyType struct {
}

func AnyConstructFunction(t *DependencyType) *AnyType {
	return &AnyType{}
}

func TestMakeDefinition(t *testing.T) {
	def, err := MakeDefinition(AnyConstructFunction, Qualifier[*DependencyType]("anyDependencyType"))
	assert.Nil(t, err)

	assert.Equal(t, "anyType", def.Name())
	assert.Equal(t, "*AnyType", def.Type().Name())
	assert.Equal(t, runtime.FuncForPC(reflect.ValueOf(AnyConstructFunction).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(def.Constructor()).Pointer()).Name())

	inputs := def.Inputs()
	assert.NotNil(t, inputs)
	assert.Len(t, inputs, 1)

	assert.Equal(t, 0, inputs[0].Index())
	assert.False(t, inputs[0].IsOptional())
	assert.Equal(t, "anyDependencyType", inputs[0].Name())
	assert.Equal(t, "*DependencyType", inputs[0].Type().Name())

	assert.Equal(t, false, def.IsPrototype())
	assert.Equal(t, true, def.IsShared())
	assert.Equal(t, "shared", def.Scope())
}
