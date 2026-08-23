package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/laravel-go-performance-kit/ingestion/internal/domain"
	"github.com/laravel-go-performance-kit/ingestion/internal/repository"
)

type EventHandler struct {
	repo *repository.EventRepository
}

func NewEventHandler(repo *repository.EventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

// Ingest handles POST /events with zero unnecessary memory allocations.
func (h *EventHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req domain.IngestEventRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid json payload",
		})
		return
	}

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Propagate context from the incoming HTTP request (handles client disconnects/cancellations)
	id, occurredAt, err := h.repo.InsertEvent(r.Context(), &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to persist event",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":          id,
			"tenant_id":   req.Tenant,
			"event_type":  req.Event,
			"occurred_at": occurredAt.Format(time.RFC3339),
		},
	})
}