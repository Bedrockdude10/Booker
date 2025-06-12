// handlers/discovery/handlers.go
package discovery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bedrockdude10/Booker/backend/utils"
)

type Handler struct {
	service *BandcampService
}

// ScrapeBandcamp triggers scraping of Boston artists from Bandcamp
func (h *Handler) ScrapeBandcamp(w http.ResponseWriter, r *http.Request) {
	// Parse optional limit parameter
	limit := 1000 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 5000 {
				parsedLimit = 5000 // cap at 5000 to be respectful
			}
			limit = parsedLimit
		}
	}

	// Start scraping
	if appErr := h.service.ScrapeBostonArtists(r.Context(), limit); appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Get updated count
	count, appErr := h.service.GetArtistCount(r.Context())
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"message":       "Bandcamp scraping completed successfully",
		"artists_total": count,
		"limit_used":    limit,
	}

	writeJSON(w, response)
}

// GetScrapedArtists returns the scraped artists
func (h *Handler) GetScrapedArtists(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit := 50 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 200 {
				parsedLimit = 200 // cap at 200 for performance
			}
			limit = parsedLimit
		}
	}

	// Get artists
	artists, appErr := h.service.GetScrapedArtists(r.Context(), limit)
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	// Get total count
	totalCount, appErr := h.service.GetArtistCount(r.Context())
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	response := map[string]interface{}{
		"data": artists,
		"meta": map[string]interface{}{
			"count":       len(artists),
			"limit":       limit,
			"total_count": totalCount,
		},
	}

	writeJSON(w, response)
}

// GetArtistCount returns just the count of scraped artists
func (h *Handler) GetArtistCount(w http.ResponseWriter, r *http.Request) {
	count, appErr := h.service.GetArtistCount(r.Context())
	if appErr != nil {
		utils.HandleError(w, appErr)
		return
	}

	response := map[string]interface{}{
		"total_artists": count,
	}

	writeJSON(w, response)
}

// writeJSON is a helper function to write JSON responses
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
