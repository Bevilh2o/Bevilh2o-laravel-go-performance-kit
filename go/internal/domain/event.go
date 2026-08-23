package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrTenantRequired = errors.New("tenant identifier is required")
	ErrEventRequired  = errors.New("event type is required")
)

// IngestEventRequest represents the incoming JSON contract matching the Laravel API schema.
type IngestEventRequest struct {
	Tenant    string          `json:"tenant"`
	Event     string          `json:"event"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp *int64          `json:"timestamp,omitempty"`
}

// Validate ensures all required contract invariants are satisfied.
func (r *IngestEventRequest) Validate() error {
	if r.Tenant == "" {
		return ErrTenantRequired
	}
	if r.Event == "" {
		return ErrEventRequired
	}
	return nil
}

// EventEntity represents the persistent record matching PostgreSQL schema.
type EventEntity struct {
	ID        int64
	TenantID  string
	EventType string
	Payload   []byte
	OccurredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}