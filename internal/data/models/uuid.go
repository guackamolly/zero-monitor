package models

import "github.com/google/uuid"

// Returns a random v4 UUID.
// If an error occurs, it panics.
func UUID() string {
	return uuid.NewString()
}
