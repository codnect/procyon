package procyon

import (
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testComponent1 struct {
}

func newTestComponent1() testComponent1 {
	return testComponent1{}
}

type testComponent2 struct {
}

func newTestComponent2(variable string) testComponent2 {
	return testComponent2{}
}

func init() {
	core.Register(newTestComponent1)
	core.Register(newTestComponent2)
}

func TestGetInstances(t *testing.T) {
	instances, err := getInstances(goo.GetType(testComponent1{}))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(instances))
}

func TestGetInstancesWithParamTypes(t *testing.T) {
	instances, err := getInstancesWithParamTypes(goo.GetType(testComponent2{}),
		[]goo.Type{goo.GetType((*string)(nil))},
		[]interface{}{"test"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(instances))

	instances, err = getInstancesWithParamTypes(goo.GetType(testComponent1{}),
		[]goo.Type{goo.GetType((*string)(nil))},
		[]interface{}{"test"})
	assert.Nil(t, err)
	assert.Nil(t, instances)
}
