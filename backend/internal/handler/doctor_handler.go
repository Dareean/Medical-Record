package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/service"
	"github.com/JinXVIII/BE-Medical-Record/pkg/helper"
)

type DoctorHandler interface {
	CreateDoctor(w http.ResponseWriter, r *http.Request)
	GetDoctorByID(w http.ResponseWriter, r *http.Request)
	GetAllDoctors(w http.ResponseWriter, r *http.Request)
	UpdateDoctor(w http.ResponseWriter, r *http.Request)
	DeleteDoctor(w http.ResponseWriter, r *http.Request)
	SearchDoctors(w http.ResponseWriter, r *http.Request)
}

type DoctorHandlerImpl struct {
	Service service.DoctorService
}

func NewDoctorHandler(service service.DoctorService) DoctorHandler {
	return &DoctorHandlerImpl{
		Service: service,
	}
}

func (h *DoctorHandlerImpl) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	if err := h.checkAdminRole(r); err != nil {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "Access denied: " + err.Error(),
			Data:    nil,
		})
		return
	}

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

	doctor, err := h.Service.CreateDoctor(r.Context(), req)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "invalid gender") {
			helper.SendJSON(w, http.StatusBadRequest, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		if strings.Contains(errorMessage, "email sudah ada") {
			helper.SendJSON(w, http.StatusConflict, domain.Response{
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
		Message: "Doctor and user created successfully",
		Data:    doctor,
	})
}

func (h *DoctorHandlerImpl) GetDoctorByID(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	if err := h.checkAdminRole(r); err != nil {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "Access denied: " + err.Error(),
			Data:    nil,
		})
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/doctors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid doctor ID",
			Data:    nil,
		})
		return
	}

	doctor, err := h.Service.GetDoctorByID(r.Context(), id)
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
		Message: "Doctor retrieved successfully",
		Data:    doctor,
	})
}

func (h *DoctorHandlerImpl) GetAllDoctors(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	if err := h.checkAdminRole(r); err != nil {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "Access denied: " + err.Error(),
			Data:    nil,
		})
		return
	}

	doctors, err := h.Service.GetAllDoctors(r.Context())
	if err != nil {
		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{
		Message: "Doctors retrieved successfully",
		Data:    doctors,
	})
}

func (h *DoctorHandlerImpl) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	if err := h.checkAdminRole(r); err != nil {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "Access denied: " + err.Error(),
			Data:    nil,
		})
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/doctors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid doctor ID",
			Data:    nil,
		})
		return
	}

	var req domain.DoctorRequest
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

	doctor, err := h.Service.UpdateDoctorByDoctorID(r.Context(), id, req)
	if err != nil {
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
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
		Message: "Doctor updated successfully",
		Data:    doctor,
	})
}

func (h *DoctorHandlerImpl) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	if err := h.checkAdminRole(r); err != nil {
		helper.SendJSON(w, http.StatusForbidden, domain.Response{
			Message: "Access denied: " + err.Error(),
			Data:    nil,
		})
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/doctors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid doctor ID",
			Data:    nil,
		})
		return
	}

	err = h.Service.DeleteDoctor(r.Context(), id)
	if err != nil {
		errorMessage := err.Error()

		// Check if doctor not found
		if strings.Contains(errorMessage, "not found") {
			helper.SendJSON(w, http.StatusNotFound, domain.Response{
				Message: "Doctor not found",
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
		Message: "Doctor deleted successfully",
		Data:    nil,
	})
}

// Helper function to check if user has admin role
func (h *DoctorHandlerImpl) checkAdminRole(r *http.Request) error {
	// Get user info from context (set by JWT middleware)
	userInfo, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		return errors.New("user not found in context")
	}

	// Check if user role is admin
	role, ok := userInfo["role"].(string)
	if !ok || role != "admin" {
		return errors.New("only admins can access this endpoint")
	}

	return nil
}

func (h *DoctorHandlerImpl) SearchDoctors(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	specIDStr := r.URL.Query().Get("specialization_id")

	specID := 0
	if specIDStr != "" {
		if v, err := strconv.Atoi(specIDStr); err == nil {
			specID = v
		}
	}

	doctors, err := h.Service.SearchDoctors(r.Context(), q, specID)
	if err != nil {
		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{
		Message: "Doctors search result",
		Data:    doctors,
	})
}

