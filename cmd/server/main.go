package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"

	"github.com/yourusername/go-app/internal/config"
	"github.com/yourusername/go-app/internal/handlers"
	customMiddleware "github.com/yourusername/go-app/internal/middleware"
	"github.com/yourusername/go-app/internal/metrics"
	"github.com/yourusername/go-app/pkg/health"
)

const version = "1.0.0"

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize metrics
	m := metrics.New()

	// Initialize router
	r := setupRouter(cfg, m)

	// Configure server
	srv := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s (environment: %s)", cfg.Server.Address(), cfg.Server.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(cfg *config.Config, m *metrics.Metrics) *chi.Mux {
	r := chi.NewRouter()

	// Basic middleware stack
	r.Use(customMiddleware.Recovery)
	r.Use(customMiddleware.Logger)
	r.Use(customMiddleware.Metrics(m))

	// Request ID middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Compression
	r.Use(middleware.Compress(5))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Rate limiting - per IP
	r.Use(httprate.LimitByIP(cfg.RateLimit.RequestsPerSecond, time.Second))

	// Health check endpoints (no rate limiting)
	r.Get("/health", health.Handler(version, time.Now()))
	r.Get("/health/ready", health.ReadinessHandler())
	r.Get("/health/live", health.LivenessHandler())

	// Metrics endpoint
	r.Get("/metrics", handlers.MetricsHandler(m))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Example endpoint
		r.Get("/hello", handlers.HelloHandler)

		// Add your API endpoints here
	})

	// 404 handler
	r.NotFound(handlers.NotFoundHandler)

	return r
}
