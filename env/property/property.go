package property

type Source interface {
	Name() string
	Source() any
	ContainsProperty(name string) bool
	Property(name string) (any, bool)
	PropertyOrDefault(name string, defaultValue any) any
	PropertyNames() []string
}

type Sources interface {
	Contains(name string) bool
	Find(name string) (Source, bool)
	Size() int
	AddFirst(source Source)
	AddLast(source Source)
	AddAtIndex(index int, source Source)
	Remove(name string)
	Replace(name string, source Source)
	PrecendenceOf(source Source) int
}
