// handlers/preferences/handlers.go
package preferences

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
// CRUD Operations
//==============================================================================

// CreateUserPreference creates new user preferences
func (h *Handler) CreateUserPreference(w http.ResponseWriter, r *http.Request) {
	var params CreateUserPreferenceParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate the struct using our validation package
	if appErr := validation.ValidateStruct(r.Context(), params); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	preference, appErr := h.service.CreateUserPreference(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(preference)
}

// GetUserPreference retrieves user preference by ID
func (h *Handler) GetUserPreference(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	preference, appErr := h.service.GetUserPreferenceByID(r.Context(), id)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, preference)
}

// GetUserPreferenceByAccount retrieves user preference by account ID
func (h *Handler) GetUserPreferenceByAccount(w http.ResponseWriter, r *http.Request) {
	accountID, appErr := parseObjectID(chi.URLParam(r, "accountId"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	preference, appErr := h.service.GetUserPreferenceByAccountID(r.Context(), accountID)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, preference)
}

// GetAllUserPreferences retrieves all user preferences with pagination
func (h *Handler) GetAllUserPreferences(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

	preferences, appErr := h.service.GetAllUserPreferences(r.Context(), page, limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Get total count for pagination metadata
	totalCount, appErr := h.service.CountUserPreferences(r.Context())
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	response := map[string]interface{}{
		"data": preferences,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"count":      len(preferences),
			"totalCount": totalCount,
			"hasMore":    int64(page*limit) < totalCount,
		},
	}

	writeJSON(w, response)
}

// UpdateUserPreference updates user preferences by ID
func (h *Handler) UpdateUserPreference(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	var params UpdateUserPreferenceParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate the struct
	if appErr := validation.ValidateStruct(r.Context(), params); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	updatedPreference, appErr := h.service.UpdateUserPreference(r.Context(), id, params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, updatedPreference)
}

// UpdateUserPreferenceByAccount updates user preferences by account ID
func (h *Handler) UpdateUserPreferenceByAccount(w http.ResponseWriter, r *http.Request) {
	accountID, appErr := parseObjectID(chi.URLParam(r, "accountId"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	var params UpdateUserPreferenceParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate the struct
	if appErr := validation.ValidateStruct(r.Context(), params); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	updatedPreference, appErr := h.service.UpdateUserPreferenceByAccountID(r.Context(), accountID, params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, updatedPreference)
}

// DeleteUserPreference deletes user preferences by ID
func (h *Handler) DeleteUserPreference(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	if appErr := h.service.DeleteUserPreference(r.Context(), id); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUserPreferenceByAccount deletes user preferences by account ID
func (h *Handler) DeleteUserPreferenceByAccount(w http.ResponseWriter, r *http.Request) {
	accountID, appErr := parseObjectID(chi.URLParam(r, "accountId"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	if appErr := h.service.DeleteUserPreferenceByAccountID(r.Context(), accountID); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//==============================================================================
// Query Operations
//==============================================================================

// GetPreferencesByGenre gets all users who prefer a specific genre
func (h *Handler) GetPreferencesByGenre(w http.ResponseWriter, r *http.Request) {
	genre := chi.URLParam(r, "genre")
	if genre == "" {
		utils.HandleError(w, utils.ValidationError("Genre parameter is required"))
		return
	}

	preferences, appErr := h.service.GetPreferencesByGenre(r.Context(), genre)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data":  preferences,
		"genre": genre,
		"count": len(preferences),
	})
}

// GetPreferencesByCity gets all users who prefer a specific city
func (h *Handler) GetPreferencesByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		utils.HandleError(w, utils.ValidationError("City parameter is required"))
		return
	}

	preferences, appErr := h.service.GetPreferencesByCity(r.Context(), city)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data":  preferences,
		"city":  city,
		"count": len(preferences),
	})
}

//==============================================================================
// Bulk Operations
//==============================================================================

// CreateOrUpdateUserPreference creates new preferences or updates existing ones
func (h *Handler) CreateOrUpdateUserPreference(w http.ResponseWriter, r *http.Request) {
	var params CreateUserPreferenceParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate the struct
	if appErr := validation.ValidateStruct(r.Context(), params); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Check if preferences already exist
	existing, _ := h.service.GetUserPreferenceByAccountID(r.Context(), params.AccountID)

	if existing != nil {
		// Update existing preferences
		updateParams := UpdateUserPreferenceParams{
			PreferredGenres: params.PreferredGenres,
			PreferredCities: params.PreferredCities,
		}

		updatedPreference, appErr := h.service.UpdateUserPreference(r.Context(), existing.ID, updateParams)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		writeJSON(w, updatedPreference)
		return
	}

	// Create new preferences
	preference, appErr := h.service.CreateUserPreference(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(preference)
}

//==============================================================================
// Analytics Endpoints
//==============================================================================

// GetPreferencesStats provides statistics about user preferences
func (h *Handler) GetPreferencesStats(w http.ResponseWriter, r *http.Request) {
	// Get total count
	totalCount, appErr := h.service.CountUserPreferences(r.Context())
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// You could add more analytics here like:
	// - Most popular genres
	// - Most popular cities
	// - User preference distribution

	stats := map[string]interface{}{
		"totalUsers": totalCount,
		"message":    "Additional analytics can be implemented here",
	}

	writeJSON(w, stats)
}

//==============================================================================
// Helper Functions
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
