package artists

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Handler to execute business logic for Artists Endpoint
*/
type Handler struct {
	service *Service
}

// CreateArtist handles the creation of a new artist
func (h *Handler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	var params CreateArtistParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	artist, err := h.service.CreateArtist(params)
	if err != nil {
		http.Error(w, "Failed to create artist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

// GetArtists handles retrieving all artists with pagination
func (h *Handler) GetArtists(w http.ResponseWriter, r *http.Request) {
	// Get page and limit from query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if pageVal, err := strconv.Atoi(pageStr); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if limitStr != "" {
		if limitVal, err := strconv.Atoi(limitStr); err == nil && limitVal > 0 {
			if limitVal > 100 {
				limitVal = 100 // Cap at 100
			}
			limit = limitVal
		}
	}

	artists, total, err := h.service.GetArtists(page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": artists,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetArtist handles retrieving a single artist by ID
func (h *Handler) GetArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artist, err := h.service.GetArtistByID(id)
	if err != nil {
		if err == primitive.ErrInvalidHex {
			http.Error(w, "Invalid artist ID format", http.StatusBadRequest)
		} else if err == mongo.ErrNoDocuments {
			http.Error(w, "Artist not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch artist", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

// UpdateArtist handles updating an existing artist
func (h *Handler) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	var params CreateArtistParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedArtist, err := h.service.UpdateArtist(id, params)
	if err != nil {
		http.Error(w, "Failed to update artist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedArtist)
}

// UpdatePartialArtist handles partial updates to an artist
func (h *Handler) UpdatePartialArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	var params CreateArtistParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedArtist, err := h.service.UpdatePartialArtist(id, params)
	if err != nil {
		http.Error(w, "Failed to update artist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedArtist)
}

// DeleteArtist handles removing an artist
func (h *Handler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteArtist(id); err != nil {
		http.Error(w, "Failed to delete artist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetArtistsByGenre handles retrieving artists by genre
func (h *Handler) GetArtistsByGenre(w http.ResponseWriter, r *http.Request) {
	genre := chi.URLParam(r, "genre")

	artists, err := h.service.GetAllArtistsByGenre(genre)
	if err != nil {
		http.Error(w, "Failed to fetch artists by genre", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  artists,
		"genre": genre,
	})
}

// GetArtistsByCity handles retrieving artists by city
func (h *Handler) GetArtistsByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	artists, err := h.service.GetArtistsByCity(city)
	if err != nil {
		http.Error(w, "Failed to fetch artists by city", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": artists,
		"city": city,
	})
}

// GetRecommendations handles general artist recommendations
func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := h.service.GetRecommendations()
	if err != nil {
		http.Error(w, "Failed to fetch recommendations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": recommendations,
	})
}

// GetRecommendationsByGenre handles genre-specific recommendations
func (h *Handler) GetRecommendationsByGenre(w http.ResponseWriter, r *http.Request) {
	genre := chi.URLParam(r, "genre")

	recommendations, err := h.service.GetRecommendationsByGenre(genre)
	if err != nil {
		http.Error(w, "Failed to fetch recommendations by genre", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  recommendations,
		"genre": genre,
	})
}

// GetRecommendationsByLocation handles location-specific recommendations
func (h *Handler) GetRecommendationsByLocation(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	recommendations, err := h.service.GetRecommendationsByLocation(city)
	if err != nil {
		http.Error(w, "Failed to fetch recommendations by city", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": recommendations,
		"city": city,
	})
}

// // SaveUserPreferences handles saving user listening preferences
// func (h *Handler) SaveUserPreferences(w http.ResponseWriter, r *http.Request) {
// 	var prefs struct {
// 		UserID   string   `json:"userId"`
// 		Genres   []string `json:"genres"`
// 		Artists  []string `json:"artists,omitempty"`
// 		Location string   `json:"location,omitempty"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	userID, err := primitive.ObjectIDFromHex(prefs.UserID)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Convert artist IDs from strings to ObjectIDs
// 	var artistIDs []primitive.ObjectID
// 	for _, artistIDStr := range prefs.Artists {
// 		artistID, err := primitive.ObjectIDFromHex(artistIDStr)
// 		if err != nil {
// 			http.Error(w, "Invalid artist ID in list", http.StatusBadRequest)
// 			return
// 		}
// 		artistIDs = append(artistIDs, artistID)
// 	}

// 	userPrefs, err := h.service.SaveUserPreferences(userID, prefs.Genres, artistIDs, prefs.Location)
// 	if err != nil {
// 		http.Error(w, "Failed to save user preferences", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(userPrefs)
// }

// // GetUserPreferences handles retrieving user listening preferences
// func (h *Handler) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
// 	userIDStr := chi.URLParam(r, "userId")
// 	userID, err := primitive.ObjectIDFromHex(userIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	prefs, err := h.service.GetUserPreferences(userID)
// 	if err != nil {
// 		if err == primitive.ErrInvalidHex {
// 			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
// 		} else if err == mongo.ErrNoDocuments {
// 			http.Error(w, "User preferences not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Failed to fetch user preferences", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(prefs)
// }

// // GetPersonalizedRecommendations handles user-specific recommendations
// func (h *Handler) GetPersonalizedRecommendations(w http.ResponseWriter, r *http.Request) {
// 	userIDStr := chi.URLParam(r, "userId")
// 	userID, err := primitive.ObjectIDFromHex(userIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	recommendations, err := h.service.GetPersonalizedRecommendations(userID)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch personalized recommendations", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"data":   recommendations,
// 		"userId": userIDStr,
// 	})
// }
