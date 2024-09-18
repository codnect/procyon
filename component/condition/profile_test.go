package condition

import (
	"codnect.io/procyon/component/container"
	"codnect.io/procyon/runtime"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnProfileCondition_MatchesConditionShouldReturnTrueIfProfileIsActivated(t *testing.T) {
	onProfileCondition := OnProfile("anyProfileName")
	objectContainer := container.New()

	environment := runtime.NewDefaultEnvironment()
	environment.SetActiveProfiles("anyProfileName")
	objectContainer.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.True(t, onProfileCondition.MatchesCondition(conditionContext))
}

func TestOnProfileCondition_MatchesConditionShouldReturnFalseIfEnvironmentObjectDoesNotExist(t *testing.T) {
	onProfileCondition := OnProfile("anyProfileName")
	objectContainer := container.New()

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onProfileCondition.MatchesCondition(conditionContext))
}

func TestOnProfileCondition_MatchesConditionShouldReturnFalseIfProfileIsNotActivated(t *testing.T) {
	onProfileCondition := OnProfile("anyProfileName")
	objectContainer := container.New()

	environment := runtime.NewDefaultEnvironment()
	environment.SetActiveProfiles("anotherProfileName")
	objectContainer.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), objectContainer)
	assert.False(t, onProfileCondition.MatchesCondition(conditionContext))
}
