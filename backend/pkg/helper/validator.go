package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation
	validate.RegisterValidation("email", validateEmail)
}

// ValidateStruct validates a struct and returns error messages
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	// Type assertion to get validation errors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fieldName := e.Field()

			// Get custom error message based on tag
			switch e.Tag() {
			case "required":
				errors[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				errors[fieldName] = "Invalid email format"
			case "min":
				errors[fieldName] = fmt.Sprintf("%s must be at least %s characters", fieldName, e.Param())
			case "max":
				errors[fieldName] = fmt.Sprintf("%s must be at most %s characters", fieldName, e.Param())
			case "oneof":
				errors[fieldName] = fmt.Sprintf("%s must be one of: %s", fieldName, e.Param())
			default:
				errors[fieldName] = fmt.Sprintf("%s is invalid", fieldName)
			}
		}
	}

	return errors
}

// validateEmail custom validation for email format
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}
