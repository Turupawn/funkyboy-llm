package config

import "os"

// TierConfig holds the provider settings for one tier.
type TierConfig struct {
	Provider string
	APIKey   string
	Model    string
	OllamaURL string // only used when Provider == "ollama"
}

// Config holds the full runtime configuration.
type Config struct {
	Port string
	XS   TierConfig
	Fast TierConfig
	SOTA TierConfig
}

func Load() *Config {
	return &Config{
		Port: envOr("PORT", "8090"),
		XS: TierConfig{
			Provider:  envOr("XS_PROVIDER", "ollama"),
			OllamaURL: envOr("XS_OLLAMA_URL", "http://ollama:11434"),
			Model:     envOr("XS_MODEL", "dolphin-unleashed"),
		},
		Fast: TierConfig{
			Provider: envOr("FAST_PROVIDER", "gemini"),
			APIKey:   os.Getenv("FAST_API_KEY"),
			Model:    envOr("FAST_MODEL", "gemini-2.0-flash"),
		},
		SOTA: TierConfig{
			Provider: envOr("SOTA_PROVIDER", "anthropic"),
			APIKey:   os.Getenv("SOTA_API_KEY"),
			Model:    envOr("SOTA_MODEL", "claude-opus-4-6"),
		},
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
