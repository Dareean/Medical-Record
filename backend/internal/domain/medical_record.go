package domain

import (
	"time"
)

type MedicalRecord struct {
	ID              int       `json:"id"`
	AppointmentID   int       `json:"appointment_id"`
	Diagnosis       string    `json:"diagnosis"`
	Prescription    string    `json:"prescription"`
	DoctorNotes     string    `json:"doctor_notes"`
	ExaminationDate time.Time `json:"examination_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	Appointment *Appointment `json:"appointment,omitempty"`
}

type MedicalRecordRequest struct {
	AppointmentID int    `json:"appointment_id" validate:"required"`
	Diagnosis     string `json:"diagnosis"`
	Prescription  string `json:"prescription"`
	DoctorNotes   string `json:"doctor_notes"`
}
