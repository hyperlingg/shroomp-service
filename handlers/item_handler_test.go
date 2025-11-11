package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"service/models"
	"service/storage"
	"testing"
	"time"
)

func createTestHandler(t *testing.T) (*ItemHandler, func()) {
	testFile := "test_handler_" + t.Name() + ".json"
	store := &storage.Store{}
	// Initialize store with test file
	os.Remove(testFile) // Clean up any existing test file

	store = storage.NewStore()
	handler := NewItemHandler(store)

	cleanup := func() {
		os.Remove("data.json") // NewStore uses "data.json" as default
	}

	return handler, cleanup
}

func TestHandleItems_POST_Success(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	now := time.Now()
	item := models.Item{
		MushroomName: "Chanterelle",
		Location:     "Forest Trail",
		Count:        5,
		DateTime:     now,
	}

	body, _ := json.Marshal(item)
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleItems(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var created models.Item
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if created.ID == "" {
		t.Error("Expected ID to be generated")
	}
	if created.MushroomName != item.MushroomName {
		t.Errorf("Expected MushroomName %s, got %s", item.MushroomName, created.MushroomName)
	}
	if created.Location != item.Location {
		t.Errorf("Expected Location %s, got %s", item.Location, created.Location)
	}
}

func TestHandleItems_POST_ValidationError(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	tests := []struct {
		name     string
		item     models.Item
		expected string
	}{
		{
			name: "missing location",
			item: models.Item{
				Count:    5,
				DateTime: time.Now(),
			},
			expected: "location is required",
		},
		{
			name: "invalid count",
			item: models.Item{
				Location: "Forest",
				Count:    0,
				DateTime: time.Now(),
			},
			expected: "count must be at least 1",
		},
		{
			name: "missing dateTime",
			item: models.Item{
				Location: "Forest",
				Count:    5,
				DateTime: time.Time{},
			},
			expected: "dateTime is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.item)
			req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleItems(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestHandleItems_POST_InvalidJSON(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleItems(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleItems_GET_Success(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	// Create some test items
	now := time.Now()
	items := []models.Item{
		{ID: "test-1", MushroomName: "Chanterelle", Location: "Forest", Count: 5, DateTime: now, CreatedAt: now, UpdatedAt: now},
		{ID: "test-2", MushroomName: "Morel", Location: "Woods", Count: 3, DateTime: now, CreatedAt: now, UpdatedAt: now},
	}

	for _, item := range items {
		handler.store.Create(item)
	}

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	w := httptest.NewRecorder()

	handler.HandleItems(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var retrieved []models.Item
	if err := json.NewDecoder(w.Body).Decode(&retrieved); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(retrieved) != len(items) {
		t.Errorf("Expected %d items, got %d", len(items), len(retrieved))
	}
}

func TestHandleItems_MethodNotAllowed(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPut, "/items", nil)
	w := httptest.NewRecorder()

	handler.HandleItems(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandleItemByID_GET_Success(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	now := time.Now()
	item := models.Item{
		ID:           "test-1",
		MushroomName: "Chanterelle",
		Location:     "Forest",
		Count:        5,
		DateTime:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	handler.store.Create(item)

	req := httptest.NewRequest(http.MethodGet, "/items/test-1", nil)
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var retrieved models.Item
	if err := json.NewDecoder(w.Body).Decode(&retrieved); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if retrieved.ID != item.ID {
		t.Errorf("Expected ID %s, got %s", item.ID, retrieved.ID)
	}
}

func TestHandleItemByID_GET_NotFound(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/items/non-existent", nil)
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleItemByID_PUT_Success(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	now := time.Now()
	item := models.Item{
		ID:           "test-1",
		MushroomName: "Original",
		Location:     "Original Location",
		Count:        5,
		DateTime:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	handler.store.Create(item)

	updatedItem := models.Item{
		MushroomName: "Updated",
		Location:     "Updated Location",
		Count:        10,
		DateTime:     now,
	}

	body, _ := json.Marshal(updatedItem)
	req := httptest.NewRequest(http.MethodPut, "/items/test-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var retrieved models.Item
	if err := json.NewDecoder(w.Body).Decode(&retrieved); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if retrieved.MushroomName != updatedItem.MushroomName {
		t.Errorf("Expected MushroomName %s, got %s", updatedItem.MushroomName, retrieved.MushroomName)
	}
}

func TestHandleItemByID_PUT_NotFound(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	now := time.Now()
	item := models.Item{
		MushroomName: "Test",
		Location:     "Location",
		Count:        5,
		DateTime:     now,
	}

	body, _ := json.Marshal(item)
	req := httptest.NewRequest(http.MethodPut, "/items/non-existent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleItemByID_PUT_ValidationError(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	now := time.Now()
	item := models.Item{
		ID:           "test-1",
		MushroomName: "Test",
		Location:     "Location",
		Count:        5,
		DateTime:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	handler.store.Create(item)

	invalidItem := models.Item{
		Location: "", // Invalid: empty location
		Count:    5,
		DateTime: now,
	}

	body, _ := json.Marshal(invalidItem)
	req := httptest.NewRequest(http.MethodPut, "/items/test-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleItemByID_PUT_InvalidJSON(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPut, "/items/test-1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleItemByID_DELETE_Success(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	now := time.Now()
	item := models.Item{
		ID:           "test-1",
		MushroomName: "Test",
		Location:     "Location",
		Count:        5,
		DateTime:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	handler.store.Create(item)

	req := httptest.NewRequest(http.MethodDelete, "/items/test-1", nil)
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify item is deleted
	_, err := handler.store.Get("test-1")
	if err != storage.ErrNotFound {
		t.Errorf("Expected item to be deleted, but still exists")
	}
}

func TestHandleItemByID_DELETE_NotFound(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodDelete, "/items/non-existent", nil)
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleItemByID_EmptyID(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/items/", nil)
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleItemByID_MethodNotAllowed(t *testing.T) {
	handler, cleanup := createTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/items/test-1", nil)
	w := httptest.NewRecorder()

	handler.HandleItemByID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestValidateSighting(t *testing.T) {
	tests := []struct {
		name      string
		item      models.Item
		expectErr bool
	}{
		{
			name: "valid item",
			item: models.Item{
				Location: "Forest",
				Count:    5,
				DateTime: time.Now(),
			},
			expectErr: false,
		},
		{
			name: "empty location",
			item: models.Item{
				Location: "",
				Count:    5,
				DateTime: time.Now(),
			},
			expectErr: true,
		},
		{
			name: "count zero",
			item: models.Item{
				Location: "Forest",
				Count:    0,
				DateTime: time.Now(),
			},
			expectErr: true,
		},
		{
			name: "zero dateTime",
			item: models.Item{
				Location: "Forest",
				Count:    5,
				DateTime: time.Time{},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSighting(&tt.item)
			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
