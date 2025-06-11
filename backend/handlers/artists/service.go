// handlers/artists/service.go - Replace old methods with filtering-first approach
package artists

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Bedrockdude10/Booker/backend/cache"
	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/Bedrockdude10/Booker/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewService receives the map of collections and initializes the service
func NewService(collections map[string]*mongo.Collection) *Service {
	return &Service{
		artists:         collections["artists"],
		userPreferences: collections["userPreferences"],
	}
}

/////////////////////////////////////////////// ARTISTS - ALL WITH FILTERING SUPPORT

// GetArtists retrieves a list of artists with optional filtering
func (s *Service) GetArtists(ctx context.Context, filters FilterParams, limit, offset int) ([]ArtistDocument, *utils.AppError) {
	// Build filter query using the function from filtering.go
	filterQuery := BuildFilterQuery(filters)

	// Set up find options
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(getDefaultSort())

	cursor, err := s.artists.Find(ctx, filterQuery, opts)
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find artists", err)
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode artists", err)
	}

	return results, nil
}

// GetArtistsByGenre - simplified to use filtering system
func (s *Service) GetArtistsByGenre(ctx context.Context, genre string, additionalFilters FilterParams) ([]ArtistDocument, *utils.AppError) {
	// Validate genre
	if !domain.HasGenre(genre) {
		return nil, utils.ValidationErrorLog(ctx, "Invalid genre", "Genre '"+genre+"' is not valid")
	}

	// Add genre to filters
	if !contains(additionalFilters.Genres, genre) {
		additionalFilters.Genres = append(additionalFilters.Genres, genre)
	}

	// Check cache
	cacheKey := fmt.Sprintf("artists:genre:%s:filters:%+v", genre, additionalFilters)
	if cached, found := cache.Get(cacheKey); found {
		if artists, ok := cached.([]ArtistDocument); ok {
			return artists, nil
		}
	}

	// Use unified filtering method
	artists, appErr := s.GetArtists(ctx, additionalFilters, 0, 0)
	if appErr != nil {
		return nil, appErr
	}

	// Cache for 15 minutes
	cache.Set(cacheKey, artists, 15*time.Minute)

	return artists, nil
}

// GetArtistsByCity - simplified to use filtering system
func (s *Service) GetArtistsByCity(ctx context.Context, city string, additionalFilters FilterParams) ([]ArtistDocument, *utils.AppError) {
	if city == "" {
		return nil, utils.ValidationErrorLog(ctx, "City is required")
	}

	// Add city to filters
	if !contains(additionalFilters.Cities, city) {
		additionalFilters.Cities = append(additionalFilters.Cities, city)
	}

	// Check cache
	cacheKey := fmt.Sprintf("artists:city:%s:filters:%+v", city, additionalFilters)
	if cached, found := cache.Get(cacheKey); found {
		if artists, ok := cached.([]ArtistDocument); ok {
			return artists, nil
		}
	}

	// Use unified filtering method
	artists, appErr := s.GetArtists(ctx, additionalFilters, 0, 0)
	if appErr != nil {
		return nil, appErr
	}

	// Cache for 10 minutes
	cache.Set(cacheKey, artists, 10*time.Minute)

	return artists, nil
}

// GetArtistByID - keep as is, no filtering needed for single artist lookup
func (s *Service) GetArtistByID(ctx context.Context, id primitive.ObjectID) (*ArtistDocument, *utils.AppError) {
	key := fmt.Sprintf("artist:%s", id.Hex())

	// Try cache
	if cached, found := cache.Get(key); found {
		if artist, ok := cached.(*ArtistDocument); ok {
			return artist, nil
		}
	}

	var artist ArtistDocument
	err := s.artists.FindOne(ctx, bson.M{"_id": id}).Decode(&artist)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "Artist")
	}
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find artist by id", err)
	}

	// Cache for 30 minutes
	cache.Set(key, &artist, 30*time.Minute)

	return &artist, nil
}

/////////////////////////////////////////////// CRUD OPERATIONS (unchanged)

// CreateArtist - keep as is
func (s *Service) CreateArtist(ctx context.Context, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	artist := ArtistDocument{
		ID:          primitive.NewObjectID(),
		Name:        params.Name,
		Genres:      params.Genres,
		Manager:     params.Manager,
		Cities:      params.Cities,
		SpotifyID:   params.SpotifyID,
		Rating:      0.0, // Default rating
		RatingCount: 0,   // No ratings yet
	}

	if _, err := s.artists.InsertOne(ctx, artist); err != nil {
		return nil, utils.Log(ctx,
			utils.DatabaseError("create artist", err),
			"Failed to create artist",
			"artist_name", params.Name,
		)
	}

	// Invalidate relevant caches
	s.invalidateFilterCaches(params.Genres, params.Cities)

	return &artist, nil
}

// UpdateArtist - keep as is but invalidate more caches
func (s *Service) UpdateArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	updateFields := bson.M{
		"name":      params.Name,
		"genres":    params.Genres,
		"manager":   params.Manager,
		"cities":    params.Cities,
		"spotifyId": params.SpotifyID,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedArtist ArtistDocument
	err := s.artists.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
		opts,
	).Decode(&updatedArtist)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "Artist")
	}
	if err != nil {
		return nil, utils.Log(ctx,
			utils.DatabaseError("update artist", err),
			"Failed to update artist",
			"artist_id", id.Hex(),
		)
	}

	// Invalidate caches
	s.invalidateFilterCaches(params.Genres, params.Cities)
	cache.Del(fmt.Sprintf("artist:%s", id.Hex()))

	return &updatedArtist, nil
}

// UpdatePartialArtist - keep as is but invalidate more caches
func (s *Service) UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	updateFields := bson.M{}

	if params.Name != "" {
		updateFields["name"] = params.Name
	}
	if len(params.Genres) > 0 {
		updateFields["genres"] = params.Genres
	}
	if params.Manager != "" {
		updateFields["manager"] = params.Manager
	}
	if len(params.Cities) > 0 {
		updateFields["cities"] = params.Cities
	}
	if params.SpotifyID != "" {
		updateFields["spotifyId"] = params.SpotifyID
	}

	if len(updateFields) == 0 {
		return s.GetArtistByID(ctx, id)
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedArtist ArtistDocument
	err := s.artists.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
		opts,
	).Decode(&updatedArtist)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "Artist")
	}
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "partial update artist", err)
	}

	// Invalidate caches
	s.invalidateFilterCaches(params.Genres, params.Cities)
	cache.Del(fmt.Sprintf("artist:%s", id.Hex()))

	return &updatedArtist, nil
}

// DeleteArtist - keep as is
func (s *Service) DeleteArtist(ctx context.Context, id primitive.ObjectID) *utils.AppError {
	result, err := s.artists.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return utils.DatabaseErrorLog(ctx, "delete artist", err)
	}

	if result.DeletedCount == 0 {
		return utils.NotFoundLog(ctx, "Artist")
	}

	// Invalidate caches
	cache.Del(fmt.Sprintf("artist:%s", id.Hex()))

	return nil
}

/////////////////////////////////////////////// HELPER FUNCTIONS

// Helper function to check if a string slice contains a value
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// invalidateFilterCaches invalidates caches that might be affected by genre/city changes
func (s *Service) invalidateFilterCaches(genres []string, cities []string) {
	// Invalidate genre-specific caches
	for _, genre := range genres {
		cache.Del(fmt.Sprintf("artists:genre:%s", genre))
	}

	// Invalidate city-specific caches
	for _, city := range cities {
		cache.Del(fmt.Sprintf("artists:city:%s", city))
	}

	// Could also invalidate general filter caches here if needed
}

// getDefaultSort returns the default sort configuration from environment
func getDefaultSort() bson.M {
	sortField := os.Getenv("DEFAULT_SORT_FIELD")
	if sortField == "" {
		sortField = "name" // fallback default
	}
	return bson.M{sortField: 1} // 1 for ascending order
}
