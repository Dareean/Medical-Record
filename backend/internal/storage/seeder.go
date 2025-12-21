package storage

import (
	"context"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedAllData(db *sql.DB) error {
	if err := SeedDefaultUsers(db); err != nil {
		log.Println("WARNING: Failed to seed default users:", err)
	}

	if err := SeedDefaultSpecializations(db); err != nil {
		log.Println("WARNING: Failed to seed default specializations:", err)
	}

	if err := SeedDefaultDoctors(db); err != nil {
		log.Println("WARNING: Failed to seed default doctors:", err)
	}

	if err := SeedDefaultPatients(db); err != nil {
		log.Println("WARNING: Failed to seed default patients:", err)
	}

	if err := SeedDefaultDoctorSchedules(db); err != nil {
		log.Println("WARNING: Failed to seed default doctor schedules:", err)
	}

	return nil
}

func SeedDefaultUsers(db *sql.DB) error {
	ctx := context.Background()

	// Hash passwords
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	doctorPassword, _ := bcrypt.GenerateFromPassword([]byte("doctor123"), bcrypt.DefaultCost)
	patientPassword, _ := bcrypt.GenerateFromPassword([]byte("patient123"), bcrypt.DefaultCost)

	users := []struct {
		name     string
		email    string
		password string
		role     string
	}{
		{"Admin User", "admin@hospital.com", string(adminPassword), "admin"},
		{"Dr. John Smith", "john.smith@hospital.com", string(doctorPassword), "doctor"},
		{"Dr. Sarah Johnson", "sarah.johnson@hospital.com", string(doctorPassword), "doctor"},
		{"Dr. Michael Brown", "michael.brown@hospital.com", string(doctorPassword), "doctor"},
		{"Patient One", "patient1@email.com", string(patientPassword), "patient"},
		{"Patient Two", "patient2@email.com", string(patientPassword), "patient"},
		{"Patient Three", "patient3@email.com", string(patientPassword), "patient"},
	}

	for _, user := range users {
		_, err := db.ExecContext(ctx,
			"INSERT IGNORE INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
			user.name, user.email, user.password, user.role)
		if err != nil {
			log.Printf("WARNING: Failed to insert user %s: %v", user.email, err)
		}
	}

	log.Println("Default users seeded")
	return nil
}

func SeedDefaultSpecializations(db *sql.DB) error {
	specializations := []string{
		"General Practitioner",
		"Pediatrician",
		"Cardiologist",
		"Dermatologist",
		"Neurologist",
		"Orthopedist",
		"Ophthalmologist",
		"Gynecologist",
		"Psychiatrist",
		"Urologist",
	}

	ctx := context.Background()
	for _, spec := range specializations {
		_, err := db.ExecContext(ctx, "INSERT IGNORE INTO specializations (name) VALUES (?)", spec)
		if err != nil {
			log.Printf("WARNING: Failed to insert specialization %s: %v", spec, err)
		}
	}

	log.Println("Default specializations inserted")
	return nil
}

func SeedDefaultDoctors(db *sql.DB) error {
	ctx := context.Background()

	doctors := []struct {
		userID           int
		specializationID int
		gender           string
		address          string
		licenseNumber    string
	}{
		{2, 1, "male", "123 Main St, City, State", "DOC001"},
		{3, 2, "female", "456 Oak Ave, City, State", "DOC002"},
		{4, 3, "male", "789 Pine Rd, City, State", "DOC003"},
	}

	for _, doctor := range doctors {
		_, err := db.ExecContext(ctx,
			"INSERT IGNORE INTO doctors (user_id, specialization_id, gender, address, license_number) VALUES (?, ?, ?, ?, ?)",
			doctor.userID, doctor.specializationID, doctor.gender, doctor.address, doctor.licenseNumber)
		if err != nil {
			log.Printf("WARNING: Failed to insert doctor %d: %v", doctor.userID, err)
		}
	}

	log.Println("Default doctors seeded")
	return nil
}

func SeedDefaultPatients(db *sql.DB) error {
	ctx := context.Background()

	patients := []struct {
		userID      int
		dateOfBirth string
		phone       string
		address     string
		bloodType   string
	}{
		{5, "1990-01-15", "555-0101", "111 Patient St, City, State", "A+"},
		{6, "1985-05-22", "555-0102", "222 Patient Ave, City, State", "B+"},
		{7, "1992-11-08", "555-0103", "333 Patient Rd, City, State", "O-"},
	}

	for _, patient := range patients {
		_, err := db.ExecContext(ctx,
			"INSERT IGNORE INTO patients (user_id, date_of_birth, phone, address, blood_type) VALUES (?, ?, ?, ?, ?)",
			patient.userID, patient.dateOfBirth, patient.phone, patient.address, patient.bloodType)
		if err != nil {
			log.Printf("WARNING: Failed to insert patient %d: %v", patient.userID, err)
		}
	}

	log.Println("Default patients seeded")
	return nil
}

func SeedDefaultDoctorSchedules(db *sql.DB) error {
	ctx := context.Background()

	schedules := []struct {
		doctorID  int
		workDay   string
		startTime string
		endTime   string
		quota     int
	}{
		// Dr. John Smith (ID: 1) - General Practitioner
		{1, "monday", "09:00", "12:00", 10},
		{1, "monday", "14:00", "17:00", 10},
		{1, "wednesday", "09:00", "12:00", 10},
		{1, "wednesday", "14:00", "17:00", 10},
		{1, "friday", "09:00", "12:00", 10},
		{1, "friday", "14:00", "17:00", 10},

		// Dr. Sarah Johnson (ID: 2) - Pediatrician
		{2, "tuesday", "08:00", "12:00", 8},
		{2, "tuesday", "13:00", "16:00", 8},
		{2, "thursday", "08:00", "12:00", 8},
		{2, "thursday", "13:00", "16:00", 8},
		{2, "saturday", "08:00", "12:00", 8},

		// Dr. Michael Brown (ID: 3) - Cardiologist
		{3, "monday", "10:00", "13:00", 12},
		{3, "tuesday", "10:00", "13:00", 12},
		{3, "wednesday", "10:00", "13:00", 12},
		{3, "thursday", "10:00", "13:00", 12},
		{3, "friday", "10:00", "13:00", 12},
	}

	for _, schedule := range schedules {
		_, err := db.ExecContext(ctx,
			"INSERT IGNORE INTO doctor_schedules (doctor_id, work_day, start_time, end_time, patient_quota) VALUES (?, ?, ?, ?, ?)",
			schedule.doctorID, schedule.workDay, schedule.startTime, schedule.endTime, schedule.quota)
		if err != nil {
			log.Printf("WARNING: Failed to insert schedule for doctor %d: %v", schedule.doctorID, err)
		}
	}

	log.Println("Default doctor schedules seeded")
	return nil
}
