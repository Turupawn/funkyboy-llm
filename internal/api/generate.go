package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/funkyboy/llm/internal/tier"
)

type generateRequest struct {
	Prompt string `json:"prompt"`
	Tier   string `json:"tier"`
}

type generateResponse struct {
	Response  string `json:"response"`
	Tier      string `json:"tier"`
	Model     string `json:"model"`
	Provider  string `json:"provider"`
	LatencyMs int64  `json:"latency_ms"`
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Prompt == "" {
		writeError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	if req.Tier == "" {
		req.Tier = tier.XS
	}

	if !tier.ValidTier(req.Tier) {
		writeError(w, http.StatusBadRequest, "invalid tier: "+req.Tier+". Must be xs, fast, or sota")
		return
	}

	p, err := s.registry.Get(req.Tier)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.log.Info("generate", "tier", req.Tier, "provider", p.Name(), "model", p.Model(), "prompt_len", len(req.Prompt))

	start := time.Now()
	response, err := p.Generate(r.Context(), req.Prompt)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		s.log.Error("generate failed", "tier", req.Tier, "error", err, "latency_ms", latency)
		writeError(w, http.StatusBadGateway, "backend error: "+err.Error())
		return
	}

	s.log.Info("generate ok", "tier", req.Tier, "latency_ms", latency, "response_len", len(response))

	writeJSON(w, http.StatusOK, generateResponse{
		Response:  response,
		Tier:      req.Tier,
		Model:     p.Model(),
		Provider:  p.Name(),
		LatencyMs: latency,
	})
}