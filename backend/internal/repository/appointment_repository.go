package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
)

type AppointmentRepository interface {
	CreateTx(ctx context.Context, tx *sql.Tx, a *domain.Appointment) error
	GetByID(ctx context.Context, id int64) (*domain.Appointment, error)
	UpdateStatusTx(ctx context.Context, tx *sql.Tx, id int64, status domain.AppointmentStatus) error
	GetByPatient(ctx context.Context, patientID int64) ([]domain.Appointment, error)
	GetByDoctor(ctx context.Context, doctorID int64) ([]domain.Appointment, error)
}

type appointmentRepoMySQL struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) AppointmentRepository {
	return &appointmentRepoMySQL{db: db}
}

var _ AppointmentRepository = (*appointmentRepoMySQL)(nil)

func (r *appointmentRepoMySQL) CreateTx(ctx context.Context, tx *sql.Tx, a *domain.Appointment) error {
	const q = `
		INSERT INTO appointments
			(patient_id, doctor_id, schedule_id, appointment_date, start_time_slot, complaint, status, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	var scheduleID interface{}
	if a.ScheduleID != nil {
		scheduleID = *a.ScheduleID
	} else {
		scheduleID = nil
	}
	res, err := tx.ExecContext(
		ctx,
		q,
		a.PatientID,
		a.DoctorID,
		scheduleID,
		a.AppointmentDate,
		a.StartTimeSlot,
		a.Complaint,
		a.Status,
	)
	if err != nil {
		return err
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	now := time.Now()
	a.ID = int(insertID)
	a.CreatedAt = now
	a.UpdatedAt = now

	return nil
}

func (r *appointmentRepoMySQL) GetByID(ctx context.Context, id int64) (*domain.Appointment, error) {
	const q = `
		SELECT a.id, a.patient_id, a.doctor_id, a.schedule_id,
		       a.appointment_date, a.start_time_slot, a.complaint,
		       a.status, a.created_at, a.updated_at,
		       du.name, du.email,
		       pu.name, pu.email
		FROM appointments a
		LEFT JOIN doctors d ON a.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		LEFT JOIN patients p ON p.id = a.patient_id
		LEFT JOIN users pu ON p.user_id = pu.id
		WHERE a.id = ?
	`

	row := r.db.QueryRowContext(ctx, q, id)
	var (
		a            domain.Appointment
		idDB         int64
		scheduleID   sql.NullInt64
		startTime    sql.NullString
		complaint    sql.NullString
		doctorName   sql.NullString
		doctorEmail  sql.NullString
		patientName  sql.NullString
		patientEmail sql.NullString
	)

	if err := row.Scan(
		&idDB,
		&a.PatientID,
		&a.DoctorID,
		&scheduleID,
		&a.AppointmentDate,
		&startTime,
		&complaint,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
		&doctorName,
		&doctorEmail,
		&patientName,
		&patientEmail,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	a.ID = int(idDB)
	if scheduleID.Valid {
		v := int(scheduleID.Int64)
		a.ScheduleID = &v
	}
	if startTime.Valid {
		a.StartTimeSlot = startTime.String
	}
	if complaint.Valid {
		a.Complaint = complaint.String
	}
	if doctorName.Valid || doctorEmail.Valid {
		a.Doctor = &domain.Doctor{
			ID: a.DoctorID,
			User: &domain.User{
				Name:  doctorName.String,
				Email: doctorEmail.String,
			},
		}
	}
	if patientName.Valid || patientEmail.Valid {
		a.Patient = &domain.User{
			Name:  patientName.String,
			Email: patientEmail.String,
		}
	}

	return &a, nil
}

func (r *appointmentRepoMySQL) UpdateStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	status domain.AppointmentStatus,
) error {
	const q = `
		UPDATE appointments
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`
	res, err := tx.ExecContext(ctx, q, status, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *appointmentRepoMySQL) GetByPatient(ctx context.Context, patientID int64) ([]domain.Appointment, error) {
	const q = `
		SELECT a.id, a.patient_id, a.doctor_id, a.schedule_id,
		       a.appointment_date, a.start_time_slot, a.complaint,
		       a.status, a.created_at, a.updated_at,
		       du.name, du.email
		FROM appointments a
		LEFT JOIN doctors d ON a.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE a.patient_id = ?
		ORDER BY a.appointment_date DESC, a.start_time_slot DESC
	`
	rows, err := r.db.QueryContext(ctx, q, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Appointment
	for rows.Next() {
		var (
			a           domain.Appointment
			idDB        int64
			scheduleID  sql.NullInt64
			startTime   sql.NullString
			complaint   sql.NullString
			doctorName  sql.NullString
			doctorEmail sql.NullString
		)
		if err := rows.Scan(
			&idDB,
			&a.PatientID,
			&a.DoctorID,
			&scheduleID,
			&a.AppointmentDate,
			&startTime,
			&complaint,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
			&doctorName,
			&doctorEmail,
		); err != nil {
			return nil, err
		}
		a.ID = int(idDB)
		if scheduleID.Valid {
			v := int(scheduleID.Int64)
			a.ScheduleID = &v
		}
		if startTime.Valid {
			a.StartTimeSlot = startTime.String
		}
		if complaint.Valid {
			a.Complaint = complaint.String
		}
		if doctorName.Valid || doctorEmail.Valid {
			a.Doctor = &domain.Doctor{
				ID: a.DoctorID,
				User: &domain.User{
					Name:  doctorName.String,
					Email: doctorEmail.String,
				},
			}
		}
		result = append(result, a)
	}

	return result, nil
}

func (r *appointmentRepoMySQL) GetByDoctor(ctx context.Context, doctorID int64) ([]domain.Appointment, error) {
	const q = `
		SELECT a.id, a.patient_id, a.doctor_id, a.schedule_id,
		       a.appointment_date, a.start_time_slot, a.complaint,
		       a.status, a.created_at, a.updated_at,
		       pu.name, pu.email
		FROM appointments a
		LEFT JOIN patients p ON p.id = a.patient_id
		LEFT JOIN users pu ON p.user_id = pu.id
		WHERE a.doctor_id = ?
		ORDER BY a.appointment_date ASC, a.start_time_slot ASC
	`
	rows, err := r.db.QueryContext(ctx, q, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Appointment
	for rows.Next() {
		var (
			a            domain.Appointment
			idDB         int64
			scheduleID   sql.NullInt64
			startTime    sql.NullString
			complaint    sql.NullString
			patientName  sql.NullString
			patientEmail sql.NullString
		)
		if err := rows.Scan(
			&idDB,
			&a.PatientID,
			&a.DoctorID,
			&scheduleID,
			&a.AppointmentDate,
			&startTime,
			&complaint,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
			&patientName,
			&patientEmail,
		); err != nil {
			return nil, err
		}
		a.ID = int(idDB)
		if scheduleID.Valid {
			v := int(scheduleID.Int64)
			a.ScheduleID = &v
		}
		if startTime.Valid {
			a.StartTimeSlot = startTime.String
		}
		if complaint.Valid {
			a.Complaint = complaint.String
		}
		if patientName.Valid || patientEmail.Valid {
			a.Patient = &domain.User{
				Name:  patientName.String,
				Email: patientEmail.String,
			}
		}
		result = append(result, a)
	}

	return result, nil
}
