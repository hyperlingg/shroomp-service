package main

import (
	"net/http"
	"os"

	"service/handlers"
	"service/logger"
	"service/storage"
)

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
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
	// Initialize the storage
	store := storage.NewStore()

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(store)

	// Setup routes
	mux := http.NewServeMux()

	// CRUD endpoints
	mux.HandleFunc("/items", itemHandler.HandleItems)
	mux.HandleFunc("/items/", itemHandler.HandleItemByID)

	// Wrap mux with CORS middleware
	handler := corsMiddleware(mux)

	// Get port from environment variable (for Cloud Run) or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	logger.Info("Server starting", map[string]interface{}{
		"port": port,
	})
	logger.Info("CORS configuration", map[string]interface{}{
		"enabled": true,
		"origins": "*",
	})
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatal("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
