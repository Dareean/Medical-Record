package handler

import (
	"net/http"
	"strings"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/internal/service"
	"github.com/JinXVIII/BE-Medical-Record/pkg/helper"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type UserHandlerImpl struct {
	Service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &UserHandlerImpl{
		Service: service,
	}
}

func (handler *UserHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest

	// Parsing body request
	if err := helper.ParseBody(r, &req); err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid request body: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Validation input
	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Validation failed",
			Data:    validationErrors,
		})
		return
	}

	createdUser, err := handler.Service.Register(r.Context(), req)
	if err != nil {
		errorMessage := err.Error()

		// Jika error terkait email sudah ada, berikan status conflict
		if strings.Contains(errorMessage, "email sudah ada") {
			helper.SendJSON(w, http.StatusConflict, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		// Jika error terkait role tidak valid
		if strings.Contains(errorMessage, "invalid role") {
			helper.SendJSON(w, http.StatusBadRequest, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		// Untuk error lainnya, berikan status internal server error
		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: errorMessage,
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusCreated, domain.Response{
		Message: "User berhasil registrasi",
		Data:    createdUser,
	})
}

func (h *UserHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest

	// Parsing body request
	if err := helper.ParseBody(r, &req); err != nil {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Invalid request body: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Validation input
	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		helper.SendJSON(w, http.StatusBadRequest, domain.Response{
			Message: "Validation failed",
			Data:    validationErrors,
		})
		return
	}

	resp, err := h.Service.Login(r.Context(), req)
	if err != nil {
		// Error handling yang konsisten dengan register
		errorMessage := err.Error()

		// Jika error terkait credentials
		if strings.Contains(errorMessage, "invalid credentials") ||
			strings.Contains(errorMessage, "user not found") ||
			strings.Contains(errorMessage, "password") {
			helper.SendJSON(w, http.StatusUnauthorized, domain.Response{
				Message: errorMessage,
				Data:    nil,
			})
			return
		}

		// Error lainnya
		helper.SendJSON(w, http.StatusInternalServerError, domain.Response{
			Message: errorMessage,
			Data:    nil,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, domain.Response{
		Message: "Login successful",
		Data:    resp,
	})
}
