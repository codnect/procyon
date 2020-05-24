package procyon

import (
	"log"
	"reflect"
	"runtime"
	"strings"
)

type Component interface{}

var (
	componentTypes = make(map[string]reflect.Type, 0)
)

func Register(components ...Component) {
	for _, component := range components {
		typ := getComponentType(component)
		if isSupportComponent(typ) {
			name := getComponentName(component)
			registerComponentType(name, typ)
		} else {
			log.Fatal("It supports only struct and function")
		}
	}
}

func registerComponentType(name string, typ reflect.Type) {
	if _, ok := componentTypes[name]; ok {
		log.Fatal("You have already registered the same component : " + name)
	}
	componentTypes[name] = typ
}

func isSupportComponent(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct || typ.Kind() == reflect.Func
}

func getComponentType(component Component) reflect.Type {
	typ := reflect.TypeOf(component)
	if typ == nil {
		log.Fatal("Type cannot be determined.")
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func getComponentName(component Component) string {
	typ := getComponentType(component)
	var name string
	if typ.Kind() == reflect.Struct {
		name := sanitizedName(typ.PkgPath())
		name = name + "$" + typ.Name()
	} else {
		funcFullName := getFullFunctionName(component)
		lastIndexForDot := strings.LastIndex(funcFullName, ".")
		funcFullName = funcFullName[0:lastIndexForDot] + "#" + funcFullName[lastIndexForDot+1:]
		name = sanitizedName(funcFullName)
	}
	return name
}

func sanitizedName(str string) string {
	name := strings.ReplaceAll(str, "/", ".")
	name = strings.ReplaceAll(name, "-", ".")
	name = strings.ReplaceAll(name, "_", ".")
	return name
}

func getFullFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
