package main

import (
	"log"
	"net/http"
	"os"

	authMiddleware "github.com/JinXVIII/BE-Medical-Record/internal/middleware"

	"github.com/JinXVIII/BE-Medical-Record/internal/handler"
	"github.com/JinXVIII/BE-Medical-Record/internal/repository"
	"github.com/JinXVIII/BE-Medical-Record/internal/service"
	"github.com/JinXVIII/BE-Medical-Record/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Could not load .env file")
	}

	// Cek koneksi database
	db, err := storage.GetConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}
	log.Println("Berhasil terhubung ke database")

	// Inisialisasi database (buat tabel jika belum ada)
	if err := storage.InitializeDatabase(db); err != nil {
		log.Fatal("Gagal menginisialisasi database: ", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET tidak bisa kosong")
	}

	// User Auth
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtSecret)
	userHandler := handler.NewUserHandler(userService)

	// Doctor Management
	doctorRepo := repository.NewDoctorRepository(db)
	doctorService := service.NewDoctorService(doctorRepo, userRepo, db)
	doctorHandler := handler.NewDoctorHandler(doctorService)

	// Doctor Profile Management
	doctorProfileHandler := handler.NewDoctorProfileHandler(doctorService)

	// Doctor Schedule Management
	scheduleRepo := repository.NewDoctorScheduleRepository(db)
	scheduleService := service.NewDoctorScheduleService(scheduleRepo, doctorRepo)
	scheduleHandler := handler.NewDoctorScheduleHandler(scheduleService, doctorService)
	doctorScheduleHandler := handler.NewDoctorScheduleHandler(scheduleService, doctorService)

	// Appointment
	appoinmentRepo := repository.NewAppointmentRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	patientService := service.NewPatientService(db, appoinmentRepo, patientRepo)
	patientHandler := handler.NewPatientHandler(patientService)
	doctorAppointmentHandler := handler.NewDoctorAppointmentHandler(patientService, doctorService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Route("/api", func(r chi.Router) {
		// User endpoints
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.AuthMiddleware)

			r.Route("/doctor", func(r chi.Router) {
				// Doctor Profile
				r.Route("/profile", func(r chi.Router) {
					r.Get("/", doctorProfileHandler.GetMyProfile)    // Get my profile
					r.Put("/", doctorProfileHandler.UpdateMyProfile) // Update my profile
				})

				// Doctor Schedules (doctors only)
				r.Route("/schedules", func(r chi.Router) {
					r.Get("/", scheduleHandler.GetMySchedules)        // Get my schedules
					r.Post("/", scheduleHandler.CreateSchedule)       // Create schedule
					r.Put("/{id}", scheduleHandler.UpdateSchedule)    // Update schedule
					r.Delete("/{id}", scheduleHandler.DeleteSchedule) // Delete schedule
				})

				r.Route("/appointments", func(r chi.Router) {
					r.Get("/", doctorAppointmentHandler.GetAppointments)
					r.Patch("/{id}", doctorAppointmentHandler.UpdateStatus)
				})
			})

			r.Route("/patient", func(r chi.Router) {
				r.Route("/doctors", func(r chi.Router) {
					r.Get("/{id}/schedules", doctorScheduleHandler.GetDoctorSchedules) //Get Schedule
					r.Get("/search", doctorHandler.SearchDoctors)                      // Search Doctors
				})

				r.Route("/appointments", func(r chi.Router) {
					r.Get("/", patientHandler.GetAppointments)                // Get Appointment
					r.Post("/", patientHandler.CreateAppointment)             // Create Appointment
					r.Get("/{id}", patientHandler.GetAppointmentDetail)       //Get Appointment detail
					r.Patch("/{id}/cancel", patientHandler.CancelAppointment) // Canceled Appointment
				})
			})

		})
	})

	// Ambil port dari environment variable yang disediakan Railway
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default untuk lokal jika variabel PORT tidak ada
    }

    log.Printf("Server berjalan di port %s", port)
    // Gunakan variabel port dengan tanda titik dua
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatal("Gagal menjalankan server: ", err)
    }
}
