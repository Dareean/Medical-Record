package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JinXVIII/BE-Medical-Record/internal/service"
	"github.com/go-chi/chi/v5"
)

type PatientHandler struct {
	service service.PatientService
}

func NewPatientHandler(s service.PatientService) *PatientHandler {
	return &PatientHandler{service: s}
}

type createAppointmentRequest struct {
	DoctorID        int64  `json:"doctor_id"`
	ScheduleID      *int64 `json:"schedule_id,omitempty"`
	AppointmentDate string `json:"appointment_date"` // YYYY-MM-DD
	StartTimeSlot   string `json:"start_time_slot"`
	Complaint       string `json:"complaint"`
}

func (h *PatientHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var req createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		http.Error(w, "invalid appointment_date format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	startTimeNormalized := req.StartTimeSlot
	if startTimeNormalized == "" {
		http.Error(w, "start_time_slot wajib diisi (HH:MM)", http.StatusBadRequest)
		return
	}
	if t, err := time.Parse("15:04", req.StartTimeSlot); err == nil {
		startTimeNormalized = t.Format("15:04:05")
	} else if tFull, errFull := time.Parse("15:04:05", req.StartTimeSlot); errFull == nil {
		startTimeNormalized = tFull.Format("15:04:05")
	} else {
		http.Error(w, "start_time_slot harus memiliki format HH:MM", http.StatusBadRequest)
		return
	}

	patientID, err := getPatientUserID(r)
	if err != nil {
		handlePatientAuthError(w, err)
		return
	}

	appointment, err := h.service.CreateAppointment(
		r.Context(),
		patientID,
		req.DoctorID,
		date,
		startTimeNormalized,
		req.Complaint,
		req.ScheduleID,
	)
	if err != nil {
		log.Printf("create appointment error (patient=%d doctor=%d schedule=%v): %v", patientID, req.DoctorID, req.ScheduleID, err)
		switch err {
		case service.ErrNotAllowed:
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(appointment)
}

func (h *PatientHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	patientID, err := getPatientUserID(r)
	if err != nil {
		handlePatientAuthError(w, err)
		return
	}

	data, err := h.service.GetAppointmentHistory(r.Context(), patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}

func (h *PatientHandler) GetAppointmentDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	data, err := h.service.GetAppointmentDetail(r.Context(), id)
	if err != nil {
		http.Error(w, "appointment not found", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}

func (h *PatientHandler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	appointmentID, _ := strconv.ParseInt(idStr, 10, 64)

	patientID, err := getPatientUserID(r)
	if err != nil {
		handlePatientAuthError(w, err)
		return
	}

	if err := h.service.CancelAppointment(r.Context(), patientID, appointmentID); err != nil {
		switch err {
		case service.ErrNotAllowed:
			http.Error(w, err.Error(), http.StatusForbidden)
		case service.ErrInvalidStatus:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case sql.ErrNoRows:
			http.Error(w, "appointment not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

var (
	errUnauthorized = errors.New("unauthorized")
	errForbidden    = errors.New("forbidden")
)

func getPatientUserID(r *http.Request) (int64, error) {
	userInfo, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		return 0, errUnauthorized
	}
	role, _ := userInfo["role"].(string)
	if role != "patient" {
		return 0, errForbidden
	}
	userIDFloat, ok := userInfo["user_id"].(float64)
	if !ok {
		return 0, errUnauthorized
	}
	return int64(userIDFloat), nil
}

func handlePatientAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errForbidden):
		http.Error(w, err.Error(), http.StatusForbidden)
	default:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}
