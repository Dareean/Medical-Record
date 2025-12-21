package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/service"
	"github.com/JinXVIII/BE-Medical-Record/pkg/helper"
	"github.com/go-chi/chi/v5"
)

type DoctorScheduleHandler interface {
	GetMySchedules(w http.ResponseWriter, r *http.Request)
	CreateSchedule(w http.ResponseWriter, r *http.Request)
	UpdateSchedule(w http.ResponseWriter, r *http.Request)
	DeleteSchedule(w http.ResponseWriter, r *http.Request)
	GetDoctorSchedules(w http.ResponseWriter, r *http.Request)
}

type DoctorScheduleHandlerImpl struct {
	ScheduleService service.DoctorScheduleService
	DoctorService   service.DoctorService
}

var (
	errSchedUserContextMissing = errors.New("user not found in context")
	errSchedDoctorRoleRequired = errors.New("access denied: only doctors can access this endpoint")
	errSchedInvalidUserID      = errors.New("invalid user ID in token")
	errSchedDoctorProfile      = errors.New("doctor not found for this user")
)

func respondScheduleAuthError(w http.ResponseWriter, err error) {
	status := http.StatusUnauthorized
	switch {
	case errors.Is(err, errSchedUserContextMissing), errors.Is(err, errSchedInvalidUserID):
		status = http.StatusUnauthorized
	case errors.Is(err, errSchedDoctorRoleRequired), errors.Is(err, errSchedDoctorProfile):
		status = http.StatusForbidden
	default:
		status = http.StatusUnauthorized
	}
	helper.SendJSON(w, status, domain.Response{
		Message: err.Error(),
		Data:    nil,
	})
}

func NewDoctorScheduleHandler(scheduleService service.DoctorScheduleService, doctorService service.DoctorService) DoctorScheduleHandler {
	return &DoctorScheduleHandlerImpl{
		ScheduleService: scheduleService,
		DoctorService:   doctorService,
	}
}

func (h *DoctorScheduleHandlerImpl) GetMySchedules(w http.ResponseWriter, r *http.Request) {
	// Get doctor ID from JWT token
	doctorID, err := h.getDoctorIDFromToken(r)
	if err != nil {
		respondScheduleAuthError(w, err)
		return
	}

	schedules, err := h.ScheduleService.GetSchedulesByDoctorID(r.Context(), doctorID)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "doctor not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: errorMessage,
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{
		Message: "Doctor schedules retrieved successfully",
		Data:    schedules,
	})
}

func (h *DoctorScheduleHandlerImpl) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	// Get doctor ID from JWT token
	doctorID, err := h.getDoctorIDFromToken(r)
	if err != nil {
		respondScheduleAuthError(w, err)
		return
	}

	var req domain.DoctorScheduleRequest

	// Parsing body request
	if err := helper.ParseBody(r, &req); err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid request body: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Set doctor ID from token
	req.DoctorID = doctorID

	// Validation
	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Validation failed",
			Data:    validationErrors,
		})
		return
	}

	schedule, err := h.ScheduleService.CreateSchedule(r.Context(), req)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "invalid work day") {
			helper.SendJSON(w, http.StatusBadRequest, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		if strings.Contains(errorMessage, "doctor not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: errorMessage,
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusCreated, domain.Response{
		Message: "Doctor schedule created successfully",
		Data:    schedule,
	})
}

func (h *DoctorScheduleHandlerImpl) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	// Get doctor ID from JWT token
	doctorID, err := h.getDoctorIDFromToken(r)
	if err != nil {
		respondScheduleAuthError(w, err)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/doctor/schedules/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid schedule ID",
			Data:    nil,
		})
		return
	}

	// Check if schedule belongs to this doctor
	existingSchedule, err := h.ScheduleService.GetScheduleByID(r.Context(), id)
	if err != nil {
		helper.SendJSON(w, http.StatusNotFound, domain.Response{
			Message: "Schedule not found",
			Data:    nil,
		})
		return
	}

	if existingSchedule.DoctorID != doctorID {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "Access denied: You can only update your own schedules",
			Data:    nil,
		})
		return
	}

	var req domain.DoctorScheduleRequest
	if err := helper.ParseBody(r, &req); err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid request body: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Set doctor ID from token
	req.DoctorID = doctorID

	// Validation
	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Validation failed",
			Data:    validationErrors,
		})
		return
	}

	schedule, err := h.ScheduleService.UpdateSchedule(r.Context(), id, req)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		if strings.Contains(errorMessage, "invalid work day") {
			helper.SendJSON(w, http.StatusBadRequest, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: errorMessage,
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{
		Message: "Doctor schedule updated successfully",
		Data:    schedule,
	})
}

func (h *DoctorScheduleHandlerImpl) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	// Get doctor ID from JWT token
	doctorID, err := h.getDoctorIDFromToken(r)
	if err != nil {
		respondScheduleAuthError(w, err)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/doctor/schedules/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid schedule ID",
			Data:    nil,
		})
		return
	}

	// Check if schedule belongs to this doctor
	existingSchedule, err := h.ScheduleService.GetScheduleByID(r.Context(), id)
	if err != nil {
		helper.SendJSON(w, http.StatusNotFound, domain.Response{
			Message: "Schedule not found",
			Data:    nil,
		})
		return
	}

	if existingSchedule.DoctorID != doctorID {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "Access denied: You can only delete your own schedules",
			Data:    nil,
		})
		return
	}

	err = h.ScheduleService.DeleteSchedule(r.Context(), id)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: errorMessage,
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{
		Message: "Doctor schedule deleted successfully",
		Data:    nil,
	})
}

// Helper function to extract doctor ID from JWT token
func (h *DoctorScheduleHandlerImpl) getDoctorIDFromToken(r *http.Request) (int, error) {
	// Get user info from context (set by JWT middleware)
	userInfo, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		return 0, errSchedUserContextMissing
	}

	// Check if user role is doctor
	role, ok := userInfo["role"].(string)
	if !ok || role != "doctor" {
		return 0, errSchedDoctorRoleRequired
	}

	// Get user ID
	userIDFloat, ok := userInfo["user_id"].(float64)
	if !ok {
		return 0, errSchedInvalidUserID
	}

	userID := int(userIDFloat)

	// Get doctor ID from user ID
	doctor, err := h.DoctorService.GetByUserID(r.Context(), userID)
	if err != nil {
		return 0, errSchedDoctorProfile
	}

	return doctor.ID, nil

}

func (h *DoctorScheduleHandlerImpl) GetDoctorSchedules(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	doctorID, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "invalid doctor id",
			Data:    nil,
		})
		return
	}

	schedules, err := h.ScheduleService.GetSchedulesByDoctorID(
		r.Context(),
		doctorID,
	)
	if err != nil {
		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{
		Message: "Doctor schedules retrieved successfully",
		Data:    schedules,
	})
}
