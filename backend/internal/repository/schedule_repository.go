package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
)

type ScheduleRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.DoctorSchedule, error)
	GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*domain.DoctorSchedule, error)
}

type scheduleRepoMySQL struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepoMySQL{db: db}
}

func (r *scheduleRepoMySQL) GetByID(ctx context.Context, id int64) (*domain.DoctorSchedule, error) {
	var s domain.DoctorSchedule
	q := `SELECT id, doctor_id, work_day, start_time, end_time, patient_quota, created_at, updated_at
	      FROM doctor_schedules WHERE id = ?`
	row := r.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(&s.ID, &s.DoctorID, &s.WorkDay, &s.StartTime, &s.EndTime, &s.PatientQuota, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, sql.ErrNoRows }
		return nil, err
	}
	return &s, nil
}

func (r *scheduleRepoMySQL) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*domain.DoctorSchedule, error) {
	var s domain.DoctorSchedule
	q := `SELECT id, doctor_id, work_day, start_time, end_time, patient_quota, created_at, updated_at
	      FROM doctor_schedules WHERE id = ? FOR UPDATE`
	row := tx.QueryRowContext(ctx, q, id)
	if err := row.Scan(&s.ID, &s.DoctorID, &s.WorkDay, &s.StartTime, &s.EndTime, &s.PatientQuota, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, sql.ErrNoRows }
		return nil, err
	}
	return &s, nil
}