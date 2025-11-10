package models

import "time"

// MushroomSighting represents a mushroom sighting record
type MushroomSighting struct {
	ID           string    `json:"id"`
	Image        *string   `json:"image,omitempty"`        // Optional base64 encoded image
	MushroomName string    `json:"mushroomName,omitempty"` // Optional user identification
	DateTime     time.Time `json:"dateTime"`               // When the mushroom was found
	Location     string    `json:"location"`               // Where the mushroom was found
	Count        int       `json:"count"`                  // Number of mushrooms found
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Item is kept for backwards compatibility, aliased to MushroomSighting
type Item = MushroomSighting
