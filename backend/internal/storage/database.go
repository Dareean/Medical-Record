package storage

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/medical_record_db?parseTime=true")
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db, nil
}

func InitializeDatabase(db *sql.DB) error {
	// Run migrations to create tables
	if err := RunMigrations(db); err != nil {
		return err
	}

	// Seed all data
	if err := SeedAllData(db); err != nil {
		log.Println("WARNING: Failed to seed data:", err)
	}

	return nil
}
