// handlers/recommendations/handlers.go - Streamlined with filtering built-in
package recommendations

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service *Service
}

//==============================================================================
// Core Recommendation Endpoints - All with filtering support
//==============================================================================

// GetPersonalizedRecommendations returns personalized recommendations for a user with filtering
func (h *Handler) GetPersonalizedRecommendations(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")

	userID, appErr := parseObjectID(userIDStr)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Parse filters from query parameters
	filters := ParseRecommendationFilters(r)

	// Validate filters
	if appErr := ValidateRecommendationFilters(filters); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	limit := parseLimit(r, 10)
	offset := parseOffset(r, 0)

	// Always use the filtering method (it handles empty filters gracefully)
	params := EnhancedRecommendationParams{
		UserID:  userID,
		Filters: filters,
		Limit:   limit,
		Offset:  offset,
	}

	recommendations, appErr := h.service.GetPersonalizedRecommendations(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, recommendations)
}

// GetRecommendationsByGenre returns recommendations for a specific genre with additional filters
func (h *Handler) GetRecommendationsByGenre(w http.ResponseWriter, r *http.Request) {
	genre := chi.URLParam(r, "genre")
	if genre == "" {
		utils.HandleError(w, utils.ValidationError("Genre parameter is required"))
		return
	}

	// Sanitize genre to lowercase immediately
	genre = strings.ToLower(strings.TrimSpace(genre))

	// Parse additional filters (these are already sanitized by ParseRecommendationFilters)
	filters := ParseRecommendationFilters(r)

	// Add the genre from URL to filters if not already present
	if !containsString(filters.Genres, genre) {
		filters.Genres = append(filters.Genres, genre)
	}

	// Now validation will work because everything is lowercase
	if appErr := ValidateRecommendationFilters(filters); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	limit := parseLimit(r, 10)
	offset := parseOffset(r, 0)

	params := EnhancedRecommendationParams{
		Filters: filters,
		Limit:   limit,
		Offset:  offset,
	}

	recommendations, appErr := h.service.GetRecommendationsByGenre(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Update metadata to indicate this was a genre-based request
	recommendations.RequestedBy = "genre"
	if recommendations.Metadata == nil {
		recommendations.Metadata = make(map[string]interface{})
	}
	recommendations.Metadata["primaryGenre"] = genre

	writeJSON(w, recommendations)
}

// GetRecommendationsByCity returns recommendations for a specific city with additional filters
func (h *Handler) GetRecommendationsByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		utils.HandleError(w, utils.ValidationError("City parameter is required"))
		return
	}

	// Parse additional filters
	filters := ParseRecommendationFilters(r)

	// Add the city from URL to filters if not already present
	if !containsString(filters.Cities, city) {
		filters.Cities = append(filters.Cities, city)
	}

	if appErr := ValidateRecommendationFilters(filters); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	limit := parseLimit(r, 10)
	offset := parseOffset(r, 0)

	params := EnhancedRecommendationParams{
		Filters: filters,
		Limit:   limit,
		Offset:  offset,
	}

	recommendations, appErr := h.service.GetRecommendationsByCity(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Update metadata to indicate this was a city-based request
	recommendations.RequestedBy = "city"
	if recommendations.Metadata == nil {
		recommendations.Metadata = make(map[string]interface{})
	}
	recommendations.Metadata["primaryCity"] = city

	writeJSON(w, recommendations)
}

// GetGeneralRecommendations returns general recommendations with filtering
func (h *Handler) GetGeneralRecommendations(w http.ResponseWriter, r *http.Request) {
	// Parse filters from query parameters
	filters := ParseRecommendationFilters(r)

	// Validate filters
	if appErr := ValidateRecommendationFilters(filters); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	limit := parseLimit(r, 10)
	offset := parseOffset(r, 0)

	params := EnhancedRecommendationParams{
		Filters: filters,
		Limit:   limit,
		Offset:  offset,
	}

	recommendations, appErr := h.service.GetGeneralRecommendations(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	recommendations.RequestedBy = "general"
	writeJSON(w, recommendations)
}

//==============================================================================
// Filtering Support Functions
//==============================================================================

// ParseRecommendationFilters extracts filter parameters from HTTP request and sanitizes them
func ParseRecommendationFilters(r *http.Request) RecommendationFilters {
	params := RecommendationFilters{}
	query := r.URL.Query()

	// Parse genres - standardized on 'genres' parameter only
	if genresStr := query.Get("genres"); genresStr != "" {
		rawGenres := strings.Split(genresStr, ",")

		// Normalize and deduplicate genres
		genreSet := make(map[string]bool)
		for _, genre := range rawGenres {
			normalized := strings.ToLower(strings.TrimSpace(genre))
			if normalized != "" && !genreSet[normalized] {
				params.Genres = append(params.Genres, normalized)
				genreSet[normalized] = true
			}
		}
	}

	// Parse cities (comma-separated or single) - cities don't need case normalization like genres
	if citiesStr := query.Get("cities"); citiesStr != "" {
		rawCities := strings.Split(citiesStr, ",")

		// Normalize and deduplicate cities
		citySet := make(map[string]bool)
		for _, city := range rawCities {
			normalized := strings.TrimSpace(city)
			if normalized != "" && !citySet[normalized] {
				params.Cities = append(params.Cities, normalized)
				citySet[normalized] = true
			}
		}
	}

	// Parse rating filters
	if minRatingStr := query.Get("minRating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			params.MinRating = minRating
		}
	}

	if maxRatingStr := query.Get("maxRating"); maxRatingStr != "" {
		if maxRating, err := strconv.ParseFloat(maxRatingStr, 64); err == nil {
			params.MaxRating = maxRating
		}
	}

	// Parse boolean filters
	if hasManagerStr := query.Get("hasManager"); hasManagerStr != "" {
		if hasManager, err := strconv.ParseBool(hasManagerStr); err == nil {
			params.HasManager = &hasManager
		}
	}

	if hasSpotifyStr := query.Get("hasSpotify"); hasSpotifyStr != "" {
		if hasSpotify, err := strconv.ParseBool(hasSpotifyStr); err == nil {
			params.HasSpotify = &hasSpotify
		}
	}

	return params
}

// ValidateRecommendationFilters validates the filter parameters
func ValidateRecommendationFilters(filters RecommendationFilters) *utils.AppError {
	// Validate genres (they should already be sanitized to lowercase)
	for _, genre := range filters.Genres {
		if !domain.HasGenre(genre) {
			return utils.ValidationError("Invalid genre: " + genre)
		}
	}

	// Validate rating range
	if filters.MinRating < 0 || filters.MinRating > 5 {
		return utils.ValidationError("MinRating must be between 0 and 5")
	}
	if filters.MaxRating < 0 || filters.MaxRating > 5 {
		return utils.ValidationError("MaxRating must be between 0 and 5")
	}
	if filters.MinRating > 0 && filters.MaxRating > 0 && filters.MinRating > filters.MaxRating {
		return utils.ValidationError("MinRating cannot be greater than MaxRating")
	}

	return nil
}

//==============================================================================
// User Interaction Endpoints (unchanged)
//==============================================================================

// TrackInteraction logs a user interaction with an artist
func (h *Handler) TrackInteraction(w http.ResponseWriter, r *http.Request) {
	var params TrackInteractionParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Validate interaction type
	if !isValidInteractionType(params.Type) {
		utils.HandleError(w, utils.ValidationError("Invalid interaction type"))
		return
	}

	appErr := h.service.TrackInteraction(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]interface{}{
		"message": "Interaction tracked successfully",
		"type":    params.Type,
	})
}

// GetUserInteractions returns recent interactions for a user
func (h *Handler) GetUserInteractions(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")

	userID, appErr := parseObjectID(userIDStr)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	limit := parseLimit(r, 50)

	interactions, appErr := h.service.GetUserInteractions(r.Context(), userID, limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, map[string]interface{}{
		"data":   interactions,
		"total":  len(interactions),
		"userId": userID.Hex(),
	})
}

// GetRecommendationsBatch handles complex recommendation requests
func (h *Handler) GetRecommendationsBatch(w http.ResponseWriter, r *http.Request) {
	var params EnhancedRecommendationParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Sanitize genres in the request body (for batch requests)
	sanitizedGenres := make([]string, 0, len(params.Filters.Genres))
	for _, genre := range params.Filters.Genres {
		sanitized := strings.ToLower(strings.TrimSpace(genre))
		if sanitized != "" {
			sanitizedGenres = append(sanitizedGenres, sanitized)
		}
	}
	params.Filters.Genres = sanitizedGenres

	// Validate filters
	if appErr := ValidateRecommendationFilters(params.Filters); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Set default limit
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var recommendations *RecommendationResponse
	var appErr *utils.AppError

	// Use the unified filtering approach
	if !params.UserID.IsZero() {
		recommendations, appErr = h.service.GetPersonalizedRecommendations(r.Context(), params)
	} else {
		recommendations, appErr = h.service.GetGeneralRecommendations(r.Context(), params)
	}

	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, recommendations)
}

// SaveRecommendation allows users to save/favorite a recommended artist
func (h *Handler) SaveRecommendation(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID   primitive.ObjectID `json:"userId"`
		ArtistID primitive.ObjectID `json:"artistId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Track the save interaction
	params := TrackInteractionParams{
		UserID:   request.UserID,
		ArtistID: request.ArtistID,
		Type:     InteractionSave,
		Metadata: map[string]interface{}{
			"source": "recommendation",
		},
	}

	appErr := h.service.TrackInteraction(r.Context(), params)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]interface{}{
		"message":  "Artist saved successfully",
		"artistId": request.ArtistID.Hex(),
	})
}

// GetRecommendationStats returns statistics about recommendations
func (h *Handler) GetRecommendationStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"message": "Recommendation stats endpoint",
		"status":  "active",
		"features": []string{
			"personalized_recommendations",
			"genre_based_recommendations",
			"city_based_recommendations",
			"interaction_tracking",
			"filtering_support",
		},
	}

	writeJSON(w, stats)
}

//==============================================================================
// Helper Functions
//==============================================================================

// containsString checks if a slice contains a string
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
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

// parseLimit extracts and validates limit parameter
func parseLimit(r *http.Request, defaultLimit int) int {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return defaultLimit
	}

	// Cap maximum limit
	maxLimit := 100
	if limit > maxLimit {
		return maxLimit
	}

	return limit
}

// parseOffset extracts and validates offset parameter
func parseOffset(r *http.Request, defaultOffset int) int {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		return defaultOffset
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return defaultOffset
	}

	return offset
}

// isValidInteractionType validates interaction types
func isValidInteractionType(interactionType InteractionType) bool {
	validTypes := []InteractionType{
		InteractionView,
		InteractionLike,
		InteractionSave,
		InteractionContact,
		InteractionSkip,
	}

	for _, validType := range validTypes {
		if interactionType == validType {
			return true
		}
	}

	return false
}

// writeJSON is a helper to write JSON responses
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
