package main

import (
	"log"
	"net/http"

	"service/handlers"
	"service/storage"
)

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

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
