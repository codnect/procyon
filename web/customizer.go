package web

import "codnect.io/procyon/app"

type ContextCustomizer struct {
}

func NewContextCustomizer() *ContextCustomizer {
	return &ContextCustomizer{}
}

func (c *ContextCustomizer) CustomizeContext(ctx app.Context) error {
	return nil
}
