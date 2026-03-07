package provider

import "context"

// Provider is the interface all LLM backends implement.
type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
	Name() string
	Model() string
	Ping(ctx context.Context) error
}
