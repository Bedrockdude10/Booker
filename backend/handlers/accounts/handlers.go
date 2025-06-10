// handlers/accounts/handlers.go
package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/Bedrockdude10/Booker/backend/validation"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service *Service
}

//==============================================================================
// Account Creation and Management
//==============================================================================

// CreateAccount handles user registration
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var params CreateAccountParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate the struct using our validation package
	if appErr := validation.ValidateStruct(r.Context(), params); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	account, appErr := h.service.CreateAccount(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetAccount retrieves a single account by ID
func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	account, appErr := h.service.GetAccountByID(r.Context(), id)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, account)
}

// GetAccountByEmail retrieves account by email address
func (h *Handler) GetAccountByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		utils.HandleError(w, utils.ValidationError("Email parameter is required"))
		return
	}

	account, appErr := h.service.GetAccountByEmail(r.Context(), email)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, account)
}

// UpdateAccount handles account updates
func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	var params UpdateAccountParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate the struct
	if appErr := validation.ValidateStruct(r.Context(), params); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	updatedAccount, appErr := h.service.UpdateAccount(r.Context(), id, params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, updatedAccount)
}

// DeactivateAccount handles soft deletion of accounts
func (h *Handler) DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	if appErr := h.service.DeactivateAccount(r.Context(), id); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ActivateAccount handles reactivation of deactivated accounts
func (h *Handler) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// You'll need to implement this method in the service
	if appErr := h.service.ActivateAccount(r.Context(), id); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]string{"message": "Account activated successfully"})
}

//==============================================================================
// Authentication
//==============================================================================

// Login handles user authentication
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate credentials struct
	if appErr := validation.ValidateStruct(r.Context(), credentials); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	account, appErr := h.service.VerifyPassword(r.Context(), credentials.Email, credentials.Password)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// In a real application, you'd generate a JWT token here
	response := map[string]interface{}{
		"message": "Login successful",
		"account": account,
		// "token": generateJWT(account), // You'll implement this later
	}

	writeJSON(w, response)
}

// RequestPasswordReset handles password reset requests
func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	if appErr := validation.ValidateStruct(r.Context(), request); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Check if account exists
	_, appErr := h.service.GetActiveAccountByEmail(r.Context(), request.Email)
	if appErr != nil {
		// Don't reveal whether email exists or not for security
		writeJSON(w, map[string]string{
			"message": "If an account with that email exists, a password reset link has been sent",
		})
		return
	}

	// In a real application, send email with reset token
	// For now, just return success message
	writeJSON(w, map[string]string{
		"message": "Password reset email sent",
	})
}

// UpdatePassword handles password updates
func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	var request struct {
		NewPassword string `json:"newPassword" validate:"required,min=8"`
		// You might also want current password for verification
		CurrentPassword string `json:"currentPassword,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	if appErr := validation.ValidateStruct(r.Context(), request); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Optional: Verify current password before updating
	if request.CurrentPassword != "" {
		account, appErr := h.service.GetAccountByID(r.Context(), id)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		// Verify current password
		if _, appErr := h.service.VerifyPassword(r.Context(), account.Email, request.CurrentPassword); appErr != nil {
			utils.HandleError(w, utils.ValidationError("Current password is incorrect"))
			return
		}
	}

	if appErr := h.service.UpdatePassword(r.Context(), id, request.NewPassword); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, map[string]string{"message": "Password updated successfully"})
}

//==============================================================================
// Admin Operations
//==============================================================================

// ListAccounts handles listing all accounts (admin only)
func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

	accounts, appErr := h.service.ListAccounts(r.Context(), page, limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Get total count for pagination metadata
	totalCount, appErr := h.service.CountAccounts(r.Context())
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	response := map[string]interface{}{
		"data": accounts,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"count":      len(accounts),
			"totalCount": totalCount,
			"hasMore":    int64(page*limit) < totalCount,
		},
	}

	writeJSON(w, response)
}

//==============================================================================
// Helper Functions (copied from artists handlers)
//==============================================================================

// parsePagination extracts page and limit from query parameters
func parsePagination(r *http.Request) (page, limit int) {
	page = 1
	limit = 10 // Default page size

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if pageVal, err := strconv.Atoi(pageStr); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limitVal, err := strconv.Atoi(limitStr); err == nil && limitVal > 0 {
			maxPageSize := 100 // Maximum page size
			if limitVal > maxPageSize {
				limitVal = maxPageSize
			}
			limit = limitVal
		}
	}

	return page, limit
}

// parseObjectID converts string to ObjectID with proper error handling
func parseObjectID(idStr string) (primitive.ObjectID, *utils.AppError) {
	if idStr == "" {
		return primitive.NilObjectID, utils.ValidationError("ID parameter is required")
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.NilObjectID, utils.ValidationError("Invalid ID format")
	}

	return id, nil
}

// writeJSON is a helper to write JSON responses
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
