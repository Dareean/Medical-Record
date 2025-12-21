package storage

import (
	"database/sql"
	"log"
	"os" // Tambahkan ini untuk membaca "kabel" Railway
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetConnection() (*sql.DB, error) {
	// 1. Ambil "kabel" koneksi dari dashboard Railway
	dsn := os.Getenv("MYSQL_URL")

	// 2. Jika kabel kosong (berarti kamu lagi jalanin di laptop sendiri)
	if dsn == "" {
		dsn = "root:@tcp(localhost:3306)/medical_record_db?parseTime=true"
		log.Println("Database: Menggunakan koneksi lokal (XAMPP)")
	} else {
		log.Println("Database: Berhasil membaca alamat dari Railway")
	}

	// 3. Gunakan variabel dsn (bukan tulisan manual lagi)
	db, err := sql.Open("mysql", dsn)
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
	if err := RunMigrations(db); err != nil {
		return err
	}

	if err := SeedAllData(db); err != nil {
		log.Println("WARNING: Failed to seed data:", err)
	}

	return nil
}