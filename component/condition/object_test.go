package condition

import (
	"codnect.io/procyon/component/container"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnObjectCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithNameExists(t *testing.T) {
	onObjectCondition := OnObject("anyObjectName")
	container := container.New()
	err := container.Singletons().Register("anyObjectName", "anyObject")
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onObjectCondition.MatchesCondition(conditionContext))
}

func TestOnObjectCondition_MatchesConditionShouldReturnTrueIfAnyDefinitionWithNameExists(t *testing.T) {
	onObjectCondition := OnObject("anyObjectName")
	objectContainer := container.New()

	definition, err := container.MakeDefinition(anyConstructorFunction, container.Named("anyObjectName"))
	assert.Nil(t, err)
	err = objectContainer.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onObjectCondition.MatchesCondition(conditionContext))
}

func TestOnObjectCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithNameDoesNotExist(t *testing.T) {
	onObjectCondition := OnObject("anyObjectName")
	objectContainer := container.New()

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onObjectCondition.MatchesCondition(conditionContext))
}

func TestOnMissingObjectCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithNameExists(t *testing.T) {
	onMissingObjectCondition := OnMissingObject("anyObjectName")
	objectContainer := container.New()
	err := objectContainer.Singletons().Register("anyObjectName", "anyObject")
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onMissingObjectCondition.MatchesCondition(conditionContext))
}

func TestOnMissingObjectCondition_MatchesConditionShouldReturnFalseIfAnyDefinitionWithNameExists(t *testing.T) {
	onMissingObjectCondition := OnMissingObject("anyObjectName")
	objectContainer := container.New()

	definition, err := container.MakeDefinition(anyConstructorFunction, container.Named("anyObjectName"))
	assert.Nil(t, err)
	err = objectContainer.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onMissingObjectCondition.MatchesCondition(conditionContext))
}

func TestOnMissingObjectCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithNameDoesNotExist(t *testing.T) {
	onMissingObjectCondition := OnMissingObject("anyObjectName")
	objectContainer := container.New()

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onMissingObjectCondition.MatchesCondition(conditionContext))
}
