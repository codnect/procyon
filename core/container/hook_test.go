package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHooks_AddReturnsErrorIfHookIsNil(t *testing.T) {
	hooks := NewHooks()
	err := hooks.Add(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "container: hook cannot be nil", err.Error())
}

func TestHooks_AddReturnsErrorIfHookIsDuplicated(t *testing.T) {
	hooks := NewHooks()
	hook := PostInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	})

	err := hooks.Add(hook)
	assert.Nil(t, err)
	assert.Equal(t, hooks.Count(), 1)

	err = hooks.Add(hook)
	assert.NotNil(t, err)
	assert.Equal(t, "container: hook already exists", err.Error())
}

func TestHooks_AddRegistersHookSuccessfully(t *testing.T) {
	hooks := NewHooks()
	hook := PostInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	})

	err := hooks.Add(hook)
	assert.Nil(t, err)
	assert.Equal(t, hooks.Count(), 1)
}

func TestHooks_CountReturnsNumberOfHooks(t *testing.T) {
	hooks := NewHooks()
	err := hooks.Add(PostInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	}))
	assert.Nil(t, err)

	err = hooks.Add(PostInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	}))
	assert.Nil(t, err)
	assert.Equal(t, hooks.Count(), 2)
}

func TestHooks_RemoveDeletesAlreadyRegisteredHook(t *testing.T) {
	hooks := NewHooks()
	hook := PreInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	})

	err := hooks.Add(hook)
	assert.Nil(t, err)
	assert.Equal(t, hooks.Count(), 1)

	hooks.Remove(hook)
	assert.Equal(t, hooks.Count(), 0)
}

func TestHooks_ToSlice(t *testing.T) {
	hooks := NewHooks()
	hook := PostInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	})

	err := hooks.Add(hook)
	assert.Nil(t, err)

	items := hooks.ToSlice()
	assert.Len(t, items, 1)
	assert.Equal(t, []*Hook{hook}, items)
}

func TestHooks_RemoveAllDeletesAllRegisteredHooks(t *testing.T) {
	hooks := NewHooks()
	hook := PostInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	})

	err := hooks.Add(hook)
	assert.Nil(t, err)
	assert.Equal(t, hooks.Count(), 1)

	hooks.RemoveAll()
	assert.Equal(t, hooks.Count(), 0)
}

func TestPreInitialization(t *testing.T) {
	hook := PreInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	})
	assert.NotNil(t, hook)
	assert.NotNil(t, hook.OnPreInitialization)
	assert.Nil(t, hook.OnPostInitialization)
}

func TestPostInitialization(t *testing.T) {
	hook := PostInitialization(func(name string, instance any) (any, error) {
		return nil, nil
	})
	assert.NotNil(t, hook)
	assert.Nil(t, hook.OnPreInitialization)
	assert.NotNil(t, hook.OnPostInitialization)
}
