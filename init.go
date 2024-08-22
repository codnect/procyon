package procyon

import (
	"codnect.io/procyon-core"
	"codnect.io/procyon-core/module"
)

func init() {
	module.Use[core.Module]()
}
