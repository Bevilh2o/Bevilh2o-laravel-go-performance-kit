package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laravel-go-performance-kit/ingestion/internal/config"
	"github.com/laravel-go-performance-kit/ingestion/internal/handler"
	"github.com/laravel-go-performance-kit/ingestion/internal/repository"
)

func main() {
	cfg := config.LoadFromEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("[INFO] Initializing PostgreSQL connection pool...")
	repo, err := repository.NewEventRepository(ctx, cfg.DBDSN)
	if err != nil {
		log.Fatalf("[FATAL] Could not connect to database: %v", err)
	}
	defer repo.Close()
	log.Println("[INFO] Database connection pool established successfully.")

	eventHandler := handler.NewEventHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/events", eventHandler.Ingest)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in a separate goroutine
	go func() {
		log.Printf("[INFO] Go Ingestion Service listening on port %s...", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] HTTP server error: %v", err)
		}
	}()

	// Block until an interrupt signal is received
	<-ctx.Done()
	log.Println("[INFO] Shutdown signal received. Commencing graceful termination...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] Graceful shutdown encountered an error: %v", err)
	}

	log.Println("[INFO] Go Ingestion Service terminated cleanly.")
}