package domain

import (
	"time"
)

type BloodType string

const (
	BloodTypeAPlus   BloodType = "A+"
	BloodTypeAMinus  BloodType = "A-"
	BloodTypeBPlus   BloodType = "B+"
	BloodTypeBMinus  BloodType = "B-"
	BloodTypeABPlus  BloodType = "AB+"
	BloodTypeABMinus BloodType = "AB-"
	BloodTypeOPlus   BloodType = "O+"
	BloodTypeOMinus  BloodType = "O-"
)

type Patient struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	BloodType   BloodType `json:"blood_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty"`
}

type PatientRequest struct {
	UserID      int    `json:"user_id" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	BloodType   string `json:"blood_type" validate:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
}

type PatientUpdateRequest struct {
	DateOfBirth string `json:"date_of_birth"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	BloodType   string `json:"blood_type" validate:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
}

func IsValidBloodType(bloodType string) bool {
	switch bloodType {
	case string(BloodTypeAPlus), string(BloodTypeAMinus),
		string(BloodTypeBPlus), string(BloodTypeBMinus),
		string(BloodTypeABPlus), string(BloodTypeABMinus),
		string(BloodTypeOPlus), string(BloodTypeOMinus):
		return true
	default:
		return false
	}
}
