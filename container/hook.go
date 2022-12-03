package container

import (
	"fmt"
	"sync"
)

type PreInitializationHook interface {
	OnPreInitialization(name string, instance any) (any, error)
}

type PostInitializationHook interface {
	OnPostInitialization(name string, instance any) (any, error)
}

type HookFunc interface {
	func(name string, instance any) (any, error)
}

func PostInitialization[F HookFunc](f F) *Hook {
	return &Hook{
		OnPostInitialization: f,
	}
}

func PreInitialization[F HookFunc](f F) *Hook {
	return &Hook{
		OnPreInitialization: f,
	}
}

type Hook struct {
	OnPreInitialization  func(string, any) (any, error)
	OnPostInitialization func(string, any) (any, error)
}

type Hooks struct {
	hooks map[*Hook]struct{}
	mu    sync.RWMutex
}

func NewHooks() *Hooks {
	return &Hooks{
		make(map[*Hook]struct{}),
		sync.RWMutex{},
	}
}

func (h *Hooks) Add(hook *Hook) error {
	defer h.mu.Unlock()
	h.mu.Lock()

	if hook == nil {
		return fmt.Errorf("container: hook cannot be nil")
	}

	if _, ok := h.hooks[hook]; ok {
		return fmt.Errorf("container: hook already exists")
	}

	h.hooks[hook] = struct{}{}
	return nil
}

func (h *Hooks) Remove(hook *Hook) {
	defer h.mu.Unlock()
	h.mu.Lock()

	if hook != nil {
		delete(h.hooks, hook)
	}
}

func (h *Hooks) ToSlice() []*Hook {
	defer h.mu.Unlock()
	h.mu.Lock()

	hooks := make([]*Hook, 0)

	for hook, _ := range h.hooks {
		hooks = append(hooks, hook)
	}

	return hooks
}

func (h *Hooks) Count() int {
	defer h.mu.Unlock()
	h.mu.Lock()
	return len(h.hooks)
}

func (h *Hooks) RemoveAll() {
	defer h.mu.Unlock()
	h.mu.Lock()

	for key := range h.hooks {
		delete(h.hooks, key)
	}
}
