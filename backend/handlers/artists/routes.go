// handlers/artists/routes.go - Simplified for admin use only
package artists

import (
	"net/http"

	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Routes maps endpoints to handlers for artist admin operations.
Note: User-facing discovery endpoints are handled by recommendations service.
*/
func Routes(r chi.Router, collections map[string]*mongo.Collection) {
	service := NewService(collections)
	handler := &Handler{service: service}

	// Mount artist routes under /api/artists (admin interface)
	r.Route("/api/artists", func(r chi.Router) {
		//==============================================================================
		// CRUD Operations (Admin Interface)
		//==============================================================================
		r.Post("/", handler.CreateArtist)             // Create new artist
		r.Get("/{id}", handler.GetArtist)             // Get single artist by ID
		r.Put("/{id}", handler.UpdateArtist)          // Full update
		r.Patch("/{id}", handler.UpdatePartialArtist) // Partial update
		r.Delete("/{id}", handler.DeleteArtist)       // Delete artist

		//==============================================================================
		// Admin Browse/Filter (Limited Use)
		//==============================================================================
		r.Get("/", handler.GetArtists) // Admin browse with filtering

		//==============================================================================
		// Utility Endpoints
		//==============================================================================
		r.Get("/genres", handler.GetAllGenres) // List all available genres
	})
}

// GetAllGenres retrieves all available genres for admin interface
func (h *Handler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres := domain.GetAllGenres()

	writeJSON(w, map[string]interface{}{
		"data":  genres,
		"count": len(genres),
	})
}
