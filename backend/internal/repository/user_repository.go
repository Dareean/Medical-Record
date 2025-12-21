package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
)

type UserRepository interface {
	Register(ctx context.Context, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, id int) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (repo *UserRepositoryImpl) Register(ctx context.Context, user domain.User) (domain.User, error) {
	var cekUser domain.User
	query := `
		SELECT id, name, email 
		FROM users 
		WHERE email = ?
	`

	err := repo.DB.QueryRowContext(ctx, query, user.Email).Scan(&cekUser.ID, &cekUser.Name, &cekUser.Email)

	// Error: email sudah ada
	if err == nil {
		log.Println("ERROR: email sudah ada")
		return user, errors.New("email sudah ada")
	}

	// Error: sql
	if err != sql.ErrNoRows {
		log.Println("ERROR:", err)
		return user, err
	}

	// Insert user baru
	result, err := repo.DB.ExecContext(ctx, "INSERT INTO users (name, email, password, role, profile_picture) VALUES (?, ?, ?, ?, ?)", user.Name, user.Email, user.Password, user.Role, user.ProfilePicture)
	if err != nil {
		log.Println("ERROR:", err)
		return user, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("ERROR:", err)
		return user, err
	}
	user.ID = int(id)

	return user, nil
}

func (repo *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user domain.User
	var profilePicture sql.NullString

	err := repo.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		log.Println("ERROR FindByEmail:", err)
		return user, err
	}

	// Handle nullable profile_picture
	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	}

	return user, nil
}

func (repo *UserRepositoryImpl) FindByID(ctx context.Context, id int) (domain.User, error) {
	query := `
		SELECT id, email, password, name, role, profile_picture, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user domain.User
	var profilePicture sql.NullString

	err := repo.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&profilePicture,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		log.Println("ERROR FindByID:", err)
		return user, err
	}

	// Handle nullable profile_picture
	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	}

	return user, nil
}

func (repo *UserRepositoryImpl) Update(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		UPDATE users 
		SET name = ?, email = ?, role = ?, profile_picture = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	// Check if transaction is available in context
	var db interface{} = repo.DB
	if tx, ok := ctx.Value("tx").(*sql.Tx); ok {
		db = tx
	}

	var result sql.Result
	var err error

	switch db := db.(type) {
	case *sql.DB:
		result, err = db.ExecContext(ctx, query, user.Name, user.Email, user.Role, user.ProfilePicture, user.ID)
	case *sql.Tx:
		result, err = db.ExecContext(ctx, query, user.Name, user.Email, user.Role, user.ProfilePicture, user.ID)
	default:
		return user, errors.New("invalid database connection")
	}

	if err != nil {
		log.Println("ERROR Update user:", err)
		return user, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("ERROR getting rows affected:", err)
		return user, err
	}

	if rowsAffected == 0 {
		return user, errors.New("user not found or no changes made")
	}

	return user, nil
}
