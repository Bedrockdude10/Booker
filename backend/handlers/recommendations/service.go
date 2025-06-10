// handlers/recommendations/service.go
package recommendations

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Bedrockdude10/Booker/backend/cache"
	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewService creates a new recommendations service
func NewService(collections map[string]*mongo.Collection) *Service {
	return &Service{
		artists:       collections["artists"],
		preferences:   collections["userPreferences"],
		interactions:  collections["userInteractions"],
		trendingCache: collections["trendingCache"],
	}
}

//==============================================================================
// Main Recommendation Methods
//==============================================================================

// GetPersonalizedRecommendations generates recommendations for a specific user
func (s *Service) GetPersonalizedRecommendations(ctx context.Context, userID primitive.ObjectID, limit int) (*RecommendationResponse, *utils.AppError) {
	if limit <= 0 {
		limit = 10
	}

	// Get user preferences
	prefs, err := s.getUserPreferences(ctx, userID)
	if err != nil {
		// If no preferences, return general recommendations
		return s.GetGeneralRecommendations(ctx, limit)
	}

	// Get user interactions to exclude already seen artists
	interactions, _ := s.getUserInteractions(ctx, userID, 100) // Get recent interactions
	excludeArtists := make([]primitive.ObjectID, 0)
	for _, interaction := range interactions {
		excludeArtists = append(excludeArtists, interaction.ArtistID)
	}

	// Generate recommendations based on preferences
	params := RecommendationParams{
		UserID:  userID,
		Genres:  prefs.PreferredGenres,
		Cities:  prefs.PreferredCities,
		Limit:   limit * 2, // Get more to filter and rank
		Exclude: excludeArtists,
	}

	recommendations, appErr := s.generateRecommendations(ctx, params)
	if appErr != nil {
		return nil, appErr
	}

	// Score and rank recommendations
	scoredRecommendations := s.scoreRecommendations(ctx, recommendations, prefs, interactions)

	// Limit results
	if len(scoredRecommendations) > limit {
		scoredRecommendations = scoredRecommendations[:limit]
	}

	return &RecommendationResponse{
		Data:        scoredRecommendations,
		Total:       len(scoredRecommendations),
		RequestedBy: "user",
		Metadata: map[string]interface{}{
			"userId":  userID.Hex(),
			"basedOn": "preferences",
		},
	}, nil
}

// GetRecommendationsByGenre returns recommendations for a specific genre
func (s *Service) GetRecommendationsByGenre(ctx context.Context, genre string, limit int) (*RecommendationResponse, *utils.AppError) {
	if !domain.HasGenre(genre) {
		return nil, utils.ValidationError("Invalid genre")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("recs:genre:%s:limit:%d", genre, limit)
	if cached, found := cache.Get(cacheKey); found {
		if response, ok := cached.(*RecommendationResponse); ok {
			return response, nil
		}
	}

	// Get artists for this genre
	artists, appErr := s.getArtistsByGenre(ctx, genre, limit*2)
	if appErr != nil {
		return nil, appErr
	}

	// Convert to recommendation results
	recommendations := make([]RecommendationResult, 0, len(artists))
	for _, artist := range artists {
		recommendations = append(recommendations, RecommendationResult{
			Artist: artist,
			Score:  1.0, // Base score for genre match
		})
	}

	// Add trending boost
	recommendations = s.addTrendingBoost(ctx, recommendations)

	// Sort by score and limit
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	response := &RecommendationResponse{
		Data:        recommendations,
		Total:       len(recommendations),
		RequestedBy: "genre",
		Metadata: map[string]interface{}{
			"genre": genre,
		},
	}

	// Cache for 30 minutes
	cache.Set(cacheKey, response, 30*time.Minute)

	return response, nil
}

// GetRecommendationsByCity returns recommendations for a specific city
func (s *Service) GetRecommendationsByCity(ctx context.Context, city string, limit int) (*RecommendationResponse, *utils.AppError) {
	// Check cache first
	cacheKey := fmt.Sprintf("recs:city:%s:limit:%d", city, limit)
	if cached, found := cache.Get(cacheKey); found {
		if response, ok := cached.(*RecommendationResponse); ok {
			return response, nil
		}
	}

	// Get artists for this city
	artists, appErr := s.getArtistsByCity(ctx, city, limit*2)
	if appErr != nil {
		return nil, appErr
	}

	// Convert to recommendation results
	recommendations := make([]RecommendationResult, 0, len(artists))
	for _, artist := range artists {
		recommendations = append(recommendations, RecommendationResult{
			Artist: artist,
			Score:  1.0, // Base score for city match
		})
	}

	// Add trending boost
	recommendations = s.addTrendingBoost(ctx, recommendations)

	// Sort by score and limit
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	response := &RecommendationResponse{
		Data:        recommendations,
		Total:       len(recommendations),
		RequestedBy: "city",
		Metadata: map[string]interface{}{
			"city": city,
		},
	}

	// Cache for 30 minutes
	cache.Set(cacheKey, response, 30*time.Minute)

	return response, nil
}

// GetGeneralRecommendations returns general recommendations (trending, popular)
func (s *Service) GetGeneralRecommendations(ctx context.Context, limit int) (*RecommendationResponse, *utils.AppError) {
	// Check cache first
	cacheKey := fmt.Sprintf("recs:general:limit:%d", limit)
	if cached, found := cache.Get(cacheKey); found {
		if response, ok := cached.(*RecommendationResponse); ok {
			return response, nil
		}
	}

	// Get trending artists
	trending, _ := s.getTrendingArtists(ctx, limit/2)

	// Get random selection of artists
	random, appErr := s.getRandomArtists(ctx, limit-len(trending))
	if appErr != nil {
		return nil, appErr
	}

	// Combine recommendations
	recommendations := make([]RecommendationResult, 0, len(trending)+len(random))

	// Add trending with higher scores
	for i, artist := range trending {
		score := 1.0 + (0.5 * float64(len(trending)-i) / float64(len(trending))) // Boost trending
		recommendations = append(recommendations, RecommendationResult{
			Artist: artist,
			Score:  score,
		})
	}

	// Add random artists
	for _, artist := range random {
		recommendations = append(recommendations, RecommendationResult{
			Artist: artist,
			Score:  0.5, // Lower score for random
		})
	}

	// Sort by score
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	response := &RecommendationResponse{
		Data:        recommendations,
		Total:       len(recommendations),
		RequestedBy: "general",
		Metadata: map[string]interface{}{
			"type": "trending_and_discovery",
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

	if _, err := s.interactions.InsertOne(ctx, interaction); err != nil {
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
// Helper Methods
//==============================================================================

// getUserPreferences retrieves user preferences
func (s *Service) getUserPreferences(ctx context.Context, userID primitive.ObjectID) (*UserPreference, *utils.AppError) {
	var prefs UserPreference
	err := s.preferences.FindOne(ctx, bson.M{"accountId": userID}).Decode(&prefs)

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

	cursor, err := s.interactions.Find(ctx, bson.M{"userId": userID}, opts)
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

// generateRecommendations generates basic recommendations based on parameters
func (s *Service) generateRecommendations(ctx context.Context, params RecommendationParams) ([]ArtistDocument, *utils.AppError) {
	filter := bson.M{}

	// Build filter based on parameters
	if len(params.Genres) > 0 || len(params.Cities) > 0 {
		andConditions := []bson.M{}

		if len(params.Genres) > 0 {
			andConditions = append(andConditions, bson.M{"genres": bson.M{"$in": params.Genres}})
		}

		if len(params.Cities) > 0 {
			andConditions = append(andConditions, bson.M{"cities": bson.M{"$in": params.Cities}})
		}

		filter["$or"] = andConditions
	}

	// Exclude specific artists
	if len(params.Exclude) > 0 {
		filter["_id"] = bson.M{"$nin": params.Exclude}
	}

	// Set up find options
	opts := options.Find()
	if params.Limit > 0 {
		opts.SetLimit(int64(params.Limit))
	}

	cursor, err := s.artists.Find(ctx, filter, opts)
	if err != nil {
		return nil, utils.DatabaseError("find artists for recommendations", err)
	}
	defer cursor.Close(ctx)

	var artists []ArtistDocument
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, utils.DatabaseError("decode recommended artists", err)
	}

	return artists, nil
}

// scoreRecommendations scores and ranks recommendations
func (s *Service) scoreRecommendations(ctx context.Context, artists []ArtistDocument, prefs *UserPreference, interactions []UserInteraction) []RecommendationResult {
	results := make([]RecommendationResult, 0, len(artists))

	for _, artist := range artists {
		score := s.calculateArtistScore(artist, prefs, interactions)
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

// calculateArtistScore calculates a score for an artist based on user preferences
func (s *Service) calculateArtistScore(artist ArtistDocument, prefs *UserPreference, interactions []UserInteraction) float64 {
	score := 0.0

	// Genre matching
	genreMatches := 0
	for _, genre := range artist.Genres {
		for _, prefGenre := range prefs.PreferredGenres {
			if genre == prefGenre {
				genreMatches++
				break
			}
		}
	}
	score += float64(genreMatches) * 0.6 // 60% weight for genre matching

	// City matching
	cityMatches := 0
	for _, city := range artist.Cities {
		for _, prefCity := range prefs.PreferredCities {
			if city == prefCity {
				cityMatches++
				break
			}
		}
	}
	score += float64(cityMatches) * 0.4 // 40% weight for city matching

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

// Cache invalidation helpers
func (s *Service) invalidateUserCaches(userID primitive.ObjectID) {
	cache.Del(fmt.Sprintf("user:recs:%s", userID.Hex()))
}

func (s *Service) invalidateTrendingCaches() {
	// This would invalidate trending-related caches
	// Implementation depends on your caching strategy
}

// Additional helper methods for getting artists
func (s *Service) getArtistsByGenre(ctx context.Context, genre string, limit int) ([]ArtistDocument, *utils.AppError) {
	opts := options.Find().SetLimit(int64(limit))
	cursor, err := s.artists.Find(ctx, bson.M{"genres": genre}, opts)
	if err != nil {
		return nil, utils.DatabaseError("find artists by genre", err)
	}
	defer cursor.Close(ctx)

	var artists []ArtistDocument
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, utils.DatabaseError("decode artists by genre", err)
	}

	return artists, nil
}

func (s *Service) getArtistsByCity(ctx context.Context, city string, limit int) ([]ArtistDocument, *utils.AppError) {
	opts := options.Find().SetLimit(int64(limit))
	cursor, err := s.artists.Find(ctx, bson.M{"cities": city}, opts)
	if err != nil {
		return nil, utils.DatabaseError("find artists by city", err)
	}
	defer cursor.Close(ctx)

	var artists []ArtistDocument
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, utils.DatabaseError("decode artists by city", err)
	}

	return artists, nil
}

func (s *Service) getRandomArtists(ctx context.Context, limit int) ([]ArtistDocument, *utils.AppError) {
	// Simple random selection - in production you might want to use MongoDB's $sample
	opts := options.Find().SetLimit(int64(limit * 2)) // Get more to randomize
	cursor, err := s.artists.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, utils.DatabaseError("find random artists", err)
	}
	defer cursor.Close(ctx)

	var artists []ArtistDocument
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, utils.DatabaseError("decode random artists", err)
	}

	// Simple shuffle and limit
	if len(artists) > limit {
		// Basic shuffle - use crypto/rand for production
		for i := range artists {
			j := i % len(artists)
			artists[i], artists[j] = artists[j], artists[i]
		}
		artists = artists[:limit]
	}

	return artists, nil
}

func (s *Service) getTrendingArtists(ctx context.Context, limit int) ([]ArtistDocument, *utils.AppError) {
	// For now, return empty slice - implement trending logic later
	return []ArtistDocument{}, nil
}

func (s *Service) addTrendingBoost(ctx context.Context, recommendations []RecommendationResult) []RecommendationResult {
	// For now, return as-is - implement trending boost later
	return recommendations
}
