package service

import (
	"context"
	"errors"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/repository"
)

type DoctorScheduleService interface {
	CreateSchedule(ctx context.Context, req domain.DoctorScheduleRequest) (domain.DoctorSchedule, error)
	GetScheduleByID(ctx context.Context, id int) (domain.DoctorSchedule, error)
	GetSchedulesByDoctorID(ctx context.Context, doctorID int) ([]domain.DoctorSchedule, error)
	GetAllSchedules(ctx context.Context) ([]domain.DoctorSchedule, error)
	UpdateSchedule(ctx context.Context, id int, req domain.DoctorScheduleRequest) (domain.DoctorSchedule, error)
	DeleteSchedule(ctx context.Context, id int) error
}

type DoctorScheduleServiceImpl struct {
	ScheduleRepo repository.DoctorScheduleRepository
	DoctorRepo   repository.DoctorRepository
}

func NewDoctorScheduleService(scheduleRepo repository.DoctorScheduleRepository, doctorRepo repository.DoctorRepository) DoctorScheduleService {
	return &DoctorScheduleServiceImpl{
		ScheduleRepo: scheduleRepo,
		DoctorRepo:   doctorRepo,
	}
}

func (s *DoctorScheduleServiceImpl) CreateSchedule(ctx context.Context, req domain.DoctorScheduleRequest) (domain.DoctorSchedule, error) {
	// Validate work day
	if !domain.IsValidWorkDay(req.WorkDay) {
		return domain.DoctorSchedule{}, errors.New("invalid work day. Valid values: monday, tuesday, wednesday, thursday, friday, saturday, sunday")
	}

	// Check if doctor exists
	_, err := s.DoctorRepo.GetByDoctorID(ctx, req.DoctorID)
	if err != nil {
		return domain.DoctorSchedule{}, errors.New("doctor not found")
	}

	// Create schedule
	schedule := domain.DoctorSchedule{
		DoctorID:     req.DoctorID,
		WorkDay:      domain.WorkDay(req.WorkDay),
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		PatientQuota: req.PatientQuota,
	}

	createdSchedule, err := s.ScheduleRepo.Create(ctx, schedule)
	if err != nil {
		return domain.DoctorSchedule{}, err
	}

	return createdSchedule, nil
}

func (s *DoctorScheduleServiceImpl) GetScheduleByID(ctx context.Context, id int) (domain.DoctorSchedule, error) {
	schedule, err := s.ScheduleRepo.GetByID(ctx, id)
	if err != nil {
		return schedule, err
	}

	// Ensure Doctor is initialized
	if schedule.Doctor == nil {
		schedule.Doctor = &domain.Doctor{}
	}
	if schedule.Doctor.User == nil {
		schedule.Doctor.User = &domain.User{}
	}

	return schedule, nil
}

func (s *DoctorScheduleServiceImpl) GetSchedulesByDoctorID(ctx context.Context, doctorID int) ([]domain.DoctorSchedule, error) {
	// Check if doctor exists
	_, err := s.DoctorRepo.GetByDoctorID(ctx, doctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}

	schedules, err := s.ScheduleRepo.GetByDoctorID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	// Ensure Doctor is initialized for each schedule
	for i := range schedules {
		if schedules[i].Doctor == nil {
			schedules[i].Doctor = &domain.Doctor{}
		}
		if schedules[i].Doctor.User == nil {
			schedules[i].Doctor.User = &domain.User{}
		}
	}

	return schedules, nil
}

func (s *DoctorScheduleServiceImpl) GetAllSchedules(ctx context.Context) ([]domain.DoctorSchedule, error) {
	schedules, err := s.ScheduleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure Doctor is initialized for each schedule
	for i := range schedules {
		if schedules[i].Doctor == nil {
			schedules[i].Doctor = &domain.Doctor{}
		}
		if schedules[i].Doctor.User == nil {
			schedules[i].Doctor.User = &domain.User{}
		}
	}

	return schedules, nil
}

func (s *DoctorScheduleServiceImpl) UpdateSchedule(ctx context.Context, id int, req domain.DoctorScheduleRequest) (domain.DoctorSchedule, error) {
	// Validate work day
	if !domain.IsValidWorkDay(req.WorkDay) {
		return domain.DoctorSchedule{}, errors.New("invalid work day. Valid values: monday, tuesday, wednesday, thursday, friday, saturday, sunday")
	}

	// Check if schedule exists
	_, err := s.ScheduleRepo.GetByID(ctx, id)
	if err != nil {
		return domain.DoctorSchedule{}, errors.New("schedule not found")
	}

	// Check if doctor exists
	_, err = s.DoctorRepo.GetByDoctorID(ctx, req.DoctorID)
	if err != nil {
		return domain.DoctorSchedule{}, errors.New("doctor not found")
	}

	// Update schedule
	schedule := domain.DoctorSchedule{
		DoctorID:     req.DoctorID,
		WorkDay:      domain.WorkDay(req.WorkDay),
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		PatientQuota: req.PatientQuota,
	}

	updatedSchedule, err := s.ScheduleRepo.Update(ctx, id, schedule)
	if err != nil {
		return domain.DoctorSchedule{}, err
	}

	return updatedSchedule, nil
}

func (s *DoctorScheduleServiceImpl) DeleteSchedule(ctx context.Context, id int) error {
	// Check if schedule exists
	_, err := s.ScheduleRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("schedule not found")
	}

	return s.ScheduleRepo.Delete(ctx, id)
}
