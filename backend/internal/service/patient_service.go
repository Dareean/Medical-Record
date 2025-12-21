package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/repository"
)

var (
	ErrNotAllowed    = errors.New("aksi tidak diizinkan")
	ErrInvalidStatus = errors.New("status appointment tidak valid")
)

type PatientService interface {
	CreateAppointment(ctx context.Context, userID int64, doctorID int64, appointmentDate time.Time, startTimeSlot string, complaint string, scheduleID *int64) (*domain.Appointment, error)
	CancelAppointment(ctx context.Context, userID, appointmentID int64) error
	GetAppointmentHistory(ctx context.Context, userID int64) ([]domain.Appointment, error)
	GetAppointmentDetail(ctx context.Context, id int64) (*domain.Appointment, error)
	GetDoctorAppointments(ctx context.Context, doctorID int64) ([]domain.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, doctorID, appointmentID int64, status domain.AppointmentStatus) error
}

type patientService struct {
	db              *sql.DB
	appointmentRepo repository.AppointmentRepository
	patientRepo     repository.PatientRepository
}

func NewPatientService(
	db *sql.DB,
	ar repository.AppointmentRepository,
	pr repository.PatientRepository,
) PatientService {
	return &patientService{
		db:              db,
		appointmentRepo: ar,
		patientRepo:     pr,
	}
}

func (s *patientService) ensurePatient(ctx context.Context, userID int64) (*domain.Patient, error) {
	patient, err := s.patientRepo.GetByUserID(ctx, userID)
	if err == nil {
		return patient, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return s.patientRepo.CreateForUser(ctx, userID)
	}
	return nil, err
}

func (s *patientService) CreateAppointment(
	ctx context.Context,
	userID int64,
	doctorID int64,
	appointmentDate time.Time,
	startTimeSlot string,
	complaint string,
	scheduleID *int64,
) (ap *domain.Appointment, err error) {
	patient, err := s.ensurePatient(ctx, userID)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var schedulePtr *int
	if scheduleID != nil {
		v := int(*scheduleID)
		schedulePtr = &v
	}
	patientEntityID := patient.ID
	ap = &domain.Appointment{
		PatientID:       patientEntityID,
		DoctorID:        int(doctorID),
		ScheduleID:      schedulePtr,
		AppointmentDate: appointmentDate,
		StartTimeSlot:   startTimeSlot,
		Complaint:       complaint,
		Status:          domain.AppointmentStatusPending,
	}

	if err = s.appointmentRepo.CreateTx(ctx, tx, ap); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ap, nil
}

func (s *patientService) CancelAppointment(ctx context.Context, userID, appointmentID int64) (err error) {
	patient, err := s.ensurePatient(ctx, userID)
	if err != nil {
		return err
	}

	ap, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}
	if ap.PatientID != patient.ID {
		return ErrNotAllowed
	}
	if ap.Status != domain.AppointmentStatusPending && ap.Status != domain.AppointmentStatusConfirmed {
		return ErrInvalidStatus
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = s.appointmentRepo.UpdateStatusTx(ctx, tx, appointmentID, domain.AppointmentStatusRejected); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *patientService) GetAppointmentHistory(ctx context.Context, userID int64) ([]domain.Appointment, error) {
	patient, err := s.ensurePatient(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.appointmentRepo.GetByPatient(ctx, int64(patient.ID))
}

func (s *patientService) GetAppointmentDetail(ctx context.Context, id int64) (*domain.Appointment, error) {
	return s.appointmentRepo.GetByID(ctx, id)
}

func (s *patientService) GetDoctorAppointments(ctx context.Context, doctorID int64) ([]domain.Appointment, error) {
	return s.appointmentRepo.GetByDoctor(ctx, doctorID)
}

func (s *patientService) UpdateAppointmentStatus(ctx context.Context, doctorID, appointmentID int64, status domain.AppointmentStatus) error {
	if !domain.IsValidAppointmentStatus(string(status)) {
		return ErrInvalidStatus
	}

	ap, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}
	if ap.DoctorID != int(doctorID) {
		return ErrNotAllowed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = s.appointmentRepo.UpdateStatusTx(ctx, tx, appointmentID, status); err != nil {
		return err
	}
	return tx.Commit()
}
