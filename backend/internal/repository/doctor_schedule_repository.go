package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
)

type DoctorScheduleRepository interface {
	Create(ctx context.Context, schedule domain.DoctorSchedule) (domain.DoctorSchedule, error)
	GetByID(ctx context.Context, id int) (domain.DoctorSchedule, error)
	GetByDoctorID(ctx context.Context, doctorID int) ([]domain.DoctorSchedule, error)
	GetAll(ctx context.Context) ([]domain.DoctorSchedule, error)
	Update(ctx context.Context, id int, schedule domain.DoctorSchedule) (domain.DoctorSchedule, error)
	Delete(ctx context.Context, id int) error
}

type DoctorScheduleRepositoryImpl struct {
	DB *sql.DB
}

func NewDoctorScheduleRepository(db *sql.DB) DoctorScheduleRepository {
	return &DoctorScheduleRepositoryImpl{DB: db}
}

func (repo *DoctorScheduleRepositoryImpl) Create(ctx context.Context, schedule domain.DoctorSchedule) (domain.DoctorSchedule, error) {
	query := `
		INSERT INTO doctor_schedules (doctor_id, work_day, start_time, end_time, patient_quota) 
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := repo.DB.ExecContext(ctx, query,
		schedule.DoctorID, schedule.WorkDay, schedule.StartTime,
		schedule.EndTime, schedule.PatientQuota)
	if err != nil {
		log.Println("ERROR creating doctor schedule:", err)
		return schedule, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("ERROR getting last insert id:", err)
		return schedule, err
	}

	schedule.ID = int(id)
	return schedule, nil
}

func (repo *DoctorScheduleRepositoryImpl) GetByID(ctx context.Context, id int) (domain.DoctorSchedule, error) {
	query := `
		SELECT ds.id, ds.doctor_id, ds.work_day, ds.start_time, ds.end_time, 
			   ds.patient_quota, ds.created_at, ds.updated_at,
			   d.id as doctor_id, d.user_id, d.specialization_id, d.gender, d.address, 
			   d.license_number, d.is_active, d.created_at as doctor_created_at, d.updated_at as doctor_updated_at,
			   u.name as doctor_name, u.email as doctor_email, u.role as doctor_role
		FROM doctor_schedules ds
		LEFT JOIN doctors d ON ds.doctor_id = d.id
		LEFT JOIN users u ON d.user_id = u.id
		WHERE ds.id = ?
	`

	var schedule domain.DoctorSchedule
	schedule.Doctor = &domain.Doctor{}
	schedule.Doctor.User = &domain.User{}

	err := repo.DB.QueryRowContext(ctx, query, id).Scan(
		&schedule.ID, &schedule.DoctorID, &schedule.WorkDay, &schedule.StartTime,
		&schedule.EndTime, &schedule.PatientQuota, &schedule.CreatedAt, &schedule.UpdatedAt,
		&schedule.Doctor.ID, &schedule.Doctor.UserID, &schedule.Doctor.SpecializationID,
		&schedule.Doctor.Gender, &schedule.Doctor.Address, &schedule.Doctor.LicenseNumber,
		&schedule.Doctor.IsActive, &schedule.Doctor.CreatedAt, &schedule.Doctor.UpdatedAt,
		&schedule.Doctor.User.Name, &schedule.Doctor.User.Email, &schedule.Doctor.User.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schedule, errors.New("doctor schedule not found")
		}
		log.Println("ERROR getting doctor schedule by id:", err)
		return schedule, err
	}

	return schedule, nil
}

func (repo *DoctorScheduleRepositoryImpl) GetByDoctorID(ctx context.Context, doctorID int) ([]domain.DoctorSchedule, error) {
	query := `
		SELECT ds.id, ds.doctor_id, ds.work_day, ds.start_time, ds.end_time, 
			   ds.patient_quota, ds.created_at, ds.updated_at,
			   d.id as doctor_id, d.user_id, d.specialization_id, d.gender, d.address, 
			   d.license_number, d.is_active, d.created_at as doctor_created_at, d.updated_at as doctor_updated_at,
			   u.name as doctor_name, u.email as doctor_email, u.role as doctor_role
		FROM doctor_schedules ds
		LEFT JOIN doctors d ON ds.doctor_id = d.id
		LEFT JOIN users u ON d.user_id = u.id
		WHERE ds.doctor_id = ?
		ORDER BY ds.work_day, ds.start_time
	`

	rows, err := repo.DB.QueryContext(ctx, query, doctorID)
	if err != nil {
		log.Println("ERROR getting doctor schedules by doctor id:", err)
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.DoctorSchedule
	for rows.Next() {
		var schedule domain.DoctorSchedule
		schedule.Doctor = &domain.Doctor{}
		schedule.Doctor.User = &domain.User{}

		err := rows.Scan(
			&schedule.ID, &schedule.DoctorID, &schedule.WorkDay, &schedule.StartTime,
			&schedule.EndTime, &schedule.PatientQuota, &schedule.CreatedAt, &schedule.UpdatedAt,
			&schedule.Doctor.ID, &schedule.Doctor.UserID, &schedule.Doctor.SpecializationID,
			&schedule.Doctor.Gender, &schedule.Doctor.Address, &schedule.Doctor.LicenseNumber,
			&schedule.Doctor.IsActive, &schedule.Doctor.CreatedAt, &schedule.Doctor.UpdatedAt,
			&schedule.Doctor.User.Name, &schedule.Doctor.User.Email, &schedule.Doctor.User.Role,
		)

		if err != nil {
			log.Println("ERROR scanning doctor schedule row:", err)
			continue
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (repo *DoctorScheduleRepositoryImpl) GetAll(ctx context.Context) ([]domain.DoctorSchedule, error) {
	query := `
		SELECT ds.id, ds.doctor_id, ds.work_day, ds.start_time, ds.end_time, 
			   ds.patient_quota, ds.created_at, ds.updated_at,
			   d.id as doctor_id, d.user_id, d.specialization_id, d.gender, d.address, 
			   d.license_number, d.is_active, d.created_at as doctor_created_at, d.updated_at as doctor_updated_at,
			   u.name as doctor_name, u.email as doctor_email, u.role as doctor_role
		FROM doctor_schedules ds
		LEFT JOIN doctors d ON ds.doctor_id = d.id
		LEFT JOIN users u ON d.user_id = u.id
		ORDER BY ds.doctor_id, ds.work_day, ds.start_time
	`

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("ERROR getting all doctor schedules:", err)
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.DoctorSchedule
	for rows.Next() {
		var schedule domain.DoctorSchedule
		schedule.Doctor = &domain.Doctor{}
		schedule.Doctor.User = &domain.User{}

		err := rows.Scan(
			&schedule.ID, &schedule.DoctorID, &schedule.WorkDay, &schedule.StartTime,
			&schedule.EndTime, &schedule.PatientQuota, &schedule.CreatedAt, &schedule.UpdatedAt,
			&schedule.Doctor.ID, &schedule.Doctor.UserID, &schedule.Doctor.SpecializationID,
			&schedule.Doctor.Gender, &schedule.Doctor.Address, &schedule.Doctor.LicenseNumber,
			&schedule.Doctor.IsActive, &schedule.Doctor.CreatedAt, &schedule.Doctor.UpdatedAt,
			&schedule.Doctor.User.Name, &schedule.Doctor.User.Email, &schedule.Doctor.User.Role,
		)

		if err != nil {
			log.Println("ERROR scanning doctor schedule row:", err)
			continue
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (repo *DoctorScheduleRepositoryImpl) Update(ctx context.Context, id int, schedule domain.DoctorSchedule) (domain.DoctorSchedule, error) {
	query := `
		UPDATE doctor_schedules 
		SET doctor_id = ?, work_day = ?, start_time = ?, end_time = ?, patient_quota = ?
		WHERE id = ?
	`

	_, err := repo.DB.ExecContext(ctx, query,
		schedule.DoctorID, schedule.WorkDay, schedule.StartTime,
		schedule.EndTime, schedule.PatientQuota, id)
	if err != nil {
		log.Println("ERROR updating doctor schedule:", err)
		return schedule, err
	}

	schedule.ID = id
	return schedule, nil
}

func (repo *DoctorScheduleRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM doctor_schedules WHERE id = ?`

	_, err := repo.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("ERROR deleting doctor schedule:", err)
		return err
	}

	return nil
}
