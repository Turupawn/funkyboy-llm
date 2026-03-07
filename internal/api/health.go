package api

import (
	"context"
	"net/http"
	"time"
)

type tierHealth struct {
	Status   string `json:"status"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Error    string `json:"error,omitempty"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	result := make(map[string]tierHealth)

	for name, p := range s.registry.All() {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		err := p.Ping(ctx)
		cancel()

		h := tierHealth{
			Provider: p.Name(),
			Model:    p.Model(),
		}
		if err != nil {
			h.Status = "error"
			h.Error = err.Error()
		} else {
			h.Status = "ok"
		}
		result[name] = h
	}

	writeJSON(w, http.StatusOK, result)
}
