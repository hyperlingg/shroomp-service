package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"service/logger"
	"service/models"
	"service/storage"

	"github.com/google/uuid"
)

type ItemHandler struct {
	store *storage.Store
}

func NewItemHandler(store *storage.Store) *ItemHandler {
	return &ItemHandler{store: store}
}

// HandleItems handles POST (create) and GET (list all) requests
func (h *ItemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createItem(w, r)
	case http.MethodGet:
		h.getAllItems(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleItemByID handles GET (read), PUT (update), and DELETE requests for specific items
func (h *ItemHandler) HandleItemByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	if id == "" {
		http.Error(w, "Item ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getItem(w, r, id)
	case http.MethodPut:
		h.updateItem(w, r, id)
	case http.MethodDelete:
		h.deleteItem(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// validateSighting validates required fields for a mushroom sighting
func validateSighting(item *models.Item) error {
	if item.Location == "" {
		return errors.New("location is required")
	}
	if item.Count < 1 {
		return errors.New("count must be at least 1")
	}
	if item.DateTime.IsZero() {
		return errors.New("dateTime is required")
	}
	return nil
}

// createItem creates a new item
func (h *ItemHandler) createItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := validateSighting(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate ID if not provided
	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := h.store.Create(item); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			http.Error(w, "Item already exists", http.StatusConflict)
			return
		}
		logger.Error("Error creating item", map[string]interface{}{
			"error":    err.Error(),
			"item_id":  item.ID,
			"location": item.Location,
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		logger.Error("Failed to encode response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// getAllItems retrieves all items
func (h *ItemHandler) getAllItems(w http.ResponseWriter, r *http.Request) {
	items := h.store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		logger.Error("Failed to encode response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// getItem retrieves a specific item by ID
func (h *ItemHandler) getItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.store.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		logger.Error("Error getting item", map[string]interface{}{
			"error":   err.Error(),
			"item_id": id,
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		logger.Error("Failed to encode response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// updateItem updates an existing item
func (h *ItemHandler) updateItem(w http.ResponseWriter, r *http.Request, id string) {
	var item models.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := validateSighting(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preserve the ID from the URL
	item.ID = id
	item.UpdatedAt = time.Now()

	if err := h.store.Update(id, item); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		logger.Error("Error updating item", map[string]interface{}{
			"error":    err.Error(),
			"item_id":  id,
			"location": item.Location,
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		logger.Error("Failed to encode response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// deleteItem removes an item
func (h *ItemHandler) deleteItem(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		logger.Error("Error deleting item", map[string]interface{}{
			"error":   err.Error(),
			"item_id": id,
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
