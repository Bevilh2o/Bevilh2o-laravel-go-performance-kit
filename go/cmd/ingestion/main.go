package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laravel-go-performance-kit/ingestion/internal/config"
	"github.com/laravel-go-performance-kit/ingestion/internal/handler"
	"github.com/laravel-go-performance-kit/ingestion/internal/middleware"
	"github.com/laravel-go-performance-kit/ingestion/internal/repository"
)

func main() {
	// Initialize JSON Structured Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.LoadFromEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("initializing_database_pool", slog.String("host", "postgres"))
	repo, err := repository.NewEventRepository(ctx, cfg.DBDSN)
	if err != nil {
		logger.Error("database_connection_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer repo.Close()
	logger.Info("database_pool_ready")

	eventHandler := handler.NewEventHandler(repo)

	// Application Multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/events", eventHandler.Ingest)

	// Observability: Register standard pprof profiling endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Compose Middleware Chain: RequestID -> Logger -> Recovery -> Router
	var rootHandler http.Handler = mux
	rootHandler = middleware.Recovery(logger)(rootHandler)
	rootHandler = middleware.StructuredLogger(logger)(rootHandler)
	rootHandler = middleware.RequestID(rootHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      rootHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background goroutine
	go func() {
		logger.Info("server_started", slog.String("port", cfg.HTTPPort))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server_fatal_error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Await OS Interrupt Signals
	<-ctx.Done()
	logger.Info("shutdown_signal_received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful_shutdown_failed", slog.String("error", err.Error()))
	}

	logger.Info("server_terminated_cleanly")
}