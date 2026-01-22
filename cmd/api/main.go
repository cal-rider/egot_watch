package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"egot-tracker/internal/config"
	"egot-tracker/internal/database"
	"egot-tracker/internal/handler"
	"egot-tracker/internal/repository"
	"egot-tracker/internal/scraper"
	"egot-tracker/internal/service"
	"egot-tracker/pkg/response"
)

// corsMiddleware adds CORS headers to allow frontend requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load .env file if present
	godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create context for database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize database connection pool
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to database")

	// Initialize repositories
	celebrityRepo := repository.NewCelebrityRepository(pool)
	awardRepo := repository.NewAwardRepository(pool)

	// Initialize Wikidata scraper
	wikidataScraper := scraper.NewWikidataScraper()

	// Initialize services
	celebrityService := service.NewCelebrityService(celebrityRepo, awardRepo, wikidataScraper)

	// Initialize handlers
	celebrityHandler := handler.NewCelebrityHandler(celebrityService)

	// Setup routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.HealthCheck(r.Context(), pool); err != nil {
			response.Error(w, http.StatusServiceUnavailable, "database unhealthy")
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "OK"})
	})

	// Celebrity search endpoint
	mux.HandleFunc("GET /api/celebrity/search", celebrityHandler.Search)

	// Celebrity autocomplete endpoint
	mux.HandleFunc("GET /api/celebrity/autocomplete", celebrityHandler.Autocomplete)

	// Close to EGOT endpoint
	mux.HandleFunc("GET /api/celebrity/close-to-egot", celebrityHandler.CloseToEGOT)

	// EGOT winners endpoint
	mux.HandleFunc("GET /api/celebrity/egot-winners", celebrityHandler.EGOTWinners)

	// No awards endpoint
	mux.HandleFunc("GET /api/celebrity/no-awards", celebrityHandler.NoAwards)

	// Create server with CORS middleware
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
