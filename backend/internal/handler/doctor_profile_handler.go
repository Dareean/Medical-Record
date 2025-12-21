package handler

import (
	"net/http"
	"strings"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/service"
	"github.com/JinXVIII/BE-Medical-Record/pkg/helper"
)

type DoctorProfileHandler interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)
	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
}

type DoctorProfileHandlerImpl struct {
	DoctorService service.DoctorService
}

func NewDoctorProfileHandler(doctorService service.DoctorService) DoctorProfileHandler {
	return &DoctorProfileHandlerImpl{
		DoctorService: doctorService,
	}
}

func (h *DoctorProfileHandlerImpl) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	// Get user info from JWT token
	userInfo, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		helper.SendJSON(w, http.StatusUnauthorized, domain.Response{
			Message: "user not found in context",
			Data:    nil,
		})
		return
	}

	// Check if user role is doctor
	role, ok := userInfo["role"].(string)
	if !ok || role != "doctor" {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "access denied: only doctors can access this endpoint",
			Data:    nil,
		})
		return
	}

	// Get user ID
	userIDFloat, ok := userInfo["user_id"].(float64)
	if !ok {
		helper.SendJSON(w, http.StatusUnauthorized, domain.Response{
			Message: "invalid user ID in token",
			Data:    nil,
		})
		return
	}

	userID := int(userIDFloat)

	doctor, err := h.DoctorService.GetByUserID(r.Context(), userID)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
				Message: "doctor profile not found",
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
		Message: "Doctor profile retrieved successfully",
		Data:    doctor,
	})
}

func (h *DoctorProfileHandlerImpl) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	// Get user info from JWT token
	userInfo, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		helper.SendJSON(w, http.StatusUnauthorized, domain.Response{
			Message: "user not found in context",
			Data:    nil,
		})
		return
	}

	// Check if user role is doctor
	role, ok := userInfo["role"].(string)
	if !ok || role != "doctor" {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "access denied: only doctors can access this endpoint",
			Data:    nil,
		})
		return
	}

	// Get user ID
	userIDFloat, ok := userInfo["user_id"].(float64)
	if !ok {
		helper.SendJSON(w, http.StatusUnauthorized, domain.Response{
			Message: "invalid user ID in token",
			Data:    nil,
		})
		return
	}

	userID := int(userIDFloat)

	var req domain.DoctorRequest

	// Parsing body request
	if err := helper.ParseBody(r, &req); err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid request body: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Validation
	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Validation failed",
			Data:    validationErrors,
		})
		return
	}

	doctor, err := h.DoctorService.UpdateDoctorByUserID(r.Context(), userID, req)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
				Message: "doctor profile not found",
				Data:    nil,
			})
			return
		}

		if strings.Contains(errorMessage, "email already exists") {
			helper.SendJSON(w, http.StatusConflict, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		if strings.Contains(errorMessage, "invalid gender") {
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
		Message: "Doctor profile updated successfully",
		Data:    doctor,
	})
}
