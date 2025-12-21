package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/service"
	"github.com/JinXVIII/BE-Medical-Record/pkg/helper"
	"github.com/go-chi/chi/v5"
)

type DoctorAppointmentHandler struct {
	service       service.PatientService
	doctorService service.DoctorService
}

var (
	errUserContextMissing  = errors.New("user not found in context")
	errDoctorRoleRequired  = errors.New("hanya dokter yang bisa mengakses endpoint ini")
	errDoctorProfileAbsent = errors.New("profil dokter tidak ditemukan")
)

func NewDoctorAppointmentHandler(ps service.PatientService, ds service.DoctorService) *DoctorAppointmentHandler {
	return &DoctorAppointmentHandler{service: ps, doctorService: ds}
}

func (h *DoctorAppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	doctorID, err := h.getDoctorIDFromToken(r)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, errUserContextMissing):
			status = http.StatusUnauthorized
		case errors.Is(err, errDoctorRoleRequired), errors.Is(err, errDoctorProfileAbsent):
			status = http.StatusForbidden
		default:
			status = http.StatusUnauthorized
		}
		helper.SendJSON(w, status, domain.Response{Message: err.Error()})
		return
	}

	data, err := h.service.GetDoctorAppointments(r.Context(), int64(doctorID))
	if err != nil {
		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{Message: err.Error()})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{Message: "appointments loaded", Data: data})
}

type updateAppointmentStatusRequest struct {
	Status string `json:"status"`
}

func (h *DoctorAppointmentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	doctorID, err := h.getDoctorIDFromToken(r)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, errUserContextMissing):
			status = http.StatusUnauthorized
		case errors.Is(err, errDoctorRoleRequired), errors.Is(err, errDoctorProfileAbsent):
			status = http.StatusForbidden
		default:
			status = http.StatusUnauthorized
		}
		helper.SendJSON(w, status, domain.Response{Message: err.Error()})
		return
	}

	idStr := chi.URLParam(r, "id")
	appointmentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{Message: "invalid appointment id"})
		return
	}

	var req updateAppointmentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{Message: "invalid request body"})
		return
	}

	statusValue, ok := domain.NormalizeAppointmentStatus(req.Status)
	if !ok {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{Message: "status tidak dikenal"})
		return
	}
	if statusValue == domain.AppointmentStatusPending {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{Message: "status tidak boleh kembali ke Pending"})
		return
	}

	if err := h.service.UpdateAppointmentStatus(r.Context(), int64(doctorID), appointmentID, statusValue); err != nil {
		switch {
		case errors.Is(err, service.ErrNotAllowed):
			helper.SendJSON(w, http.StatusForbidden, domain.Response{Message: err.Error()})
		case errors.Is(err, service.ErrInvalidStatus):
			helper.SendJSON(w, http.StatusBadRequest, domain.Response{Message: err.Error()})
		case errors.Is(err, sql.ErrNoRows):
			helper.SendJSON(w, http.StatusNotFound, domain.Response{Message: "appointment not found"})
		default:
			helper.SendJSON(w, http.StatusInternalServerError, domain.Response{Message: err.Error()})
		}
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{Message: "status updated"})
}

func (h *DoctorAppointmentHandler) getDoctorIDFromToken(r *http.Request) (int, error) {
	userInfo, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		return 0, errUserContextMissing
	}

	role, _ := userInfo["role"].(string)
	if role != "doctor" {
		return 0, errDoctorRoleRequired
	}

	userIDFloat, ok := userInfo["user_id"].(float64)
	if !ok {
		return 0, errUserContextMissing
	}

	doctor, err := h.doctorService.GetByUserID(r.Context(), int(userIDFloat))
	if err != nil {
		return 0, errDoctorProfileAbsent
	}

	return doctor.ID, nil
}
