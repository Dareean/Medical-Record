package storage

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	// Create users table
	if err := CreateUsersTable(db); err != nil {
		return err
	}

	// Create specializations table
	if err := CreateSpecializationsTable(db); err != nil {
		return err
	}

	// Create doctors table
	if err := CreateDoctorsTable(db); err != nil {
		return err
	}

	// Create patients table
	if err := CreatePatientsTable(db); err != nil {
		return err
	}

	// Create doctor_schedules table
	if err := CreateDoctorSchedulesTable(db); err != nil {
		return err
	}

	// Create appointments table
	if err := CreateAppointmentsTable(db); err != nil {
		return err
	}

	// Create medical_records table
	if err := CreateMedicalRecordsTable(db); err != nil {
		return err
	}

	log.Println("All migrations completed successfully")
	return nil
}

func CreateUsersTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role ENUM('admin', 'doctor', 'patient') DEFAULT 'patient',
		profile_picture VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	ctx := context.Background()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println("ERROR creating users table:", err)
		return err
	}

	log.Println("Users table created or already exists")
	return nil
}

func CreateSpecializationsTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS specializations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	ctx := context.Background()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println("ERROR creating specializations table:", err)
		return err
	}

	log.Println("Specializations table created or already exists")
	return nil
}

func CreateDoctorsTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS doctors (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		specialization_id INT NOT NULL,
		gender ENUM('male', 'female') NOT NULL,
		address TEXT,
		license_number VARCHAR(255) UNIQUE,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (specialization_id) REFERENCES specializations(id) ON DELETE RESTRICT
	)`

	ctx := context.Background()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println("ERROR creating doctors table:", err)
		return err
	}

	log.Println("Doctors table created or already exists")
	return nil
}

func CreatePatientsTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS patients (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL UNIQUE,
		date_of_birth DATE,
		phone VARCHAR(20),
		address TEXT,
		blood_type ENUM('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-'),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	ctx := context.Background()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println("ERROR creating patients table:", err)
		return err
	}

	log.Println("Patients table created or already exists")
	return nil
}

func CreateDoctorSchedulesTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS doctor_schedules (
		id INT AUTO_INCREMENT PRIMARY KEY,
		doctor_id INT NOT NULL,
		work_day ENUM('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday') NOT NULL,
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		patient_quota INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
		UNIQUE KEY unique_schedule (doctor_id, work_day, start_time, end_time)
	)`

	ctx := context.Background()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println("ERROR creating doctor_schedules table:", err)
		return err
	}

	log.Println("Doctor schedules table created or already exists")
	return nil
}

func CreateAppointmentsTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS appointments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		patient_id INT NOT NULL,
		doctor_id INT NOT NULL,
		schedule_id INT NULL,
		appointment_date DATE NOT NULL,
		start_time_slot TIME NOT NULL,
		complaint TEXT,
		status ENUM('Pending', 'Confirmed', 'Rejected', 'Completed') DEFAULT 'Pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
		FOREIGN KEY (schedule_id) REFERENCES doctor_schedules(id) ON DELETE CASCADE
	)`

	ctx := context.Background()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println("ERROR creating appointments table:", err)
		return err
	}

	if err := ensureAppointmentColumns(db); err != nil {
		return err
	}

	log.Println("Appointments table created or already exists")
	return nil
}

func ensureAppointmentColumns(db *sql.DB) error {
	ctx := context.Background()
	statements := []string{
		"ALTER TABLE appointments ADD COLUMN IF NOT EXISTS start_time_slot TIME NOT NULL DEFAULT '00:00:00' AFTER appointment_date",
		"ALTER TABLE appointments ADD COLUMN IF NOT EXISTS complaint TEXT AFTER start_time_slot",
		"ALTER TABLE appointments MODIFY COLUMN schedule_id INT NULL",
	}

	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			if strings.Contains(err.Error(), "Duplicate column name") {
				continue
			}
			log.Println("ERROR ensuring appointments column:", err)
			return err
		}
	}

	return nil
}

func CreateMedicalRecordsTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS medical_records (
		id INT AUTO_INCREMENT PRIMARY KEY,
		appointment_id INT NOT NULL,
		diagnosis TEXT,
		prescription TEXT,
		doctor_notes TEXT,
		examination_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE
	)`

	ctx := context.Background()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Println("ERROR creating medical_records table:", err)
		return err
	}

	log.Println("Medical records table created or already exists")
	return nil
}
