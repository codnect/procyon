package procyon

import (
	"errors"
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
)

func getInstances(typ goo.Type) (result []interface{}, err error) {
	var types []goo.Type
	types, err = core.GetComponentTypes(typ)
	if err != nil {
		return
	}
	for _, t := range types {
		var instance interface{}
		instance, err = peas.CreateInstance(t, []interface{}{})
		if err != nil {
			return
		}
		if instance != nil {
			result = append(result, instance)
		} else {
			err = errors.New("Instance cannot be created by using the method " + t.GetName())
		}
	}
	return
}

func getInstancesWithParamTypes(typ goo.Type, parameterTypes []goo.Type, args []interface{}) (result []interface{}, err error) {
	var types []goo.Type
	types, err = core.GetComponentTypesWithParam(typ, parameterTypes)
	if err != nil {
		return
	}
	var instances []interface{}
	for _, t := range types {
		var instance interface{}
		instance, err = peas.CreateInstance(t, args)
		if err != nil {
			return
		}
		instances = append(instances, instance)
	}
	return instances, nil
}
