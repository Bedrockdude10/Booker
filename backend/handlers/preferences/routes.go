package preferences

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Router maps endpoints to handlers for user preferences operations
*/
func Routes(r chi.Router, collections map[string]*mongo.Collection) {
	service := NewService(collections)
	handler := &Handler{service: service}

	// Mount preferences routes under /api/preferences
	r.Route("/api/preferences", func(r chi.Router) {
		// Basic CRUD operations
		r.Post("/", handler.CreateUserPreference)       // Create new user preferences
		r.Get("/", handler.GetAllUserPreferences)       // Get all preferences (admin/analytics)
		r.Get("/{id}", handler.GetUserPreference)       // Get preference by ID
		r.Put("/{id}", handler.UpdateUserPreference)    // Update preference by ID
		r.Delete("/{id}", handler.DeleteUserPreference) // Delete preference by ID

		// Account-based operations (most common usage)
		r.Route("/account/{accountId}", func(r chi.Router) {
			r.Get("/", handler.GetUserPreferenceByAccount)       // Get preferences for account
			r.Put("/", handler.UpdateUserPreferenceByAccount)    // Update preferences for account
			r.Delete("/", handler.DeleteUserPreferenceByAccount) // Delete preferences for account
		})

		// Bulk operations
		r.Post("/upsert", handler.CreateOrUpdateUserPreference) // Create or update preferences

		// Query operations for analytics/recommendations
		r.Get("/genre/{genre}", handler.GetPreferencesByGenre) // Get all users who prefer a genre
		r.Get("/city/{city}", handler.GetPreferencesByCity)    // Get all users who prefer a city

		// Analytics endpoints
		r.Get("/stats", handler.GetPreferencesStats) // Get preference statistics

		// Health check for preferences service
		r.Get("/health", handler.HealthCheck)
	})
}

// Health check specifically for preferences service
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]interface{}{
		"status":  "healthy",
		"service": "preferences",
	})
}
