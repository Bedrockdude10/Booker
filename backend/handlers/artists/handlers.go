// handlers/artists/handlers.go - Updated to use shared domain types
package artists

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Bedrockdude10/Booker/backend/domain/artists"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service *Service
}

//==============================================================================
// CRUD Operations for Admin Interface
//==============================================================================

// CreateArtist creates a new artist (admin endpoint)
func (h *Handler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	var params artists.CreateArtistParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	artist, appErr := h.service.CreateArtist(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

// GetArtist retrieves a single artist by ID (admin endpoint)
func (h *Handler) GetArtist(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	artist, appErr := h.service.GetArtistByID(r.Context(), id)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, artist)
}

// UpdateArtist performs a full update of an artist (admin endpoint)
func (h *Handler) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	var params artists.CreateArtistParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	updatedArtist, appErr := h.service.UpdateArtist(r.Context(), id, params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, updatedArtist)
}

// UpdatePartialArtist performs a partial update of an artist (admin endpoint)
func (h *Handler) UpdatePartialArtist(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	var params artists.CreateArtistParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	updatedArtist, appErr := h.service.UpdatePartialArtist(r.Context(), id, params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, updatedArtist)
}

// DeleteArtist deletes an artist (admin endpoint)
func (h *Handler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseObjectID(chi.URLParam(r, "id"))
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	if appErr := h.service.DeleteArtist(r.Context(), id); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//==============================================================================
// Admin Browse/Filter Endpoints (Limited Use)
//==============================================================================

// GetArtists provides admin interface for browsing artists with filtering
func (h *Handler) GetArtists(w http.ResponseWriter, r *http.Request) {
	// Parse filters using shared domain logic
	filters := artists.ParseFilterParams(r)

	// Validate using shared domain validation
	if appErr := artists.ValidateFilterParams(filters); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Parse pagination
	page, limit := parsePagination(r)
	offset := (page - 1) * limit

	artistsList, appErr := h.service.GetArtists(r.Context(), filters, limit, offset)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	response := map[string]interface{}{
		"data": artistsList,
		"meta": map[string]interface{}{
			"page":    page,
			"limit":   limit,
			"count":   len(artistsList),
			"hasMore": len(artistsList) == limit,
			"filters": filters,
		},
	}

	writeJSON(w, response)
}

//==============================================================================
// Helper Functions
//==============================================================================

// parsePagination extracts page and limit from query parameters
func parsePagination(r *http.Request) (page, limit int) {
	page = 1
	limit = getDefaultPageSize()

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if pageVal, err := strconv.Atoi(pageStr); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limitVal, err := strconv.Atoi(limitStr); err == nil && limitVal > 0 {
			maxPageSize := getMaxPageSize()
			if limitVal > maxPageSize {
				limitVal = maxPageSize
			}
			limit = limitVal
		}
	}

	return page, limit
}

// getDefaultPageSize returns the default page size from environment
func getDefaultPageSize() int {
	if defaultStr := os.Getenv("DEFAULT_PAGE_SIZE"); defaultStr != "" {
		if defaultVal, err := strconv.Atoi(defaultStr); err == nil && defaultVal > 0 {
			return defaultVal
		}
	}
	return 10 // fallback default
}

// getMaxPageSize returns the maximum page size from environment
func getMaxPageSize() int {
	if maxStr := os.Getenv("MAX_PAGE_SIZE"); maxStr != "" {
		if maxVal, err := strconv.Atoi(maxStr); err == nil && maxVal > 0 {
			return maxVal
		}
	}
	return 100 // fallback default
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
