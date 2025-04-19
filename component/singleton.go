package component

// SingletonRegistry defines methods for managing singleton instances within the component system.
type SingletonRegistry interface {
	// RegisterSingleton registers a singleton instance with the given name.
	// Returns an error if a singleton with the same name already exists.
	RegisterSingleton(name string, instance any) error

	// ContainsSingleton checks whether a singleton with the specified name exists.
	ContainsSingleton(name string) bool

	// Singleton retrieves the singleton instance associated with the given name.
	// Returns the instance and a boolean indicating its existence.
	Singleton(name string) (any, bool)

	// RemoveSingleton removes the singleton instance associated with the specified name.
	RemoveSingleton(name string)
}
