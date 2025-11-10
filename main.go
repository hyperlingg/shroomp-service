package main

import (
	"log"
	"net/http"

	"service/handlers"
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

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("CORS enabled for all origins")
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
