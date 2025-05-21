package artists

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Router maps endpoints to handlers for artist-related operations
*/
func Routes(r chi.Router, collections map[string]*mongo.Collection) {
	service := NewService(collections)
	handler := Handler{service}

	// Mount artist routes under /api/v1/artists
	r.Route("/api/v1/artists", func(r chi.Router) {
		// Basic CRUD operations
		r.Post("/", handler.CreateArtist)
		r.Get("/", handler.GetArtists)
		r.Get("/{id}", handler.GetArtist)
		r.Put("/{id}", handler.UpdateArtist)
		r.Patch("/{id}", handler.UpdatePartialArtist)
		r.Delete("/{id}", handler.DeleteArtist)

		// Specialized routes
		r.Get("/genre/{genre}", handler.GetArtistsByGenre)
		r.Get("/city/{city}", handler.GetArtistsByCity)

		// Recommendation routes
		r.Get("/recommendations", handler.GetRecommendations)
		r.Get("/recommendations/genre/{genre}", handler.GetRecommendationsByGenre)
		r.Get("/recommendations/city/{city}", handler.GetRecommendationsByLocation)

		// // User preference routes
		// r.Post("/preferences", handler.SaveUserPreferences)
		// r.Get("/preferences/{userId}", handler.GetUserPreferences)
		// r.Get("/recommendations/user/{userId}", handler.GetPersonalizedRecommendations)
	})
}
