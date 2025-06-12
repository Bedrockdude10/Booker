// domain/artists/filtering.go
package artists

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// FilterParams represents filtering options for artists
type FilterParams struct {
	Genres     []string `json:"genres,omitempty"`
	Cities     []string `json:"cities,omitempty"`
	MinRating  float64  `json:"minRating,omitempty"`
	MaxRating  float64  `json:"maxRating,omitempty"`
	HasManager *bool    `json:"hasManager,omitempty"`
	HasSpotify *bool    `json:"hasSpotify,omitempty"`
}

// ParseFilterParams extracts filter parameters from HTTP request
func ParseFilterParams(r *http.Request) FilterParams {
	params := FilterParams{}
	query := r.URL.Query()

	// Parse genres (comma-separated)
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

	// Parse cities (comma-separated)
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

// ValidateFilterParams validates the filter parameters
func ValidateFilterParams(filters FilterParams) *utils.AppError {
	// Validate genres
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

// BuildFilterQuery constructs MongoDB filter based on FilterParams
func BuildFilterQuery(filters FilterParams) bson.M {
	query := bson.M{}
	andConditions := []bson.M{}

	// Genre filtering (OR logic within genres)
	if len(filters.Genres) > 0 {
		andConditions = append(andConditions, bson.M{
			"genres": bson.M{"$in": filters.Genres},
		})
	}

	// City filtering (OR logic within cities)
	if len(filters.Cities) > 0 {
		andConditions = append(andConditions, bson.M{
			"cities": bson.M{"$in": filters.Cities},
		})
	}

	// Rating filtering
	if filters.MinRating > 0 || filters.MaxRating > 0 {
		ratingQuery := bson.M{}
		if filters.MinRating > 0 {
			ratingQuery["$gte"] = filters.MinRating
		}
		if filters.MaxRating > 0 {
			ratingQuery["$lte"] = filters.MaxRating
		}
		andConditions = append(andConditions, bson.M{
			"rating": ratingQuery,
		})
	}

	// Manager filtering - updated to use nested structure
	if filters.HasManager != nil {
		if *filters.HasManager {
			andConditions = append(andConditions, bson.M{
				"contactInfo.manager": bson.M{"$exists": true, "$ne": ""},
			})
		} else {
			andConditions = append(andConditions, bson.M{
				"$or": []bson.M{
					{"contactInfo.manager": bson.M{"$exists": false}},
					{"contactInfo.manager": ""},
				},
			})
		}
	}

	// Spotify filtering - updated to use nested structure
	if filters.HasSpotify != nil {
		if *filters.HasSpotify {
			andConditions = append(andConditions, bson.M{
				"contactInfo.social.spotify": bson.M{"$exists": true, "$ne": ""},
			})
		} else {
			andConditions = append(andConditions, bson.M{
				"$or": []bson.M{
					{"contactInfo.social.spotify": bson.M{"$exists": false}},
					{"contactInfo.social.spotify": ""},
				},
			})
		}
	}

	// Combine all conditions
	if len(andConditions) > 0 {
		query["$and"] = andConditions
	}

	return query
}
