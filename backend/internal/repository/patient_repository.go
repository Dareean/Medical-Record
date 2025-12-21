package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
)

// PatientRepository exposes helpers for working with the patients table.
type PatientRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*domain.Patient, error)
	CreateForUser(ctx context.Context, userID int64) (*domain.Patient, error)
}

type patientRepoMySQL struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) PatientRepository {
	return &patientRepoMySQL{db: db}
}

func (r *patientRepoMySQL) GetByUserID(ctx context.Context, userID int64) (*domain.Patient, error) {
	const q = `
        SELECT id, user_id, date_of_birth, phone, address, blood_type,
               created_at, updated_at
        FROM patients
        WHERE user_id = ?
    `

	var (
		patient   domain.Patient
		dob       sql.NullTime
		phone     sql.NullString
		address   sql.NullString
		bloodType sql.NullString
	)

	err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&patient.ID,
		&patient.UserID,
		&dob,
		&phone,
		&address,
		&bloodType,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if dob.Valid {
		patient.DateOfBirth = dob.Time
	}
	if phone.Valid {
		patient.Phone = phone.String
	}
	if address.Valid {
		patient.Address = address.String
	}
	if bloodType.Valid {
		patient.BloodType = domain.BloodType(bloodType.String)
	}

	return &patient, nil
}

func (r *patientRepoMySQL) CreateForUser(ctx context.Context, userID int64) (*domain.Patient, error) {
	const q = `
        INSERT INTO patients (user_id, created_at, updated_at)
        VALUES (?, NOW(), NOW())
    `

	res, err := r.db.ExecContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	patient := &domain.Patient{
		ID:        int(insertID),
		UserID:    int(userID),
		CreatedAt: now,
		UpdatedAt: now,
	}
	return patient, nil
}
