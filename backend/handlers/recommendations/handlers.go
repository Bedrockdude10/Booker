// handlers/recommendations/handlers.go
package recommendations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service *Service
}

//==============================================================================
// Core Recommendation Endpoints
//==============================================================================

// GetPersonalizedRecommendations returns personalized recommendations for a user
func (h *Handler) GetPersonalizedRecommendations(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")

	userID, appErr := parseObjectID(userIDStr)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	limit := parseLimit(r, 10)

	recommendations, appErr := h.service.GetPersonalizedRecommendations(r.Context(), userID, limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, recommendations)
}

// GetRecommendationsByGenre returns recommendations for a specific genre
func (h *Handler) GetRecommendationsByGenre(w http.ResponseWriter, r *http.Request) {
	genre := chi.URLParam(r, "genre")
	if genre == "" {
		utils.HandleError(w, utils.ValidationError("Genre parameter is required"))
		return
	}

	limit := parseLimit(r, 10)

	recommendations, appErr := h.service.GetRecommendationsByGenre(r.Context(), genre, limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, recommendations)
}

// GetRecommendationsByCity returns recommendations for a specific city
func (h *Handler) GetRecommendationsByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		utils.HandleError(w, utils.ValidationError("City parameter is required"))
		return
	}

	limit := parseLimit(r, 10)

	recommendations, appErr := h.service.GetRecommendationsByCity(r.Context(), city, limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, recommendations)
}

// GetGeneralRecommendations returns general recommendations (trending, popular)
func (h *Handler) GetGeneralRecommendations(w http.ResponseWriter, r *http.Request) {
	limit := parseLimit(r, 10)

	recommendations, appErr := h.service.GetGeneralRecommendations(r.Context(), limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, recommendations)
}

//==============================================================================
// User Interaction Endpoints
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

//==============================================================================
// Batch Recommendation Endpoints
//==============================================================================

// GetRecommendationsBatch handles complex recommendation requests
func (h *Handler) GetRecommendationsBatch(w http.ResponseWriter, r *http.Request) {
	var params RecommendationParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.HandleError(w, utils.ValidationError("Invalid request body"))
		return
	}

	// Set default limit
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var recommendations *RecommendationResponse
	var appErr *utils.AppError

	// Determine recommendation type based on parameters
	if !params.UserID.IsZero() {
		// Personalized recommendations
		recommendations, appErr = h.service.GetPersonalizedRecommendations(r.Context(), params.UserID, params.Limit)
	} else if len(params.Genres) > 0 && len(params.Cities) > 0 {
		// Mixed genre/city recommendations (implement if needed)
		recommendations, appErr = h.service.GetGeneralRecommendations(r.Context(), params.Limit)
	} else if len(params.Genres) > 0 {
		// Genre-based recommendations (first genre only for simplicity)
		recommendations, appErr = h.service.GetRecommendationsByGenre(r.Context(), params.Genres[0], params.Limit)
	} else if len(params.Cities) > 0 {
		// City-based recommendations (first city only for simplicity)
		recommendations, appErr = h.service.GetRecommendationsByCity(r.Context(), params.Cities[0], params.Limit)
	} else {
		// General recommendations
		recommendations, appErr = h.service.GetGeneralRecommendations(r.Context(), params.Limit)
	}

	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	writeJSON(w, recommendations)
}

//==============================================================================
// Favorite/Save Endpoints
//==============================================================================

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

//==============================================================================
// Analytics Endpoints
//==============================================================================

// GetRecommendationStats returns statistics about recommendations
func (h *Handler) GetRecommendationStats(w http.ResponseWriter, r *http.Request) {
	// This would return analytics about recommendation performance
	// For now, return basic stats

	stats := map[string]interface{}{
		"message": "Recommendation stats endpoint",
		"status":  "active",
		"features": []string{
			"personalized_recommendations",
			"genre_based_recommendations",
			"city_based_recommendations",
			"interaction_tracking",
		},
	}

	writeJSON(w, stats)
}

//==============================================================================
// Helper Functions
//==============================================================================

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
