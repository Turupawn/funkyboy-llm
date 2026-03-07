package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnthropicProvider calls the Anthropic Messages API.
type AnthropicProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewAnthropic(apiKey, model string) *AnthropicProvider {
	if model == "" {
		model = "claude-opus-4-6"
	}
	return &AnthropicProvider{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (a *AnthropicProvider) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody, err := json.Marshal(map[string]any{
		"model":      a.model,
		"max_tokens": 4096,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read anthropic response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("decode anthropic response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("anthropic returned empty content")
	}

	return result.Content[0].Text, nil
}

func (a *AnthropicProvider) Name() string  { return "anthropic" }
func (a *AnthropicProvider) Model() string { return a.model }

func (a *AnthropicProvider) Ping(ctx context.Context) error {
	_, err := a.Generate(ctx, "ping")
	return err
}
