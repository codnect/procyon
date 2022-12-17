package app

import "github.com/procyon-projects/procyon/container"

func init() {
	container.Register(newStartupListener)
}
