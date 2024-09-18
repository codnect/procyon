package component

type AnyInterface interface {
	AnyMethod()
}

type AnotherInterface interface {
	AnotherMethod()
}

type AnyType struct {
}

func anyConstructorFunction() *AnyType {
	return &AnyType{}
}

func (a *AnyType) AnyMethod() {

}

type AnotherType struct {
}

func anotherConstructorFunction() *AnotherType {
	return &AnotherType{}
}

func (a *AnotherType) AnyMethod() {

}

func (a *AnotherType) AnotherMethod() {

}
