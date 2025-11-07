package models

import "time"

// Item represents the main data structure for CRUD operations
// Replace this struct with your own custom type as needed
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
