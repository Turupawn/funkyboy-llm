package tier

import (
	"fmt"

	"github.com/funkyboy/llm/internal/provider"
)

const (
	XS   = "xs"
	Fast = "fast"
	SOTA = "sota"
)

// Registry maps tier names to providers.
type Registry struct {
	tiers map[string]provider.Provider
}

// New creates a registry from a map of tier→provider.
func New(providers map[string]provider.Provider) *Registry {
	return &Registry{tiers: providers}
}

// Get returns the provider for a tier.
func (r *Registry) Get(name string) (provider.Provider, error) {
	p, ok := r.tiers[name]
	if !ok {
		return nil, fmt.Errorf("unknown tier: %s", name)
	}
	return p, nil
}

// All returns all registered tier names and providers.
func (r *Registry) All() map[string]provider.Provider {
	return r.tiers
}

// ValidTier returns true if the tier name is known.
func ValidTier(name string) bool {
	return name == XS || name == Fast || name == SOTA
}
