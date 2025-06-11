// handlers/accounts/handlers.go
package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/Bedrockdude10/Booker/backend/validation"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service    *Service
	jwtService *JWTService
}

// NewHandler creates a new accounts handler
func NewHandler(service *Service, jwtService *JWTService) *Handler {
	return &Handler{
		service:    service,
		jwtService: jwtService,
	}
}

// Response structures
type LoginResponse struct {
	Token   string  `json:"token"`
	Account Account `json:"user"`
}

type RegisterResponse struct {
	Token   string  `json:"token"`
	Account Account `json:"user"`
	Message string  `json:"message"`
}

//==============================================================================
// Authentication Handlers
//==============================================================================

// Register creates a new user account and returns a JWT token
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var params CreateAccountParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate the struct using our validation package
	if appErr := validation.ValidateStruct(r.Context(), &params); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	account, appErr := h.service.CreateAccount(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(account)
	if err != nil {
		utils.HandleError(w, utils.InternalError("Failed to generate token", err))
		return
	}

	// Return successful response with token
	response := RegisterResponse{
		Token:   token,
		Account: *account,
		Message: "Account created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication and returns JWT token
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
	if appErr := validation.ValidateStruct(r.Context(), &credentials); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	account, appErr := h.service.VerifyPassword(r.Context(), credentials.Email, credentials.Password)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(account)
	if err != nil {
		utils.HandleError(w, utils.InternalError("Failed to generate token", err))
		return
	}

	// Return successful response with token
	response := LoginResponse{
		Token:   token,
		Account: *account,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken generates a new token from an existing valid token
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.HandleError(w, utils.ValidationError("Authorization header required"))
		return
	}

	// Extract token from "Bearer <token>" format
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		utils.HandleError(w, utils.ValidationError("Invalid authorization header format"))
		return
	}

	// Refresh the token
	newToken, err := h.jwtService.RefreshToken(tokenParts[1])
	if err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid or expired token"))
		return
	}

	// Return new token
	response := map[string]string{
		"token": newToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//==============================================================================
// Account Management Handlers
//==============================================================================

// CreateAccount handles user registration (legacy - use Register instead)
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	h.Register(w, r)
}

// GetAccount retrieves a single account by ID or current user's account
func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	// If no ID in URL, get current user's account from JWT context
	if idParam == "" {
		claims, ok := r.Context().Value("user").(*Claims)
		if !ok {
			utils.HandleError(w, utils.ValidationError("User not found in context"))
			return
		}

		account, appErr := h.service.GetAccountByID(r.Context(), claims.UserID)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		writeJSON(w, account)
		return
	}

	// Otherwise, get account by ID (admin function)
	id, appErr := parseObjectID(idParam)
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
	idParam := chi.URLParam(r, "id")

	// If no ID in URL, update current user's account from JWT context
	if idParam == "" {
		claims, ok := r.Context().Value("user").(*Claims)
		if !ok {
			utils.HandleError(w, utils.ValidationError("User not found in context"))
			return
		}

		var params UpdateAccountParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			utils.HandleError(w, utils.ValidationError("Invalid request body"))
			return
		}

		// Validate the struct
		if appErr := validation.ValidateStruct(r.Context(), &params); appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		updatedAccount, appErr := h.service.UpdateAccount(r.Context(), claims.UserID, params)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		writeJSON(w, updatedAccount)
		return
	}

	// Otherwise, update account by ID (admin function)
	id, appErr := parseObjectID(idParam)
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
	if appErr := validation.ValidateStruct(r.Context(), &params); appErr != nil {
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

// ChangePassword allows users to change their password
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get current user from JWT context or ID parameter
	var userID primitive.ObjectID
	var currentEmail string

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		// Use current user from JWT
		claims, ok := r.Context().Value("user").(*Claims)
		if !ok {
			utils.HandleError(w, utils.ValidationError("User not found in context"))
			return
		}
		userID = claims.UserID
		currentEmail = claims.Email
	} else {
		// Admin changing someone else's password
		id, appErr := parseObjectID(idParam)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}
		userID = id

		// Get account to get email
		account, appErr := h.service.GetAccountByID(r.Context(), userID)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}
		currentEmail = account.Email
	}

	var changePasswordRequest struct {
		CurrentPassword string `json:"currentPassword,omitempty"`
		NewPassword     string `json:"newPassword" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&changePasswordRequest); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	if appErr := validation.ValidateStruct(r.Context(), &changePasswordRequest); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Verify current password if provided (for user self-service)
	if changePasswordRequest.CurrentPassword != "" {
		_, appErr := h.service.VerifyPassword(r.Context(), currentEmail, changePasswordRequest.CurrentPassword)
		if appErr != nil {
			utils.HandleError(w, utils.ValidationError("Current password is incorrect"))
			return
		}
	}

	if appErr := h.service.UpdatePassword(r.Context(), userID, changePasswordRequest.NewPassword); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, map[string]string{"message": "Password updated successfully"})
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

	if appErr := validation.ValidateStruct(r.Context(), &request); appErr != nil {
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

// UpdatePassword handles password updates (legacy - use ChangePassword instead)
func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	h.ChangePassword(w, r)
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
