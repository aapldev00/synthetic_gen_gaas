package generator

import (
	"sync"
)

var (
	// mu protects the providers map from concurrent write access.
	mu        sync.RWMutex
	providers = make(map[string]GeneratorFactory)
)

// Register adds a new data provider factory to the central registry.
// This function is typically called during the application's initialization
// phase by the individual provider packages.
//
// If a provider with the same name is already registered, it will be overwritten.
func Register(name string, factory GeneratorFactory) {
	mu.Lock()
	defer mu.Unlock()
	providers[name] = factory
}

// GetFactory retrieves a GeneratorFactory by its provider name.
// It returns the factory and a boolean indicating whether the provider was found.
//
// This operation is thread-safe and optimized for concurrent reads.
func GetFactory(name string) (GeneratorFactory, bool) {
	mu.RLock()
	defer mu.RUnlock()
	factory, ok := providers[name]
	return factory, ok
}

// ListProviders returns a list of all currently registered provider names.
// Useful for debugging, logging, or discovery services.
func ListProviders() []string {
	mu.RLock()
	defer mu.RUnlock()

	list := make([]string, 0, len(providers))
	for name := range providers {
		list = append(list, name)
	}
	return list
}
