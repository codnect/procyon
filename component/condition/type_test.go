package condition

import (
	"codnect.io/procyon/component/container"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AnyType struct {
}

func anyConstructorFunction() AnyType {
	return AnyType{}
}

func TestOnTypeCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithTypeExists(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	objectContainer := container.New()
	err := objectContainer.Singletons().Register("anyObject", AnyType{})
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnTypeCondition_MatchesConditionShouldReturnTrueIfAnyDefinitionWithTypeExists(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	objectContainer := container.New()

	definition, err := container.MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = objectContainer.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnTypeCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithTypeDoesNotExist(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	objectContainer := container.New()

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithTypeExists(t *testing.T) {
	onMissingTypeCondition := OnMissingType[AnyType]()
	objectContainer := container.New()
	err := objectContainer.Singletons().Register("anyObject", AnyType{})
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnFalseIfAnyDefinitionWithTypeExists(t *testing.T) {
	onMissingTypeCondition := OnMissingType[AnyType]()
	objectContainer := container.New()

	definition, err := container.MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = objectContainer.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithTypeDoesNotExist(t *testing.T) {
	onMissingTypeCondition := OnMissingObject("anyObjectName")
	objectContainer := container.New()

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}
