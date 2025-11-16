package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akolybelnikov/flashcards/config"
	"github.com/akolybelnikov/flashcards/db"
	"github.com/akolybelnikov/flashcards/handlers"
	"github.com/akolybelnikov/flashcards/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	// Initialize database connection
	dbConn, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize flashcard components
	flashcardRepo := db.NewPostgresFlashcardRepository(dbConn)
	if flashcardRepo == nil {
		log.Fatal("Failed to initialize flashcard repository")
	}

	// Initialize LLM client if API key is provided
	var llmClient services.LLMClient
	if cfg.OpenAIAPIKey != "" {
		client, err := services.NewOpenAIClient(cfg.OpenAIAPIKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize AI translation: %v", err)
			log.Println("AI translation features will be disabled")
		} else {
			llmClient = client
			log.Println("AI translation enabled")
		}
	} else {
		log.Println("AI translation disabled (OPENAI_API_KEY not set)")
	}

	// Initialize translation cache using go-cache with TTL and cleanup interval from config
	translationCache := services.NewGoCacheTranslationCache(cfg.TranslationCacheTTL, cfg.TranslationCacheCleanupInterval)
	log.Printf("Translation cache initialized with TTL: %s, cleanup interval: %s", cfg.TranslationCacheTTL, cfg.TranslationCacheCleanupInterval)

	// Start a cache cleanup goroutine (kept for compatibility; StartCleanup will noop if interval <= 0)
	stopCleanup := make(chan struct{})
	translationCache.StartCleanup(cfg.TranslationCacheCleanupInterval, stopCleanup)
	log.Printf("Translation cache cleanup started (interval: %s)", cfg.TranslationCacheCleanupInterval)

	flashcardService := services.NewFlashcardService(flashcardRepo, llmClient, translationCache)
	if flashcardService == nil {
		log.Fatal("Failed to initialize flashcard service")
	}

	flashcardHandler := handlers.NewFlashcardHandler(flashcardService)

	router := mux.NewRouter()

	// Add panic recovery middleware first so it can catch panics from other middlewares/handlers
	router.Use(recoverMiddleware)
	router.Use(corsMiddleware)
	router.Use(jsonMiddleware)

	flashcardHandler.RegisterRoutes(router)

	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	addr := ":" + cfg.Port
	fmt.Printf("Server starting on port %s\n", cfg.Port)

	// Create HTTP server
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop cache cleanup
	close(stopCleanup)
	log.Println("Translation cache cleanup stopped")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal_server_error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status": "healthy"}`))
	if err != nil {
		return
	}
}
