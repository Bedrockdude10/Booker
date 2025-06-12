// handlers/recommendations/service.go - Updated to use service composition
package recommendations

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Bedrockdude10/Booker/backend/cache"
	"github.com/Bedrockdude10/Booker/backend/domain/artists"
	artistsService "github.com/Bedrockdude10/Booker/backend/handlers/artists"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewService creates a new recommendations service with composed artists service
func NewService(collections map[string]*mongo.Collection) *Service {
	// Create artists service for composition
	artistsSvc := artistsService.NewService(collections)

	return &Service{
		artistsService:   artistsSvc,
		preferencesCol:   collections["userPreferences"],
		interactionsCol:  collections["userInteractions"],
		trendingCacheCol: collections["trendingCache"],
	}
}

//==============================================================================
// Main Recommendation Methods - Using Service Composition
//==============================================================================

// GetPersonalizedRecommendations generates personalized recommendations with filtering
func (s *Service) GetPersonalizedRecommendations(ctx context.Context, params EnhancedRecommendationParams) (*RecommendationResponse, *utils.AppError) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Get user preferences
	prefs, err := s.getUserPreferences(ctx, params.UserID)
	if err != nil {
		// If no preferences, use general recommendations with filters
		return s.GetGeneralRecommendations(ctx, params)
	}

	// Merge user preferences with explicit filters
	mergedFilters := s.mergeUserPreferencesWithFilters(prefs, params.Filters)

	// Get user interactions to exclude already seen artists
	interactions, _ := s.getUserInteractions(ctx, params.UserID, 100)
	excludeArtists := make([]primitive.ObjectID, 0)
	for _, interaction := range interactions {
		excludeArtists = append(excludeArtists, interaction.ArtistID)
	}

	// Use artists service to get filtered results
	// Get more than needed to account for exclusions
	rawArtists, appErr := s.artistsService.GetArtists(ctx, mergedFilters, params.Limit*2, params.Offset)
	if appErr != nil {
		return nil, appErr
	}

	// Filter out excluded artists
	filteredArtists := s.excludeInteractedArtists(rawArtists, excludeArtists)

	// Score based on user preferences and interactions
	personalizedResults := s.scorePersonalizedRecommendations(ctx, filteredArtists, prefs, interactions, params.Filters)

	// Limit results
	if len(personalizedResults) > params.Limit {
		personalizedResults = personalizedResults[:params.Limit]
	}

	return &RecommendationResponse{
		Data:        personalizedResults,
		Total:       len(personalizedResults),
		RequestedBy: "user",
		HasMore:     len(personalizedResults) == params.Limit,
		Metadata: map[string]interface{}{
			"userId":  params.UserID.Hex(),
			"basedOn": "preferences_and_filters",
			"filters": params.Filters,
			"preferences": map[string]interface{}{
				"genres": prefs.PreferredGenres,
				"cities": prefs.PreferredCities,
			},
		},
	}, nil
}

// GetRecommendationsByGenre returns recommendations for a specific genre with filtering
func (s *Service) GetRecommendationsByGenre(ctx context.Context, params EnhancedRecommendationParams) (*RecommendationResponse, *utils.AppError) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Validate that we have at least one genre
	if len(params.Filters.Genres) == 0 {
		return nil, utils.ValidationError("Genre filter is required")
	}

	// Validate filters using shared validation
	if appErr := artists.ValidateFilterParams(params.Filters); appErr != nil {
		return nil, appErr
	}

	// Check cache first
	cacheKey := fmt.Sprintf("recs:genre:%+v:filters:%+v", params.Filters.Genres, params.Filters)
	if cached, found := cache.Get(cacheKey); found {
		if response, ok := cached.(*RecommendationResponse); ok {
			return response, nil
		}
	}

	// Use artists service to get filtered artists
	rawArtists, appErr := s.artistsService.GetArtists(ctx, params.Filters, params.Limit*2, params.Offset)
	if appErr != nil {
		return nil, appErr
	}

	// Convert to recommendation results and score
	recommendations := s.scoreArtistsForRecommendations(rawArtists, params.Filters)

	// Add trending boost and sort
	recommendations = s.addTrendingBoost(ctx, recommendations)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if len(recommendations) > params.Limit {
		recommendations = recommendations[:params.Limit]
	}

	response := &RecommendationResponse{
		Data:        recommendations,
		Total:       len(recommendations),
		RequestedBy: "genre",
		HasMore:     len(recommendations) == params.Limit,
		Metadata: map[string]interface{}{
			"genres":  params.Filters.Genres,
			"filters": params.Filters,
		},
	}

	// Cache for 30 minutes
	cache.Set(cacheKey, response, 30*time.Minute)
	return response, nil
}

// GetRecommendationsByCity returns recommendations for a specific city with filtering
func (s *Service) GetRecommendationsByCity(ctx context.Context, params EnhancedRecommendationParams) (*RecommendationResponse, *utils.AppError) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Validate that we have at least one city
	if len(params.Filters.Cities) == 0 {
		return nil, utils.ValidationError("City filter is required")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("recs:city:%+v:filters:%+v", params.Filters.Cities, params.Filters)
	if cached, found := cache.Get(cacheKey); found {
		if response, ok := cached.(*RecommendationResponse); ok {
			return response, nil
		}
	}

	// Use artists service to get filtered artists
	rawArtists, appErr := s.artistsService.GetArtists(ctx, params.Filters, params.Limit*2, params.Offset)
	if appErr != nil {
		return nil, appErr
	}

	// Convert to recommendation results and score
	recommendations := s.scoreArtistsForRecommendations(rawArtists, params.Filters)

	// Add trending boost and sort
	recommendations = s.addTrendingBoost(ctx, recommendations)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if len(recommendations) > params.Limit {
		recommendations = recommendations[:params.Limit]
	}

	response := &RecommendationResponse{
		Data:        recommendations,
		Total:       len(recommendations),
		RequestedBy: "city",
		HasMore:     len(recommendations) == params.Limit,
		Metadata: map[string]interface{}{
			"cities":  params.Filters.Cities,
			"filters": params.Filters,
		},
	}

	// Cache for 30 minutes
	cache.Set(cacheKey, response, 30*time.Minute)
	return response, nil
}

// GetGeneralRecommendations returns general recommendations with filtering
func (s *Service) GetGeneralRecommendations(ctx context.Context, params EnhancedRecommendationParams) (*RecommendationResponse, *utils.AppError) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Check cache first
	cacheKey := fmt.Sprintf("recs:general:filters:%+v:limit:%d", params.Filters, params.Limit)
	if cached, found := cache.Get(cacheKey); found {
		if response, ok := cached.(*RecommendationResponse); ok {
			return response, nil
		}
	}

	// Use artists service to get filtered artists
	rawArtists, appErr := s.artistsService.GetArtists(ctx, params.Filters, params.Limit*2, params.Offset)
	if appErr != nil {
		return nil, appErr
	}

	// Convert to recommendation results and score
	recommendations := s.scoreArtistsForRecommendations(rawArtists, params.Filters)

	// Add trending boost and sort
	recommendations = s.addTrendingBoost(ctx, recommendations)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if len(recommendations) > params.Limit {
		recommendations = recommendations[:params.Limit]
	}

	response := &RecommendationResponse{
		Data:        recommendations,
		Total:       len(recommendations),
		RequestedBy: "general",
		HasMore:     len(recommendations) == params.Limit,
		Metadata: map[string]interface{}{
			"type":    "discovery",
			"filters": params.Filters,
		},
	}

	// Cache for 15 minutes
	cache.Set(cacheKey, response, 15*time.Minute)
	return response, nil
}

//==============================================================================
// User Interaction Methods
//==============================================================================

// TrackInteraction logs a user interaction with an artist
func (s *Service) TrackInteraction(ctx context.Context, params TrackInteractionParams) *utils.AppError {
	interaction := UserInteraction{
		ID:        primitive.NewObjectID(),
		UserID:    params.UserID,
		ArtistID:  params.ArtistID,
		Type:      params.Type,
		Timestamp: time.Now(),
		Metadata:  params.Metadata,
	}

	if _, err := s.interactionsCol.InsertOne(ctx, interaction); err != nil {
		return utils.DatabaseErrorLog(ctx, "track interaction", err)
	}

	// Invalidate relevant caches
	s.invalidateUserCaches(params.UserID)
	s.invalidateTrendingCaches()

	return nil
}

// GetUserInteractions retrieves recent interactions for a user
func (s *Service) GetUserInteractions(ctx context.Context, userID primitive.ObjectID, limit int) ([]UserInteraction, *utils.AppError) {
	return s.getUserInteractions(ctx, userID, limit)
}

//==============================================================================
// Recommendation Scoring Methods (Recommendation-Specific Logic)
//==============================================================================

// scoreArtistsForRecommendations converts artists to recommendation results with scoring
func (s *Service) scoreArtistsForRecommendations(rawArtists []artists.ArtistDocument, filters artists.FilterParams) []RecommendationResult {
	recommendations := make([]RecommendationResult, 0, len(rawArtists))

	for _, artist := range rawArtists {
		score := s.calculateFilteredScore(artist, filters)
		recommendations = append(recommendations, RecommendationResult{
			Artist: artist,
			Score:  score,
		})
	}

	return recommendations
}

// excludeInteractedArtists removes artists the user has already interacted with
func (s *Service) excludeInteractedArtists(artists []artists.ArtistDocument, exclude []primitive.ObjectID) []artists.ArtistDocument {
	if len(exclude) == 0 {
		return artists
	}

	excludeSet := make(map[primitive.ObjectID]bool)
	for _, id := range exclude {
		excludeSet[id] = true
	}

	filtered := make([]artistsService.ArtistDocument, 0, len(artists))
	for _, artist := range artists {
		if !excludeSet[artist.ID] {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}

// mergeUserPreferencesWithFilters combines user preferences with explicit filters
func (s *Service) mergeUserPreferencesWithFilters(prefs *UserPreferenceAlias, filters artists.FilterParams) artists.FilterParams {
	merged := filters // Start with explicit filters

	// Add user preferred genres if no explicit genres provided
	if len(filters.Genres) == 0 && len(prefs.PreferredGenres) > 0 {
		merged.Genres = prefs.PreferredGenres
	}

	// Add user preferred cities if no explicit cities provided
	if len(filters.Cities) == 0 && len(prefs.PreferredCities) > 0 {
		merged.Cities = prefs.PreferredCities
	}

	return merged
}

// calculateFilteredScore calculates score based on how well an artist matches filters
func (s *Service) calculateFilteredScore(artist artists.ArtistDocument, filters artists.FilterParams) float64 {
	score := 1.0 // Base score

	// Genre match boost
	if len(filters.Genres) > 0 {
		genreMatches := 0
		for _, artistGenre := range artist.Genres {
			for _, filterGenre := range filters.Genres {
				if artistGenre == filterGenre {
					genreMatches++
					break
				}
			}
		}
		score += float64(genreMatches) * 0.3
	}

	// City match boost
	if len(filters.Cities) > 0 {
		cityMatches := 0
		for _, artistCity := range artist.Cities {
			for _, filterCity := range filters.Cities {
				if artistCity == filterCity {
					cityMatches++
					break
				}
			}
		}
		score += float64(cityMatches) * 0.2
	}

	// Manager boost - using new ContactInfo structure
	if filters.HasManager != nil && *filters.HasManager && artist.ContactInfo.Manager != "" {
		score += 0.1
	}

	// Spotify boost - using new ContactInfo structure
	if filters.HasSpotify != nil && *filters.HasSpotify && artist.ContactInfo.Social.Spotify != "" {
		score += 0.1
	}

	return score
}

// scorePersonalizedRecommendations scores recommendations based on user preferences + filters
func (s *Service) scorePersonalizedRecommendations(ctx context.Context, artists []artists.ArtistDocument, prefs *UserPreferenceAlias, interactions []UserInteraction, filters artists.FilterParams) []RecommendationResult {
	results := make([]RecommendationResult, 0, len(artists))

	for _, artist := range artists {
		// Start with filter-based score
		score := s.calculateFilteredScore(artist, filters)

		// Add personalization boost
		score += s.calculatePersonalizationScore(artist, prefs, interactions)

		results = append(results, RecommendationResult{
			Artist: artist,
			Score:  score,
		})
	}

	// Sort by score (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// calculatePersonalizationScore calculates additional score based on user preferences
func (s *Service) calculatePersonalizationScore(artist artists.ArtistDocument, prefs *UserPreferenceAlias, interactions []UserInteraction) float64 {
	score := 0.0

	// Preference-based scoring
	genreMatches := 0
	for _, genre := range artist.Genres {
		for _, prefGenre := range prefs.PreferredGenres {
			if genre == prefGenre {
				genreMatches++
				break
			}
		}
	}
	score += float64(genreMatches) * 0.4 // 40% weight for preferred genres

	cityMatches := 0
	for _, city := range artist.Cities {
		for _, prefCity := range prefs.PreferredCities {
			if city == prefCity {
				cityMatches++
				break
			}
		}
	}
	score += float64(cityMatches) * 0.3 // 30% weight for preferred cities

	// Favorite artists boost
	for _, favArtist := range prefs.FavoriteArtists {
		if favArtist == artist.ID {
			score += 1.0 // Big boost for favorites
			break
		}
	}

	// Interaction history penalty (to avoid showing same artists)
	for _, interaction := range interactions {
		if interaction.ArtistID == artist.ID {
			switch interaction.Type {
			case InteractionSkip:
				score -= 0.3 // Reduce score for skipped artists
			case InteractionView:
				score -= 0.1 // Small penalty for already viewed
			}
		}
	}

	return math.Max(0, score) // Ensure non-negative score
}

//==============================================================================
// Helper Methods (Database Access)
//==============================================================================

// getUserPreferences retrieves user preferences
func (s *Service) getUserPreferences(ctx context.Context, userID primitive.ObjectID) (*UserPreferenceAlias, *utils.AppError) {
	var prefs UserPreferenceAlias
	err := s.preferencesCol.FindOne(ctx, bson.M{"accountId": userID}).Decode(&prefs)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFound("User preferences")
	}
	if err != nil {
		return nil, utils.DatabaseError("find user preferences", err)
	}

	return &prefs, nil
}

// getUserInteractions retrieves user interactions
func (s *Service) getUserInteractions(ctx context.Context, userID primitive.ObjectID, limit int) ([]UserInteraction, *utils.AppError) {
	opts := options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetLimit(int64(limit))

	cursor, err := s.interactionsCol.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, utils.DatabaseError("find user interactions", err)
	}
	defer cursor.Close(ctx)

	var interactions []UserInteraction
	if err := cursor.All(ctx, &interactions); err != nil {
		return nil, utils.DatabaseError("decode user interactions", err)
	}

	return interactions, nil
}

//==============================================================================
// Cache Management
//==============================================================================

// invalidateUserCaches invalidates user-specific caches
func (s *Service) invalidateUserCaches(userID primitive.ObjectID) {
	cache.Del(fmt.Sprintf("user:recs:%s", userID.Hex()))
}

// invalidateTrendingCaches invalidates trending-related caches
func (s *Service) invalidateTrendingCaches() {
	// This would invalidate trending-related caches
	// Implementation depends on your caching strategy
}

// addTrendingBoost adds trending boost to recommendations (placeholder)
func (s *Service) addTrendingBoost(ctx context.Context, recommendations []RecommendationResult) []RecommendationResult {
	// For now, return as-is - implement trending boost later
	return recommendations
}
