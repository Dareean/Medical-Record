package service

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user domain.RegisterRequest) (domain.User, error)
	Login(ctx context.Context, credentials domain.LoginRequest) (domain.LoginResponse, error)
	GenerateToken(user domain.User) (string, error)
}

type UserServiceImpl struct {
	Repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &UserServiceImpl{
		Repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, req domain.RegisterRequest) (domain.User, error) {
	user := domain.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		Role:           req.Role,
		ProfilePicture: req.ProfilePicture,
	}

	// Set default role jika kosong
	if user.Role == "" {
		user.Role = domain.RolePatient
	}

	// Validasi role
	if !domain.IsValidRole(string(user.Role)) {
		log.Println("ERROR: invalid role:", user.Role)
		return user, errors.New("invalid role. Valid roles: admin, doctor, patient")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("ERROR:", err)
		return user, err
	}
	user.Password = string(hashedPassword)
	return service.Repo.Register(ctx, user)
}

func (service *UserServiceImpl) Login(ctx context.Context, credentials domain.LoginRequest) (domain.LoginResponse, error) {
	if credentials.Email == "" || credentials.Password == "" {
		return domain.LoginResponse{}, errors.New("email and password are required")
	}

	user, err := service.Repo.FindByEmail(ctx, credentials.Email)
	if err != nil {
		log.Println("ERROR: user not found:", err)
		return domain.LoginResponse{}, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		log.Println("ERROR: password mismatch:", err)
		return domain.LoginResponse{}, errors.New("password wrong")
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		log.Println("ERROR: failed to generate token:", err)
		return domain.LoginResponse{}, errors.New("failed to generate authentication token")
	}

	// Log berhasil login
	log.Printf("User %s logged in successfully", user.Email)

	response := domain.LoginResponse{
		User: domain.User{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return response, nil
}

func (service *UserServiceImpl) GenerateToken(user domain.User) (string, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT secret not configured")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
