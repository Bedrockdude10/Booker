// validation/validate.go
package validation

import (
	"context"
	"fmt"
	"strings"

	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register your custom genre validator
	validate.RegisterValidation("validgenres", validateGenres)
	validate.RegisterValidation("validrole", validateRole) // Custom validator for roles during account creation
}

// Simple validation function
func ValidateStruct(ctx context.Context, s interface{}) *utils.AppError {
	if err := validate.Struct(s); err != nil {
		var messages []string
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, formatError(err))
		}
		return utils.ValidationErrorLog(ctx, "Validation failed", strings.Join(messages, "; "))
	}
	return nil
}

// Custom validator for genres
func validateGenres(fl validator.FieldLevel) bool {
	genres := fl.Field().Interface().([]string)
	for _, genre := range genres {
		if !domain.HasGenre(genre) {
			return false
		}
	}
	return true
}

// Custom validator for roles
func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return domain.HasRole(role)
}

func formatError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must have at least %s items", err.Field(), err.Param())
	case "validgenres":
		return fmt.Sprintf("%s contains invalid genres", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
