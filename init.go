package procyon

import (
	"codnect.io/procyon-core/component"
)

func init() {
	component.Register(newConfigContextConfigurer, component.Named("procyonConfigContextConfigurer"))
}
