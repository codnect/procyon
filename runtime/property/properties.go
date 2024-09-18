package property

// Properties is a marker interface used to mark structs as property structs.
type Properties interface {
	noPropertiesMethodYet()
}
