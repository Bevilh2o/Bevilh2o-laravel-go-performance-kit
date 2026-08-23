package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/laravel-go-performance-kit/ingestion/internal/domain"
)

type EventRepository struct {
	pool *pgxpool.Pool
}

// NewEventRepository initializes a connection pool to PostgreSQL.
func NewEventRepository(ctx context.Context, dsn string) (*EventRepository, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	// Performance tuning: configure pool limits for high concurrency
	config.MaxConns = 50
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	// Ping database to verify connection readiness
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &EventRepository{pool: pool}, nil
}

// Close gracefully closes the connection pool.
func (r *EventRepository) Close() {
	r.pool.Close()
}

// InsertEvent persists an event entity into PostgreSQL matching Laravel's schema.
func (r *EventRepository) InsertEvent(ctx context.Context, req *domain.IngestEventRequest) (int64, time.Time, error) {
	var occurredAt time.Time
	if req.Timestamp != nil {
		occurredAt = time.Unix(*req.Timestamp, 0).UTC()
	} else {
		occurredAt = time.Now().UTC()
	}

	now := time.Now().UTC()

	var payloadData []byte
	if len(req.Payload) > 0 {
		payloadData = req.Payload
	}

	query := `
		INSERT INTO events (tenant_id, event_type, payload, occurred_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := r.pool.QueryRow(ctx, query, req.Tenant, req.Event, payloadData, occurredAt, now, now).Scan(&id)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("insert event: %w", err)
	}

	return id, occurredAt, nil
}