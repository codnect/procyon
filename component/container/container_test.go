package container

/*
import (
	"codnect.io/procyon/component/container"
	"codnect.io/procyon/component/filter"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestObjectContainer_GetObjectShouldReturnErrorIfItIsInvokedWithoutFilters(t *testing.T) {
	ctx := context.Background()
	container := NewContainer()
	got, err := container.GetObject(ctx)
	assert.Equal(t, "at least one filter must be used", err.Error())
	assert.Nil(t, got)
}

func TestObjectContainer_GetObjectShouldReturnObjectIfObjectExists(t *testing.T) {
	ctx := context.Background()

	container := NewContainer()
	mockSingletonRegistry := &MockSingletonRegistry{}
	container.singletons = mockSingletonRegistry

	anyObject := &AnyType{}
	mockSingletonRegistry.On("Find", mock.AnythingOfType("[]filter.Filter")).
		Return(anyObject, nil)

	got, err := container.GetObject(ctx, filter.ByName("anyName"))
	assert.Nil(t, err)
	assert.Equal(t, anyObject, got)
}

func TestObjectContainer_GetObjectShouldReturnErrorIfSingletonRegistryReturnErrorOtherThanNotFound(t *testing.T) {
	ctx := context.Background()

	container := NewContainer()
	mockSingletonRegistry := &MockSingletonRegistry{}
	container.singletons = mockSingletonRegistry

	mockSingletonRegistry.On("Find", mock.AnythingOfType("[]filter.Filter")).
		Return(nil, errors.New("anyError"))

	got, err := container.GetObject(ctx, filter.ByName("anyName"))
	assert.Equal(t, "anyError", err.Error())
	assert.Nil(t, got)
}

func TestObjectContainer_GetObjectShouldReturnSingletonObjectIfSingletonDefinitionExists(t *testing.T) {
	ctx := context.Background()

	container := NewContainer()
	mockSingletonRegistry := &MockSingletonRegistry{}
	container.singletons = mockSingletonRegistry

	mockDefinitionRegistry := &MockDefinitionRegistry{}
	container.definitions = mockDefinitionRegistry

	anyObject := &AnyType{}
	anyDefinition, _ := MakeDefinition(anyConstructorFunction)

	mockSingletonRegistry.On("Find", mock.AnythingOfType("[]filter.Filter")).
		Return(nil, &ObjectNotFoundError{})
	mockSingletonRegistry.On("OrElseCreate", anyDefinition.Name(), mock.AnythingOfType("component.ObjectProvider")).
		Return(anyObject, nil)
	mockDefinitionRegistry.On("Find", mock.AnythingOfType("[]filter.Filter")).
		Return(anyDefinition, nil)

	got, err := container.GetObject(ctx, filter.ByName("anyName"))
	assert.Nil(t, err)
	assert.Equal(t, anyObject, got)
}

func TestObjectContainer_GetObjectShouldReturnPrototypeObjectIfPrototypeDefinitionExists(t *testing.T) {
	ctx := context.Background()

	container := NewContainer()
	mockSingletonRegistry := &MockSingletonRegistry{}
	container.singletons = mockSingletonRegistry

	mockDefinitionRegistry := &MockDefinitionRegistry{}
	container.definitions = mockDefinitionRegistry

	anyObject := &AnyType{}
	anyDefinition, _ := MakeDefinition(anyConstructorFunction, Scoped(container2.PrototypeScope))

	mockSingletonRegistry.On("Find", mock.AnythingOfType("[]filter.Filter")).
		Return(nil, &ObjectNotFoundError{})
	mockDefinitionRegistry.On("Find", mock.AnythingOfType("[]filter.Filter")).
		Return(anyDefinition, nil)

	got, err := container.GetObject(ctx, filter.ByName("anyName"))
	assert.Nil(t, err)
	assert.Equal(t, anyObject, got)
}
*/
