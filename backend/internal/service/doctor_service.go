package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type DoctorService interface {
	CreateDoctor(ctx context.Context, req domain.DoctorRequest) (domain.Doctor, error)
	GetDoctorByID(ctx context.Context, id int) (domain.Doctor, error)
	GetByUserID(ctx context.Context, userID int) (domain.Doctor, error)
	GetAllDoctors(ctx context.Context) ([]domain.Doctor, error)
	UpdateDoctorByDoctorID(ctx context.Context, id int, req domain.DoctorRequest) (domain.Doctor, error)
	UpdateDoctorByUserID(ctx context.Context, userID int, req domain.DoctorRequest) (domain.Doctor, error)
	DeleteDoctor(ctx context.Context, id int) error
	SearchDoctors(ctx context.Context, keyword string, specializationID int) ([]domain.Doctor, error)
}

type DoctorServiceImpl struct {
	DoctorRepo repository.DoctorRepository
	UserRepo   repository.UserRepository
	DB         *sql.DB
}

func NewDoctorService(doctorRepo repository.DoctorRepository, userRepo repository.UserRepository, db *sql.DB) DoctorService {
	return &DoctorServiceImpl{
		DoctorRepo: doctorRepo,
		UserRepo:   userRepo,
		DB:         db,
	}
}

func (s *DoctorServiceImpl) CreateDoctor(ctx context.Context, req domain.DoctorRequest) (domain.Doctor, error) {
	// Validate gender
	if !domain.IsValidGender(req.Gender) {
		return domain.Doctor{}, errors.New("invalid gender. Valid values: male, female")
	}

	// Create user first with doctor role
	user := domain.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       "password",
		Role:           domain.RoleDoctor,
		ProfilePicture: "",
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("ERROR hashing password:", err)
		return domain.Doctor{}, err
	}
	user.Password = string(hashedPassword)

	// Create doctor with created user ID using same transaction
	doctor := domain.Doctor{
		UserID:           0, // Will be set in repository
		SpecializationID: req.SpecializationID,
		Gender:           domain.Gender(req.Gender),
		Address:          req.Address,
		LicenseNumber:    req.LicenseNumber,
		IsActive:         true,
	}

	// Create both user and doctor in repository with transaction
	createdDoctor, err := s.DoctorRepo.CreateWithUser(ctx, user, doctor)
	if err != nil {
		return domain.Doctor{}, err
	}

	return createdDoctor, nil
}

func (s *DoctorServiceImpl) GetDoctorByID(ctx context.Context, id int) (domain.Doctor, error) {
	doctor, err := s.DoctorRepo.GetByDoctorID(ctx, id)
	if err != nil {
		return doctor, err
	}

	// Ensure User is initialized
	if doctor.User == nil {
		doctor.User = &domain.User{}
	}

	return doctor, nil
}

func (s *DoctorServiceImpl) GetByUserID(ctx context.Context, userID int) (domain.Doctor, error) {
	return s.DoctorRepo.GetByUserId(ctx, userID)
}

func (s *DoctorServiceImpl) GetAllDoctors(ctx context.Context) ([]domain.Doctor, error) {
	doctors, err := s.DoctorRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure User is initialized for each doctor
	for i := range doctors {
		if doctors[i].User == nil {
			doctors[i].User = &domain.User{}
		}
	}

	return doctors, nil
}

func (s *DoctorServiceImpl) UpdateDoctorByDoctorID(ctx context.Context, id int, req domain.DoctorRequest) (domain.Doctor, error) {
	// Get existing doctor by doctor ID
	existingDoctor, err := s.DoctorRepo.GetByDoctorID(ctx, id)
	if err != nil {
		return domain.Doctor{}, err
	}

	// Get existing user using doctor's user ID
	existingUser, err := s.UserRepo.FindByID(ctx, existingDoctor.UserID)
	if err != nil {
		log.Println("ERROR finding user:", err)
		return domain.Doctor{}, err
	}

	// Update user fields if provided
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" {
		// Check if email already exists for another user
		existingUserByEmail, err := s.UserRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingUserByEmail.ID != existingUser.ID {
			return domain.Doctor{}, errors.New("email already exists")
		}
		existingUser.Email = req.Email
	}

	// Update doctor fields if provided
	if req.SpecializationID != 0 {
		existingDoctor.SpecializationID = req.SpecializationID
	}
	if req.Gender != "" {
		if !domain.IsValidGender(req.Gender) {
			return domain.Doctor{}, errors.New("invalid gender. Valid values: male, female")
		}
		existingDoctor.Gender = domain.Gender(req.Gender)
	}
	if req.Address != "" {
		existingDoctor.Address = req.Address
	}
	if req.LicenseNumber != "" {
		existingDoctor.LicenseNumber = req.LicenseNumber
	}

	// Use UpdateWithUser to update both user and doctor in one transaction
	updatedDoctor, err := s.DoctorRepo.UpdateWithUser(ctx, existingUser, existingDoctor)
	if err != nil {
		return domain.Doctor{}, err
	}

	return updatedDoctor, nil
}

func (s *DoctorServiceImpl) UpdateDoctorByUserID(ctx context.Context, userID int, req domain.DoctorRequest) (domain.Doctor, error) {
	// Get existing doctor by user ID
	existingDoctor, err := s.DoctorRepo.GetByUserId(ctx, userID)
	if err != nil {
		return domain.Doctor{}, err
	}

	// Get existing user
	existingUser, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		log.Println("ERROR finding user:", err)
		return domain.Doctor{}, err
	}

	// Update user fields if provided
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" {
		// Check if email already exists for another user
		existingUserByEmail, err := s.UserRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingUserByEmail.ID != userID {
			return domain.Doctor{}, errors.New("email already exists")
		}
		existingUser.Email = req.Email
	}

	// Update doctor fields if provided
	if req.SpecializationID != 0 {
		existingDoctor.SpecializationID = req.SpecializationID
	}
	if req.Gender != "" {
		if !domain.IsValidGender(req.Gender) {
			return domain.Doctor{}, errors.New("invalid gender. Valid values: male, female")
		}
		existingDoctor.Gender = domain.Gender(req.Gender)
	}
	if req.Address != "" {
		existingDoctor.Address = req.Address
	}
	if req.LicenseNumber != "" {
		existingDoctor.LicenseNumber = req.LicenseNumber
	}

	// Update both user and doctor in repository with transaction
	updatedDoctor, err := s.DoctorRepo.UpdateWithUser(ctx, existingUser, existingDoctor)
	if err != nil {
		return domain.Doctor{}, err
	}

	return updatedDoctor, nil
}

func (s *DoctorServiceImpl) DeleteDoctor(ctx context.Context, id int) error {
	return s.DoctorRepo.Delete(ctx, id)
}

func (s *DoctorServiceImpl) SearchDoctors(ctx context.Context, keyword string, specializationID int) ([]domain.Doctor, error) {
	// trimming sederhana
	keyword = strings.TrimSpace(keyword)
	return s.DoctorRepo.Search(ctx, keyword, specializationID)
}
