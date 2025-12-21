package domain

import "time"

// Gender enum
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

func IsValidGender(gender string) bool {
	switch Gender(gender) {
	case GenderMale, GenderFemale:
		return true
	default:
		return false
	}
}

// Doctor represents doctor information
type Doctor struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	SpecializationID int       `json:"specialization_id" `
	Gender           Gender    `json:"gender"`
	Address          string    `json:"address"`
	LicenseNumber    string    `json:"license_number"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relations
	User           *User           `json:"user,omitempty"`
	Specialization *Specialization `json:"specialization,omitempty"`
}

// Request for Doctor
type DoctorRequest struct {
	Name             string `json:"name" validate:"required"`
	Email            string `json:"email" validate:"required,email"`
	SpecializationID int    `json:"specialization_id" validate:"required"`
	Gender           string `json:"gender" validate:"required,oneof=male female"`
	Address          string `json:"address"`
	LicenseNumber    string `json:"license_number" validate:"required"`
}
