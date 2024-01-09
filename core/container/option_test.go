package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptional_PanicsIfDependencyTypeIsNotPresentInConstructorParameterList(t *testing.T) {
	assert.Panics(t, func() {
		Register(AnyConstructFunction, Optional[string]())
	})
}

func TestOptionalAt_PanicsIfIndexIsLessThanZero(t *testing.T) {
	assert.Panics(t, func() {
		Register(AnyConstructFunction, OptionalAt(-1))
	})
}

func TestOptionalAt_PanicsIfIndexIsNotBetweenParameterIndexRange(t *testing.T) {
	assert.Panics(t, func() {
		Register(AnyConstructFunction, OptionalAt(3))
	})
}

func TestQualifier_PanicsIfDependencyTypeIsNotPresentInConstructorParameterList(t *testing.T) {
	assert.Panics(t, func() {
		Register(AnyConstructFunction, Qualifier[string]("anyName"))
	})
}

func TestQualifierAt_PanicsIfIndexIsLessThanZero(t *testing.T) {
	assert.Panics(t, func() {
		Register(AnyConstructFunction, QualifierAt(-1, "anyName"))
	})
}

func TestQualifierAt_PanicsIfIndexIsNotBetweenParameterIndexRange(t *testing.T) {
	assert.Panics(t, func() {
		Register(AnyConstructFunction, QualifierAt(3, "anyName"))
	})
}
