package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type HealthHandler struct {
	client *mongo.Client
}

func NewHealthHandler(client *mongo.Client) *HealthHandler {
	return &HealthHandler{client: client}
}

func (h *HealthHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.liveness)
	mux.HandleFunc("/readyz", h.readiness)
}

func (h *HealthHandler) liveness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"service": "relay",
		"status":  "ok",
	})
}

func (h *HealthHandler) readiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.client.Ping(ctx, readpref.Primary()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not_ready",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":   "ready",
		"database": "connected",
	})
}

func writeJSON(w http.ResponseWriter, status int, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
