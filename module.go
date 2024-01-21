package procyon

import (
	"codnect.io/procyon-core/module"
	"codnect.io/reflector"
)

func Use[M module.Module]() {
	typ := reflector.TypeOf[M]()
	if reflector.IsStruct(typ) {
		moduleStruct := reflector.ToStruct(typ)
		instance, err := moduleStruct.Instantiate()

		if err != nil {
			panic(err)
		}

		m := instance.Elem().(module.Module)
		m.InitModule()
	}
}
