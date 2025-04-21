package component

// Container provides a unified interface for managing components,
// including definition registration, instance resolving, custom scopes,
// lifecycle management, and manual bindings.
type Container interface {
	// DefinitionRegistry provides access to component definitions and their metadata.
	DefinitionRegistry

	// SingletonRegistry manages singleton instances of components.
	SingletonRegistry

	// Resolver resolves component instances by type or name.
	Resolver

	// Binder allows manual binding of instances to specific types.
	Binder

	// ScopeRegistry manages custom scopes.
	ScopeRegistry

	// LifecycleManager manages lifecycle hooks.
	LifecycleManager
}
