// handlers/recommendations/routes.go
package recommendations

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Routes maps endpoints to handlers for recommendation operations
*/
func Routes(r chi.Router, collections map[string]*mongo.Collection) {
	service := NewService(collections)
	handler := &Handler{service: service}

	// Mount recommendation routes under /api/recommendations
	r.Route("/api/recommendations", func(r chi.Router) {

		//==============================================================================
		// Core Recommendation Endpoints
		//==============================================================================

		// General recommendations (no authentication required)
		r.Get("/", handler.GetGeneralRecommendations)

		// Genre-based recommendations
		r.Get("/genre/{genre}", handler.GetRecommendationsByGenre)

		// City-based recommendations
		r.Get("/city/{city}", handler.GetRecommendationsByCity)

		// Personalized recommendations (requires user ID)
		r.Get("/user/{userId}", handler.GetPersonalizedRecommendations)

		// Batch recommendations (complex queries via POST)
		r.Post("/batch", handler.GetRecommendationsBatch)

		//==============================================================================
		// User Interaction Endpoints
		//==============================================================================

		// Track user interactions (views, likes, saves, etc.)
		r.Post("/interactions", handler.TrackInteraction)

		// Get user interaction history
		r.Get("/interactions/user/{userId}", handler.GetUserInteractions)

		// Save/favorite recommendations
		r.Post("/save", handler.SaveRecommendation)

		//==============================================================================
		// Analytics & Stats Endpoints
		//==============================================================================

		// Recommendation statistics and performance metrics
		r.Get("/stats", handler.GetRecommendationStats)

		//==============================================================================
		// Health Check
		//==============================================================================

		// Health check for recommendations service
		r.Get("/health", handler.HealthCheck)
	})
}

// HealthCheck for recommendations service
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]interface{}{
		"status":  "healthy",
		"service": "recommendations",
		"version": "1.0",
		"features": map[string]bool{
			"personalized_recommendations": true,
			"genre_recommendations":        true,
			"city_recommendations":         true,
			"interaction_tracking":         true,
			"batch_processing":             true,
		},
	})
}
