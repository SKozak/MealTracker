package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SKozak/mealtracker/modules/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	//load env vars
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with system env vars")
	}

	// Initialize structured logger (slog is stdlib since Go 1.21)
	// Similar to SLF4J but built-in
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("configuration loaded",
		"env", cfg.App.Env,
		"addr", cfg.Server.Addr(),
	)

	// Setup router
	// Chi is similar to having a lightweight Spring MVC without the magic
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID) // Adds X-Request-ID header
	r.Use(middleware.RealIP)    // Sets RemoteAddr from X-Forwarded-For
	r.Use(middleware.Logger)    // Logs requests (we'll replace with slog later)
	r.Use(middleware.Recoverer) // Recovers from panics (like @ExceptionHandler)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check - always have this for k8s/docker
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Placeholder for future routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head><title>Meal Tracker</title></head>
			<body>
				<h1>🍽️ Meal Tracker</h1>
				<p>Coming soon...</p>
			</body>
			</html>
		`))
	})

	// Create HTTP server with timeouts
	// In Java you'd configure this in application.yml or Tomcat settings
	server := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown handling
	// This pattern is idiomatic Go - explicit signal handling
	go func() {
		slog.Info("server starting", "addr", cfg.Server.Addr())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")

}
