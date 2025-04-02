package runtime

// Module is an interface that represents a module.
// It provides a method to initialize the module.
type Module interface {
	// InitModule method initializes the module.
	// It returns an error if the initialization fails.
	InitModule() error
}
