// handlers/discovery/routes.go
package discovery

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// Routes sets up the discovery endpoints
func Routes(r chi.Router, collections map[string]*mongo.Collection) {
	service := NewBandcampService(collections["scrapedArtists"])
	handler := &Handler{service: service}

	// Mount discovery routes under /api/discovery
	r.Route("/api/discovery", func(r chi.Router) {
		// Bandcamp scraping endpoints
		r.Post("/scrape/bandcamp", handler.ScrapeBandcamp)
		r.Get("/artists", handler.GetScrapedArtists)
		r.Get("/artists/count", handler.GetArtistCount)

		// Future endpoints for other sources
		// r.Post("/scrape/spotify", handler.ScrapeSpotify)
		// r.Post("/enrich/spotify", handler.EnrichWithSpotify)
	})
}
