package container

import (
	"context"
	"fmt"
	"sync"
)

type ctxHolder struct{}

var ctxHolderKey = &ctxHolder{}

type holder struct {
	currentlyInCreation map[string]struct{}
	mu                  sync.RWMutex
}

func newHolder() *holder {
	return &holder{}
}

func (h *holder) isCurrentlyInCreation(name string) bool {
	defer h.mu.Unlock()
	h.mu.Lock()
	if _, ok := h.currentlyInCreation[name]; ok {
		return true
	}

	return false
}

func (h *holder) beforeCreation(name string) error {
	defer h.mu.Unlock()
	h.mu.Lock()

	if _, ok := h.currentlyInCreation[name]; ok {
		return fmt.Errorf("container: instance with name %s is currently in creation, maybe it has got circular dependency cycle", name)
	}

	h.currentlyInCreation[name] = struct{}{}
	return nil
}

func (h *holder) afterCreation(name string) {
	defer h.mu.Unlock()
	h.mu.Lock()
	delete(h.currentlyInCreation, name)
}

func holderFromContext(ctx context.Context) *holder {
	return ctx.Value(ctxHolderKey).(*holder)
}

func contextWithHolder(parent context.Context) context.Context {
	h := parent.Value(ctxHolderKey)

	if h != nil {
		return parent
	}

	return context.WithValue(parent, ctxHolderKey, newHolder())
}
