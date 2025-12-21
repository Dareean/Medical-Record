package domain

import (
	"strings"
	"time"
)

type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "Pending"
	AppointmentStatusConfirmed AppointmentStatus = "Confirmed"
	AppointmentStatusRejected  AppointmentStatus = "Rejected"
	AppointmentStatusCompleted AppointmentStatus = "Completed"
)

type Appointment struct {
	ID              int               `json:"id"`
	PatientID       int               `json:"patient_id"`
	DoctorID        int               `json:"doctor_id"`
	ScheduleID      *int              `json:"schedule_id,omitempty"`
	AppointmentDate time.Time         `json:"appointment_date"`
	StartTimeSlot   string            `json:"start_time_slot"`
	Complaint       string            `json:"complaint"`
	Status          AppointmentStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`

	// Relations
	Patient  *User           `json:"patient,omitempty"`
	Doctor   *Doctor         `json:"doctor,omitempty"`
	Schedule *DoctorSchedule `json:"schedule,omitempty"`
}

type AppointmentRequest struct {
	PatientID       int    `json:"patient_id" validate:"required"`
	DoctorID        int    `json:"doctor_id" validate:"required"`
	AppointmentDate string `json:"appointment_date" validate:"required"`
	StartTimeSlot   string `json:"start_time_slot" validate:"required"`
	Complaint       string `json:"complaint"`
}

type AppointmentUpdateRequest struct {
	Status AppointmentStatus `json:"status" validate:"required,oneof=Pending Confirmed Rejected Completed"`
}

func NormalizeAppointmentStatus(status string) (AppointmentStatus, bool) {
	switch strings.TrimSpace(strings.ToLower(status)) {
	case strings.ToLower(string(AppointmentStatusPending)):
		return AppointmentStatusPending, true
	case strings.ToLower(string(AppointmentStatusConfirmed)), "approved", "approve",
		"approving", "accept", "accepted":
		return AppointmentStatusConfirmed, true
	case strings.ToLower(string(AppointmentStatusRejected)), "reject", "rejected":
		return AppointmentStatusRejected, true
	case strings.ToLower(string(AppointmentStatusCompleted)), "complete", "completed":
		return AppointmentStatusCompleted, true
	default:
		return AppointmentStatusPending, false
	}
}

func IsValidAppointmentStatus(status string) bool {
	_, ok := NormalizeAppointmentStatus(status)
	return ok
}
