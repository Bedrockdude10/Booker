// utils/errors.go - Improved version
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeNotFound   ErrorType = "not_found"
	ErrorTypeDatabase   ErrorType = "database"
	ErrorTypeAuth       ErrorType = "authentication"
	ErrorTypeInternal   ErrorType = "internal"
	ErrorTypeExternal   ErrorType = "external_api"
)

// AppError represents a structured application error
type AppError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"-"`
	Err        error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Simplified error constructors
func ValidationError(message string, details ...string) *AppError {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:       ErrorTypeValidation,
		Message:    message,
		Details:    detail,
		StatusCode: http.StatusBadRequest,
	}
}

func NotFound(resource string) *AppError {
	return &AppError{
		Type:       ErrorTypeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func DatabaseError(operation string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeDatabase,
		Message:    fmt.Sprintf("Database operation failed: %s", operation),
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func InternalError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// Auto-logging error constructors (use these when you want automatic logging)
func ValidationErrorLog(ctx context.Context, message string, details ...string) *AppError {
	err := ValidationError(message, details...)
	logError(ctx, err, "Validation error occurred")
	return err
}

func NotFoundLog(ctx context.Context, resource string) *AppError {
	err := NotFound(resource)
	logError(ctx, err, "Resource not found")
	return err
}

func DatabaseErrorLog(ctx context.Context, operation string, err error) *AppError {
	appErr := DatabaseError(operation, err)
	logError(ctx, appErr, "Database error occurred")
	return appErr
}

func InternalErrorLog(ctx context.Context, message string, err error) *AppError {
	appErr := InternalError(message, err)
	logError(ctx, appErr, "Internal error occurred")
	return appErr
}

// Centralized logging function
func logError(ctx context.Context, appErr *AppError, msg string) {
	logger := slog.Default()

	// Create log attributes
	logArgs := []any{
		"error_type", appErr.Type,
		"message", appErr.Message,
		"status_code", appErr.StatusCode,
	}

	if appErr.Details != "" {
		logArgs = append(logArgs, "details", appErr.Details)
	}

	// Log with appropriate level based on error type
	switch appErr.Type {
	case ErrorTypeValidation:
		logger.WarnContext(ctx, msg, logArgs...)
	case ErrorTypeNotFound:
		logger.InfoContext(ctx, msg, logArgs...)
	case ErrorTypeDatabase, ErrorTypeInternal:
		logger.ErrorContext(ctx, msg, logArgs...)
	default:
		logger.WarnContext(ctx, msg, logArgs...)
	}
}

// Manual logging with custom attributes (for when you need more control)
func Log(ctx context.Context, err *AppError, msg string, attrs ...any) *AppError {
	logger := slog.Default()

	// Base attributes
	logArgs := []any{
		"error_type", err.Type,
		"message", err.Message,
		"status_code", err.StatusCode,
	}

	// Add custom attributes
	logArgs = append(logArgs, attrs...)

	// Log with appropriate level
	switch err.Type {
	case ErrorTypeValidation:
		logger.WarnContext(ctx, msg, logArgs...)
	case ErrorTypeNotFound:
		logger.InfoContext(ctx, msg, logArgs...)
	case ErrorTypeDatabase, ErrorTypeInternal:
		logger.ErrorContext(ctx, msg, logArgs...)
	default:
		logger.WarnContext(ctx, msg, logArgs...)
	}

	return err
}

// HTTPError writes standardized error responses
func HTTPError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"type":    err.Type,
			"message": err.Message,
		},
	}

	if err.Details != "" {
		response["error"].(map[string]interface{})["details"] = err.Details
	}

	json.NewEncoder(w).Encode(response)
}

// Helper function to handle and respond with errors in one line
func HandleError(w http.ResponseWriter, err *AppError) {
	HTTPError(w, err)
}

// Recovery middleware for panics
func ErrorHandleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered",
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
				)

				appErr := InternalError("Internal server error", fmt.Errorf("panic: %v", err))
				HTTPError(w, appErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Convenience function for common MongoDB error patterns
func HandleMongoError(ctx context.Context, err error, operation string, resourceType string) *AppError {
	if err == nil {
		return nil
	}

	// Import mongo driver to check specific errors
	// if err == mongo.ErrNoDocuments {
	//     return NotFoundLog(ctx, resourceType)
	// }

	return DatabaseErrorLog(ctx, operation, err)
}
