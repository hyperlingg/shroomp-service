package storage

import (
	"encoding/json"
	"errors"
	"os"
	"service/models"
	"sync"
)

var (
	ErrNotFound      = errors.New("item not found")
	ErrAlreadyExists = errors.New("item already exists")
)

// Store provides thread-safe storage for items with JSON file persistence
type Store struct {
	mu       sync.RWMutex
	items    map[string]models.Item
	filepath string
}

// NewStore creates a new storage instance and loads existing data from file
func NewStore() *Store {
	s := &Store{
		items:    make(map[string]models.Item),
		filepath: "data.json",
	}
	s.load()
	return s
}

// load reads items from the JSON file
func (s *Store) load() error {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's ok
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &s.items)
}

// save writes items to the JSON file
func (s *Store) save() error {
	data, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, data, 0644)
}

// Create adds a new item to the store
func (s *Store) Create(item models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[item.ID]; exists {
		return ErrAlreadyExists
	}

	s.items[item.ID] = item
	return s.save()
}

// Get retrieves an item by ID
func (s *Store) Get(id string) (models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return models.Item{}, ErrNotFound
	}

	return item, nil
}

// GetAll retrieves all items
func (s *Store) GetAll() []models.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]models.Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return items
}

// Update modifies an existing item
func (s *Store) Update(id string, item models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return ErrNotFound
	}

	s.items[id] = item
	return s.save()
}

// Delete removes an item by ID
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return ErrNotFound
	}

	delete(s.items, id)
	return s.save()
}
