package web

import "github.com/procyon-projects/procyon/app"

type ContextCustomizer struct {
}

func NewContextCustomizer() *ContextCustomizer {
	return &ContextCustomizer{}
}

func (c *ContextCustomizer) CustomizeContext(ctx app.Context) error {
	return nil
}
