package domain

import "time"

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Role           UserRole  `json:"role"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleDoctor  UserRole = "doctor"
	RolePatient UserRole = "patient"
)

func IsValidRole(role string) bool {
	switch UserRole(role) {
	case RoleAdmin, RoleDoctor, RolePatient:
		return true
	default:
		return false
	}
}

type RegisterRequest struct {
	Name           string   `json:"name" validate:"required"`
	Email          string   `json:"email" validate:"required,email"`
	Password       string   `json:"password" validate:"required,min=6"`
	Role           UserRole `json:"role" validate:"omitempty,oneof=admin doctor patient"`
	ProfilePicture string   `json:"profile_picture"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	User      User      `json:"user"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
