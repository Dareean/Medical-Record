package domain

import (
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

func IsValidAppointmentStatus(status string) bool {
	switch status {
	case string(AppointmentStatusPending), string(AppointmentStatusConfirmed),
		string(AppointmentStatusRejected), string(AppointmentStatusCompleted):
		return true
	default:
		return false
	}
}
