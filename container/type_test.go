package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestType_Name(t *testing.T) {
	typ := TypeOf[AnyType]()
	assert.NotNil(t, typ)
	assert.Equal(t, "AnyType", typ.Name())

	typ = TypeOfAny(AnyType{})
	assert.NotNil(t, typ)
	assert.Equal(t, "AnyType", typ.Name())
}

func TestType_PackageName(t *testing.T) {
	typ := TypeOf[AnyType]()
	assert.NotNil(t, typ)
	assert.Equal(t, "container", typ.PackageName())

	typ = TypeOfAny(AnyType{})
	assert.NotNil(t, typ)
	assert.Equal(t, "container", typ.PackageName())
}

func TestType_ReflectType(t *testing.T) {
	typ := TypeOf[AnyType]()
	assert.NotNil(t, typ)
	assert.NotNil(t, typ.ReflectType())
	assert.Equal(t, "AnyType", typ.ReflectType().Name())

	typ = TypeOfAny(AnyType{})
	assert.NotNil(t, typ)
	assert.NotNil(t, typ.ReflectType())
	assert.Equal(t, "AnyType", typ.ReflectType().Name())
}
