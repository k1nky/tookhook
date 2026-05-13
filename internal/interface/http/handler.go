package http

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k1nky/tookhook/internal/domain/entity"
)

// HandleWebhook handles incoming webhook requests.
func (s *Server) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "endpoint name is required", http.StatusBadRequest)
		return
	}

	// Get endpoint configuration
	endpoint, err := s.repo.GetByName(r.Context(), name)
	if err != nil {
		if err == entity.ErrEndpointNotFound {
			http.Error(w, "endpoint not found", http.StatusNotFound)
			return
		}
		s.logger.Error("failed to get endpoint", "error", err, "endpoint", name)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Check if endpoint is disabled
	if endpoint.Disabled {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("failed to read request body", "error", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Create webhook task
	task := &entity.WebhookTask{
		EndpointName: name,
		Payload:      body,
		ContentType:  r.Header.Get("Content-Type"),
		Headers:      r.Header,
		ID:           middleware.GetReqID(r.Context()),
	}

	// Enqueue task
	if err := s.queue.Enqueue(r.Context(), task); err != nil {
		s.logger.Error("failed to enqueue task", "error", err, "endpoint", name)
		http.Error(w, "failed to enqueue request", http.StatusInternalServerError)
		return
	}

	s.logger.Info("webhook enqueued", "endpoint", name, "size", len(body))
	w.WriteHeader(http.StatusCreated)
}

// HandleHealth handles health check requests.
func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// HandleReload handles configuration reload requests.
func (s *Server) HandleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.logger.Info("reloading configuration")
	if err := s.repo.Reload(r.Context()); err != nil {
		s.logger.Error("failed to reload configuration", "error", err)
		http.Error(w, "failed to reload configuration", http.StatusInternalServerError)
		return
	}

	s.logger.Info("configuration reloaded successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
