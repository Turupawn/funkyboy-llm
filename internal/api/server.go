package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/funkyboy/llm/internal/tier"
)

type Server struct {
	registry *tier.Registry
	log      *slog.Logger
	Mux      *http.ServeMux
}

func NewServer(registry *tier.Registry, log *slog.Logger) *Server {
	s := &Server{
		registry: registry,
		log:      log,
		Mux:      http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Mux.HandleFunc("POST /v1/generate", s.handleGenerate)
	s.Mux.HandleFunc("GET /health", s.handleHealth)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
