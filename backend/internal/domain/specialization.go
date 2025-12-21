package domain

import "time"

// Specialization represents medical specialization
type Specialization struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request for Specialization
type SpecializationRequest struct {
	Name string `json:"name" validate:"required"`
}
