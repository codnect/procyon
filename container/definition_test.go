package container

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"runtime"
	"testing"
)

type AnyType struct {
	*mock.Mock
	t *DependencyType
}

func (a *AnyType) String() string {
	return ""
}

func (a *AnyType) PostConstruct() error {
	a.Mock.Called()
	return nil
}

type DependencyType struct {
}

func AnyConstructFunction(t *DependencyType) *AnyType {
	m := &mock.Mock{}
	m.On("PostConstruct").Return()

	return &AnyType{
		m,
		t,
	}
}

func AnyConstructFunctionReturningNonPointerValue(t *DependencyType) AnyType {
	return AnyType{}
}

func TestMakeDefinition_WithConstructorReturningPointerType(t *testing.T) {
	def, err := MakeDefinition(AnyConstructFunction, Qualifier[*DependencyType]("anyDependencyType"),
		Scope(PrototypeScope),
		QualifierAt(0, "anotherDependencyType"),
		OptionalAt(0))
	assert.Nil(t, err)

	assert.Equal(t, "anyType", def.Name())
	assert.Equal(t, "*AnyType", def.Type().Name())
	assert.Equal(t, "*AnyType", def.reflectorType().Name())
	assert.Equal(t, runtime.FuncForPC(reflect.ValueOf(AnyConstructFunction).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(def.Constructor()).Pointer()).Name())

	inputs := def.Inputs()
	assert.NotNil(t, inputs)
	assert.Len(t, inputs, 1)

	assert.Equal(t, 0, inputs[0].Index())
	assert.True(t, inputs[0].IsOptional())
	assert.Equal(t, "anotherDependencyType", inputs[0].Name())
	assert.Equal(t, "*DependencyType", inputs[0].Type().Name())

	assert.Equal(t, true, def.IsPrototype())
	assert.Equal(t, false, def.IsShared())
	assert.Equal(t, "prototype", def.Scope())
}

func TestMakeDefinition_WithConstructorReturningNonPointerType(t *testing.T) {
	def, err := MakeDefinition(AnyConstructFunctionReturningNonPointerValue, Qualifier[*DependencyType]("anyDependencyType"))
	assert.Nil(t, err)

	assert.Equal(t, "anyType", def.Name())
	assert.Equal(t, "AnyType", def.Type().Name())
	assert.Equal(t, "AnyType", def.reflectorType().Name())
	assert.Equal(t, runtime.FuncForPC(reflect.ValueOf(AnyConstructFunctionReturningNonPointerValue).Pointer()).Name(),
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

func TestMakeDefinition_WithOptionError(t *testing.T) {
	def, err := MakeDefinition(AnyConstructFunction, QualifierAt(2, ""))
	assert.NotNil(t, err)
	assert.Equal(t, "container: could not find any input at index 2", err.Error())
	assert.Nil(t, def)
}

func TestMakeDefinition_WithoutConstructor(t *testing.T) {
	def, err := MakeDefinition(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "container: constructor should not be nil", err.Error())
	assert.Nil(t, def)
}

func TestMakeDefinition_WithNonFunction(t *testing.T) {
	def, err := MakeDefinition(struct{}{})
	assert.NotNil(t, err)
	assert.Equal(t, "container: constructor should be a function", err.Error())
	assert.Nil(t, def)
}

func TestMakeDefinition_WithFunctionReturningMultipleValues(t *testing.T) {
	def, err := MakeDefinition(func() (string, error) {
		return "", nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, "container: constructor can only be a function returning one result", err.Error())
	assert.Nil(t, def)
}
