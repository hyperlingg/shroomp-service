package storage

import (
	"os"
	"service/models"
	"testing"
	"time"
)

func TestStore_Create(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	now := time.Now()
	item := models.Item{
		ID:           "test-1",
		MushroomName: "Boletus edulis",
		Location:     "Pacific Northwest",
		Count:        3,
		DateTime:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := store.Create(item)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify item exists
	retrieved, err := store.Get(item.ID)
	if err != nil {
		t.Fatalf("Get after Create failed: %v", err)
	}

	if retrieved.ID != item.ID {
		t.Errorf("Expected ID %s, got %s", item.ID, retrieved.ID)
	}
	if retrieved.MushroomName != item.MushroomName {
		t.Errorf("Expected MushroomName %s, got %s", item.MushroomName, retrieved.MushroomName)
	}
}

func TestStore_CreateDuplicate(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	now := time.Now()
	item := models.Item{
		ID:       "test-1",
		Location: "Test Location",
		Count:    1,
		DateTime: now,
	}

	// First create should succeed
	err := store.Create(item)
	if err != nil {
		t.Fatalf("First Create failed: %v", err)
	}

	// Second create should fail
	err = store.Create(item)
	if err != ErrAlreadyExists {
		t.Errorf("Expected ErrAlreadyExists, got %v", err)
	}
}

func TestStore_Get(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	now := time.Now()
	item := models.Item{
		ID:       "test-1",
		Location: "Test Location",
		Count:    1,
		DateTime: now,
	}

	if err := store.Create(item); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.Get(item.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.ID != item.ID {
		t.Errorf("Expected ID %s, got %s", item.ID, retrieved.ID)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	_, err := store.Get("non-existent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestStore_GetAll(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	now := time.Now()
	items := []models.Item{
		{ID: "test-1", MushroomName: "Item 1", Location: "Location 1", Count: 1, DateTime: now},
		{ID: "test-2", MushroomName: "Item 2", Location: "Location 2", Count: 2, DateTime: now},
		{ID: "test-3", MushroomName: "Item 3", Location: "Location 3", Count: 3, DateTime: now},
	}

	for _, item := range items {
		if err := store.Create(item); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	allItems := store.GetAll()
	if len(allItems) != len(items) {
		t.Errorf("Expected %d items, got %d", len(items), len(allItems))
	}

	// Verify all items are present
	idMap := make(map[string]bool)
	for _, item := range allItems {
		idMap[item.ID] = true
	}

	for _, item := range items {
		if !idMap[item.ID] {
			t.Errorf("Expected item %s not found in GetAll results", item.ID)
		}
	}
}

func TestStore_GetAllEmpty(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	allItems := store.GetAll()
	if len(allItems) != 0 {
		t.Errorf("Expected empty slice, got %d items", len(allItems))
	}
}

func TestStore_Update(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	now := time.Now()
	item := models.Item{
		ID:           "test-1",
		MushroomName: "Original Name",
		Location:     "Original Location",
		Count:        1,
		DateTime:     now,
	}

	if err := store.Create(item); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	updatedItem := models.Item{
		ID:           "test-1",
		MushroomName: "Updated Name",
		Location:     "Updated Location",
		Count:        5,
		DateTime:     now,
	}

	err := store.Update(item.ID, updatedItem)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved, _ := store.Get(item.ID)
	if retrieved.MushroomName != updatedItem.MushroomName {
		t.Errorf("Expected MushroomName %s, got %s", updatedItem.MushroomName, retrieved.MushroomName)
	}
	if retrieved.Location != updatedItem.Location {
		t.Errorf("Expected Location %s, got %s", updatedItem.Location, retrieved.Location)
	}
	if retrieved.Count != updatedItem.Count {
		t.Errorf("Expected Count %d, got %d", updatedItem.Count, retrieved.Count)
	}
}

func TestStore_UpdateNotFound(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	now := time.Now()
	item := models.Item{
		ID:       "non-existent",
		Location: "Test",
		Count:    1,
		DateTime: now,
	}

	err := store.Update(item.ID, item)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestStore_Delete(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	now := time.Now()
	item := models.Item{
		ID:       "test-1",
		Location: "Test Location",
		Count:    1,
		DateTime: now,
	}

	if err := store.Create(item); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err := store.Delete(item.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify item is gone
	_, err = store.Get(item.ID)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestStore_DeleteNotFound(t *testing.T) {
	store := createTestStore(t)
	defer cleanupTestStore(store)

	err := store.Delete("non-existent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestStore_Persistence(t *testing.T) {
	testFile := "test_data.json"

	// Create store and add items
	store1 := &Store{
		items:    make(map[string]models.Item),
		filepath: testFile,
	}
	defer os.Remove(testFile)

	now := time.Now()
	items := []models.Item{
		{ID: "test-1", MushroomName: "Item 1", Location: "Location 1", Count: 1, DateTime: now},
		{ID: "test-2", MushroomName: "Item 2", Location: "Location 2", Count: 2, DateTime: now},
	}

	for _, item := range items {
		err := store1.Create(item)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Create new store instance and load from file
	store2 := &Store{
		items:    make(map[string]models.Item),
		filepath: testFile,
	}
	err := store2.load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify all items were loaded
	if len(store2.items) != len(items) {
		t.Errorf("Expected %d items after load, got %d", len(items), len(store2.items))
	}

	for _, item := range items {
		retrieved, err := store2.Get(item.ID)
		if err != nil {
			t.Errorf("Failed to get item %s after load: %v", item.ID, err)
		}
		if retrieved.MushroomName != item.MushroomName {
			t.Errorf("Expected MushroomName %s, got %s", item.MushroomName, retrieved.MushroomName)
		}
		if retrieved.Location != item.Location {
			t.Errorf("Expected Location %s, got %s", item.Location, retrieved.Location)
		}
	}
}

func TestStore_PersistenceEmptyFile(t *testing.T) {
	testFile := "test_empty.json"
	defer os.Remove(testFile)

	// Create empty file
	if err := os.WriteFile(testFile, []byte{}, 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	store := &Store{
		items:    make(map[string]models.Item),
		filepath: testFile,
	}

	err := store.load()
	if err != nil {
		t.Fatalf("Load empty file failed: %v", err)
	}

	if len(store.items) != 0 {
		t.Errorf("Expected 0 items from empty file, got %d", len(store.items))
	}
}

func TestStore_PersistenceNonExistentFile(t *testing.T) {
	store := &Store{
		items:    make(map[string]models.Item),
		filepath: "non_existent_file.json",
	}

	err := store.load()
	if err != nil {
		t.Fatalf("Load non-existent file should not fail: %v", err)
	}

	if len(store.items) != 0 {
		t.Errorf("Expected 0 items from non-existent file, got %d", len(store.items))
	}
}

// Helper functions

func createTestStore(t *testing.T) *Store {
	testFile := "test_" + t.Name() + ".json"
	return &Store{
		items:    make(map[string]models.Item),
		filepath: testFile,
	}
}

func cleanupTestStore(store *Store) {
	os.Remove(store.filepath)
}
