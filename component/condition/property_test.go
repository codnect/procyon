package condition

import (
	"codnect.io/procyon/component/container"
	"codnect.io/procyon/runtime"
	"codnect.io/procyon/runtime/property"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueIfPropertyExists(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName")
	objectContainer := container.New()

	anyPropertySource := property.NewMapSource("anyPropertySource", map[string]interface{}{
		"anyPropertyName": true,
	})

	environment := runtime.NewDefaultEnvironment()
	environment.PropertySources().AddLast(anyPropertySource)
	objectContainer.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueEvenIfPropertyDoesNotExistAndMatchIfMissingIsCalled(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName").MatchIfMissing(true)
	objectContainer := container.New()

	environment := runtime.NewDefaultEnvironment()
	objectContainer.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueIfPropertyValueEqualsToGivenValue(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName").HavingValue("anyPropertyValue")
	objectContainer := container.New()

	anyPropertySource := property.NewMapSource("anyPropertySource", map[string]interface{}{
		"anyPropertyName": "anyPropertyValue",
	})

	environment := runtime.NewDefaultEnvironment()
	environment.PropertySources().AddLast(anyPropertySource)
	objectContainer.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueIfPropertyValueDoesNotEqualToGivenValue(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName").HavingValue("anotherPropertyValue")
	objectContainer := container.New()

	anyPropertySource := property.NewMapSource("anyPropertySource", map[string]interface{}{
		"anyPropertyName": "anyPropertyValue",
	})

	environment := runtime.NewDefaultEnvironment()
	environment.PropertySources().AddLast(anyPropertySource)
	objectContainer.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnFalseIfEnvironmentObjectDoesNotExist(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName")
	objectContainer := container.New()

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onPropertyCondition.MatchesCondition(conditionContext))
}
