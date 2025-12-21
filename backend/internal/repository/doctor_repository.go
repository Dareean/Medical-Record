package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
)

type DoctorRepository interface {
	CreateWithUser(ctx context.Context, user domain.User, doctor domain.Doctor) (domain.Doctor, error)
	GetAll(ctx context.Context) ([]domain.Doctor, error)
	GetByDoctorID(ctx context.Context, doctorID int) (domain.Doctor, error)
	GetByUserId(ctx context.Context, userID int) (domain.Doctor, error)
	UpdateWithUser(ctx context.Context, user domain.User, doctor domain.Doctor) (domain.Doctor, error)
	Delete(ctx context.Context, id int) error

	Search(ctx context.Context, keyword string, specializationID int) ([]domain.Doctor, error)
}

type DoctorRepositoryImpl struct {
	DB *sql.DB
}

func NewDoctorRepository(db *sql.DB) DoctorRepository {
	return &DoctorRepositoryImpl{DB: db}
}

func (repo *DoctorRepositoryImpl) CreateWithUser(ctx context.Context, user domain.User, doctor domain.Doctor) (domain.Doctor, error) {
	// Start transaction
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("ERROR starting transaction:", err)
		return domain.Doctor{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	// Check if email already exists
	var existingUserID int
	checkEmailQuery := `SELECT id FROM users WHERE email = ?`
	err = tx.QueryRowContext(ctx, checkEmailQuery, user.Email).Scan(&existingUserID)
	if err == nil {
		tx.Rollback()
		return domain.Doctor{}, errors.New("email already exists")
	}
	if err != sql.ErrNoRows {
		log.Println("ERROR checking email:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}

	// Insert user first
	insertUserQuery := `
		INSERT INTO users (name, email, password, role, profile_picture) 
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, insertUserQuery, user.Name, user.Email, user.Password, user.Role, user.ProfilePicture)
	if err != nil {
		log.Println("ERROR creating user:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}

	// Check if user was actually inserted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("ERROR getting rows affected for user insert:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return domain.Doctor{}, errors.New("failed to create user")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Println("ERROR getting user last insert id:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}

	// Insert doctor with created user ID
	doctor.UserID = int(userID)
	insertDoctorQuery := `
		INSERT INTO doctors (user_id, specialization_id, gender, address, license_number, is_active) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	doctorResult, err := tx.ExecContext(ctx, insertDoctorQuery,
		doctor.UserID, doctor.SpecializationID, doctor.Gender,
		doctor.Address, doctor.LicenseNumber, doctor.IsActive)
	if err != nil {
		log.Println("ERROR creating doctor:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}

	// Check if doctor was actually inserted
	doctorRowsAffected, err := doctorResult.RowsAffected()
	if err != nil {
		log.Println("ERROR getting rows affected for doctor insert:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}
	if doctorRowsAffected == 0 {
		tx.Rollback()
		return domain.Doctor{}, errors.New("failed to create doctor")
	}

	doctorID, err := doctorResult.LastInsertId()
	if err != nil {
		log.Println("ERROR getting doctor last insert id:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Println("ERROR committing transaction:", err)
		return domain.Doctor{}, err
	}

	// Get the created doctor with user data from database
	createdDoctor, err := repo.GetByDoctorID(ctx, int(doctorID))
	if err != nil {
		log.Println("ERROR retrieving created doctor:", err)
		return domain.Doctor{}, err
	}

	return createdDoctor, nil
}

func (repo *DoctorRepositoryImpl) GetAll(ctx context.Context) ([]domain.Doctor, error) {
	query := `
			SELECT d.id, d.user_id, d.specialization_id, d.gender, d.address, 
				   d.license_number, d.is_active, d.created_at, d.updated_at,
				   u.name, u.email, u.role, u.profile_picture,
				   s.name as specialization_name
			FROM doctors d
			LEFT JOIN users u ON d.user_id = u.id
			LEFT JOIN specializations s ON d.specialization_id = s.id
			ORDER BY d.created_at DESC
		`

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("ERROR getting all doctors:", err)
		return nil, err
	}
	defer rows.Close()

	var doctors []domain.Doctor
	for rows.Next() {
		var doctor domain.Doctor
		doctor.User = &domain.User{}
		var specializationName sql.NullString
		var profilePicture sql.NullString

		err := rows.Scan(
			&doctor.ID, &doctor.UserID, &doctor.SpecializationID, &doctor.Gender,
			&doctor.Address, &doctor.LicenseNumber, &doctor.IsActive,
			&doctor.CreatedAt, &doctor.UpdatedAt,
			&doctor.User.Name, &doctor.User.Email, &doctor.User.Role,
			&profilePicture, &specializationName,
		)

		if err != nil {
			log.Println("ERROR scanning doctor row:", err)
			continue
		}

		// Handle nullable profile_picture
		if profilePicture.Valid {
			doctor.User.ProfilePicture = profilePicture.String
		}

		if specializationName.Valid {
			doctor.Specialization = &domain.Specialization{
				Name: specializationName.String,
			}
		}

		doctors = append(doctors, doctor)
	}

	return doctors, nil
}

func (repo *DoctorRepositoryImpl) GetByDoctorID(ctx context.Context, id int) (domain.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, d.specialization_id, d.gender, d.address, 
			   d.license_number, d.is_active, d.created_at, d.updated_at,
			   u.name, u.email, u.role, u.profile_picture,
			   s.name as specialization_name
		FROM doctors d
		LEFT JOIN users u ON d.user_id = u.id
		LEFT JOIN specializations s ON d.specialization_id = s.id
		WHERE d.id = ?
	`

	var doctor domain.Doctor
	doctor.User = &domain.User{}
	var specializationName sql.NullString
	var profilePicture sql.NullString

	err := repo.DB.QueryRowContext(ctx, query, id).Scan(
		&doctor.ID, &doctor.UserID, &doctor.SpecializationID, &doctor.Gender,
		&doctor.Address, &doctor.LicenseNumber, &doctor.IsActive,
		&doctor.CreatedAt, &doctor.UpdatedAt,
		&doctor.User.Name, &doctor.User.Email, &doctor.User.Role,
		&profilePicture, &specializationName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return doctor, errors.New("doctor not found")
		}
		log.Println("ERROR getting doctor by id:", err)
		return doctor, err
	}

	// Handle nullable profile_picture
	if profilePicture.Valid {
		doctor.User.ProfilePicture = profilePicture.String
	}

	if specializationName.Valid {
		doctor.Specialization = &domain.Specialization{
			Name: specializationName.String,
		}
	}

	return doctor, nil
}

func (repo *DoctorRepositoryImpl) GetByUserId(ctx context.Context, userID int) (domain.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, d.specialization_id, d.gender, d.address, 
			   d.license_number, d.is_active, d.created_at, d.updated_at,
			   u.name, u.email, u.role, u.profile_picture,
			   s.name as specialization_name
		FROM doctors d
		LEFT JOIN users u ON d.user_id = u.id
		LEFT JOIN specializations s ON d.specialization_id = s.id
		WHERE d.user_id = ?
	`

	var doctor domain.Doctor
	doctor.User = &domain.User{}
	var specializationName sql.NullString
	var profilePicture sql.NullString

	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(
		&doctor.ID, &doctor.UserID, &doctor.SpecializationID, &doctor.Gender,
		&doctor.Address, &doctor.LicenseNumber, &doctor.IsActive,
		&doctor.CreatedAt, &doctor.UpdatedAt,
		&doctor.User.Name, &doctor.User.Email, &doctor.User.Role,
		&profilePicture, &specializationName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return doctor, errors.New("doctor not found")
		}
		log.Println("ERROR getting doctor by user id:", err)
		return doctor, err
	}

	// Handle nullable profile_picture
	if profilePicture.Valid {
		doctor.User.ProfilePicture = profilePicture.String
	}

	if specializationName.Valid {
		doctor.Specialization = &domain.Specialization{
			Name: specializationName.String,
		}
	}

	return doctor, nil
}

func (repo *DoctorRepositoryImpl) UpdateWithUser(ctx context.Context, user domain.User, doctor domain.Doctor) (domain.Doctor, error) {
	// Start transaction
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("ERROR starting transaction:", err)
		return domain.Doctor{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	// Check if email already exists for another user
	if user.Email != "" {
		var existingUserID int
		checkEmailQuery := `SELECT id FROM users WHERE email = ? AND id != ?`
		err = tx.QueryRowContext(ctx, checkEmailQuery, user.Email, user.ID).Scan(&existingUserID)
		if err == nil {
			tx.Rollback()
			return domain.Doctor{}, errors.New("email already exists")
		}
		if err != sql.ErrNoRows {
			log.Println("ERROR checking email:", err)
			tx.Rollback()
			return domain.Doctor{}, err
		}
	}

	// Update user
	updateUserQuery := `
		UPDATE users 
		SET name = COALESCE(?, name), 
			email = COALESCE(?, email), 
			profile_picture = COALESCE(?, profile_picture),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	userResult, err := tx.ExecContext(ctx, updateUserQuery, user.Name, user.Email, user.ProfilePicture, user.ID)
	if err != nil {
		log.Println("ERROR updating user:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}

	// Check if user was actually updated
	userRowsAffected, err := userResult.RowsAffected()
	if err != nil {
		log.Println("ERROR getting rows affected for user update:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}
	if userRowsAffected == 0 {
		tx.Rollback()
		return domain.Doctor{}, errors.New("user not found or no changes made")
	}

	// Update doctor
	updateDoctorQuery := `
		UPDATE doctors 
		SET specialization_id = COALESCE(?, specialization_id), 
			gender = COALESCE(?, gender), 
			address = COALESCE(?, address), 
			license_number = COALESCE(?, license_number),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	doctorResult, err := tx.ExecContext(ctx, updateDoctorQuery,
		doctor.SpecializationID, doctor.Gender, doctor.Address, doctor.LicenseNumber, doctor.ID)
	if err != nil {
		log.Println("ERROR updating doctor:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}

	// Check if doctor was actually updated
	doctorRowsAffected, err := doctorResult.RowsAffected()
	if err != nil {
		log.Println("ERROR getting rows affected for doctor update:", err)
		tx.Rollback()
		return domain.Doctor{}, err
	}
	if doctorRowsAffected == 0 {
		tx.Rollback()
		return domain.Doctor{}, errors.New("doctor not found or no changes made")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Println("ERROR committing transaction:", err)
		return domain.Doctor{}, err
	}

	// Get the updated doctor with user data from database
	updatedDoctor, err := repo.GetByDoctorID(ctx, doctor.ID)
	if err != nil {
		log.Println("ERROR retrieving updated doctor:", err)
		return domain.Doctor{}, err
	}

	return updatedDoctor, nil
}

func (repo *DoctorRepositoryImpl) Delete(ctx context.Context, id int) error {
	// First check if doctor exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM doctors WHERE id = ?)`
	err := repo.DB.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		log.Println("ERROR checking doctor existence:", err)
		return err
	}

	if !exists {
		return errors.New("doctor not found")
	}

	// Delete the doctor
	query := `DELETE FROM doctors WHERE id = ?`

	result, err := repo.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("ERROR deleting doctor:", err)
		return err
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("ERROR getting rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("doctor not found")
	}

	return nil
}

func (r *DoctorRepositoryImpl) Search(ctx context.Context, keyword string, specializationID int) ([]domain.Doctor, error) {

	query := `
		SELECT 
			d.id, d.user_id, d.specialization_id, d.gender, d.address,
			d.license_number, d.is_active,
			u.id, u.name, u.email
		FROM doctors d
		JOIN users u ON u.id = d.user_id
		WHERE u.name LIKE ?
	`

	args := []any{"%" + keyword + "%"}

	if specializationID != 0 {
		query += " AND d.specialization_id = ?"
		args = append(args, specializationID)
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []domain.Doctor

	for rows.Next() {
		var d domain.Doctor
		var u domain.User

		err := rows.Scan(
			&d.ID,
			&d.UserID,
			&d.SpecializationID,
			&d.Gender,
			&d.Address,
			&d.LicenseNumber,
			&d.IsActive,
			&u.ID,
			&u.Name,
			&u.Email,
		)
		if err != nil {
			return nil, err
		}

		d.User = &u
		doctors = append(doctors, d)
	}

	return doctors, nil
}

