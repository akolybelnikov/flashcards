package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	flashcardService := services.NewFlashcardService(flashcardRepo)
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

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
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
