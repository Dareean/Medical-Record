package domain

import "time"

// WorkDay enum
type WorkDay string

const (
	WorkDayMonday    WorkDay = "monday"
	WorkDayTuesday   WorkDay = "tuesday"
	WorkDayWednesday WorkDay = "wednesday"
	WorkDayThursday  WorkDay = "thursday"
	WorkDayFriday    WorkDay = "friday"
	WorkDaySaturday  WorkDay = "saturday"
	WorkDaySunday    WorkDay = "sunday"
)

func IsValidWorkDay(day string) bool {
	switch WorkDay(day) {
	case
		WorkDayMonday, WorkDayTuesday,
		WorkDayWednesday, WorkDayThursday,
		WorkDayFriday, WorkDaySaturday, WorkDaySunday:
		return true
	default:
		return false
	}
}

// DoctorSchedule represents doctor working schedule
type DoctorSchedule struct {
	ID           int       `json:"id"`
	DoctorID     int       `json:"doctor_id"`
	WorkDay      WorkDay   `json:"work_day"`
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
	PatientQuota int       `json:"patient_quota"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Doctor *Doctor `json:"doctor,omitempty"`
}

// Request for DoctorSchedule
type DoctorScheduleRequest struct {
	DoctorID     int    `json:"doctor_id" validate:"required"`
	WorkDay      string `json:"work_day" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime    string `json:"start_time" validate:"required"`
	EndTime      string `json:"end_time" validate:"required"`
	PatientQuota int    `json:"patient_quota"`
}
