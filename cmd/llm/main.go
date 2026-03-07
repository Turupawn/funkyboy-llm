package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/funkyboy/llm/internal/api"
	"github.com/funkyboy/llm/internal/config"
	"github.com/funkyboy/llm/internal/provider"
	"github.com/funkyboy/llm/internal/tier"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("FUNKYBOY-LLM starting")

	cfg := config.Load()

	providers := make(map[string]provider.Provider)

	// Build each tier's provider
	if p, err := buildProvider(cfg.XS); err != nil {
		log.Warn("xs tier unavailable", "error", err)
	} else {
		providers[tier.XS] = p
		log.Info("tier ready", "tier", "xs", "provider", p.Name(), "model", p.Model())
	}

	if p, err := buildProvider(cfg.Fast); err != nil {
		log.Warn("fast tier unavailable", "error", err)
	} else {
		providers[tier.Fast] = p
		log.Info("tier ready", "tier", "fast", "provider", p.Name(), "model", p.Model())
	}

	if p, err := buildProvider(cfg.SOTA); err != nil {
		log.Warn("sota tier unavailable", "error", err)
	} else {
		providers[tier.SOTA] = p
		log.Info("tier ready", "tier", "sota", "provider", p.Name(), "model", p.Model())
	}

	if len(providers) == 0 {
		log.Error("no tiers configured — at least one provider must be available")
		os.Exit(1)
	}

	registry := tier.New(providers)
	srv := api.NewServer(registry, log)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: srv.Mux,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Info("shutdown signal received", "signal", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		httpServer.Shutdown(shutdownCtx)
	}()

	log.Info("FUNKYBOY-LLM online", "port", cfg.Port, "tiers", len(providers))
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("server error", "err", err)
		os.Exit(1)
	}

	log.Info("FUNKYBOY-LLM shutdown complete")
}

func buildProvider(tc config.TierConfig) (provider.Provider, error) {
	switch tc.Provider {
	case "ollama":
		return provider.NewOllama(tc.OllamaURL, tc.Model), nil
	case "gemini":
		if tc.APIKey == "" {
			return nil, fmt.Errorf("API key required for gemini")
		}
		return provider.NewGemini(tc.APIKey, tc.Model), nil
	case "anthropic":
		if tc.APIKey == "" {
			return nil, fmt.Errorf("API key required for anthropic")
		}
		return provider.NewAnthropic(tc.APIKey, tc.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", tc.Provider)
	}
}
