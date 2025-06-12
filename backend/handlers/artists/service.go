// handlers/artists/service.go - Updated to use shared domain types
package artists

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Bedrockdude10/Booker/backend/cache"
	"github.com/Bedrockdude10/Booker/backend/domain/artists"
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

/////////////////////////////////////////////// CORE ARTIST OPERATIONS

// GetArtists retrieves a list of artists with optional filtering
// This is the main method used by recommendations service
func (s *Service) GetArtists(ctx context.Context, filters artists.FilterParams, limit, offset int) ([]artists.ArtistDocument, *utils.AppError) {
	// Use shared filtering logic from domain
	filterQuery := artists.BuildFilterQuery(filters)

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

	var results []artists.ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode artists", err)
	}

	return results, nil
}

// GetArtistByID retrieves a single artist by ID
func (s *Service) GetArtistByID(ctx context.Context, id primitive.ObjectID) (*artists.ArtistDocument, *utils.AppError) {
	key := fmt.Sprintf("artist:%s", id.Hex())

	// Try cache first
	if cached, found := cache.Get(key); found {
		if artist, ok := cached.(*artists.ArtistDocument); ok {
			return artist, nil
		}
	}

	var artist artists.ArtistDocument
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

/////////////////////////////////////////////// CRUD OPERATIONS FOR ADMIN

// CreateArtist creates a new artist
func (s *Service) CreateArtist(ctx context.Context, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError) {
	artist := artists.ArtistDocument{
		ID:          primitive.NewObjectID(),
		Name:        params.Name,
		Genres:      params.Genres,
		Cities:      params.Cities,
		ContactInfo: params.ContactInfo,
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

// UpdateArtist performs a full update of an artist
func (s *Service) UpdateArtist(ctx context.Context, id primitive.ObjectID, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError) {
	updateFields := bson.M{
		"name":        params.Name,
		"genres":      params.Genres,
		"cities":      params.Cities,
		"contactInfo": params.ContactInfo,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedArtist artists.ArtistDocument
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

// UpdatePartialArtist performs a partial update of an artist
func (s *Service) UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError) {
	updateFields := bson.M{}

	if params.Name != "" {
		updateFields["name"] = params.Name
	}
	if len(params.Genres) > 0 {
		updateFields["genres"] = params.Genres
	}
	if len(params.Cities) > 0 {
		updateFields["cities"] = params.Cities
	}

	// Handle ContactInfo updates - only update if provided
	if !isEmptyContactInfo(params.ContactInfo) {
		updateFields["contactInfo"] = params.ContactInfo
	}

	if len(updateFields) == 0 {
		return s.GetArtistByID(ctx, id)
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedArtist artists.ArtistDocument
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

// DeleteArtist deletes an artist
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

/////////////////////////////////////////////// HELPER FUNCTIONS (PRIVATE)

// Helper function to check if ContactInfo is empty
func isEmptyContactInfo(contactInfo artists.ContactInfo) bool {
	return contactInfo.Manager == "" &&
		contactInfo.ManagerInfo == "" &&
		contactInfo.BookingInfo == "" &&
		contactInfo.LabelName == "" &&
		contactInfo.LabelURL == "" &&
		isEmptySocialMediaLinks(contactInfo.Social)
}

// Helper function to check if SocialMediaLinks is empty
func isEmptySocialMediaLinks(social artists.SocialMediaLinks) bool {
	return social.Spotify == "" &&
		social.AppleMusic == "" &&
		social.Bandcamp == "" &&
		social.Instagram == "" &&
		social.YouTube == "" &&
		social.Facebook == "" &&
		social.Twitter == "" &&
		social.TikTok == "" &&
		social.Website == "" &&
		social.SoundCloud == "" &&
		social.Discogs == "" &&
		social.Beatport == "" &&
		social.Deezer == "" &&
		social.Pandora == "" &&
		social.Email == "" &&
		social.Phone == ""
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
}

// getDefaultSort returns the default sort configuration from environment
func getDefaultSort() bson.M {
	sortField := os.Getenv("DEFAULT_SORT_FIELD")
	if sortField == "" {
		sortField = "name" // fallback default
	}
	return bson.M{sortField: 1} // 1 for ascending order
}
