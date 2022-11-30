package container

import "sync"

type PreInitializeHook interface {
	OnPreInitialization(name string, instance any) (any, error)
}

type PostInitializeHook interface {
	OnPostInitialization(name string, instance any) (any, error)
}

type InitializeHook interface {
	OnInitialization() error
}

type Hook interface {
}

type Hooks struct {
	hooks map[string]Hook
	mu    sync.RWMutex
}

func newHooks() *Hooks {
	return &Hooks{
		make(map[string]Hook),
		sync.RWMutex{},
	}
}

func (h *Hooks) Add(hook Hook) error {
	return nil
}

func (h *Hooks) Remove(hook Hook) {

}

func (h *Hooks) Hooks() []Hook {
	return nil
}

func (h *Hooks) Count() int {
	return 0
}

func (h *Hooks) RemoveAll() {

}
