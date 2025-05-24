// handlers/artists/service.go
package artists

import (
	"context"
	"os"
	"strconv"

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

// GetArtists - just improved error handling, same signature
func (s *Service) GetArtists(ctx context.Context, page int, limit int) ([]ArtistDocument, int64, *utils.AppError) {
	skip := (page - 1) * limit

	totalCount, err := s.artists.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, utils.DatabaseErrorLog(ctx, "count artists", err)
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(getDefaultSort())

	cursor, err := s.artists.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, utils.DatabaseErrorLog(ctx, "find artists", err)
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, utils.DatabaseErrorLog(ctx, "decode artists", err)
	}

	return results, totalCount, nil
}

// GetArtistByID - same signature, improved error handling
func (s *Service) GetArtistByID(ctx context.Context, id primitive.ObjectID) (*ArtistDocument, *utils.AppError) {
	var artist ArtistDocument
	err := s.artists.FindOne(ctx, bson.M{"_id": id}).Decode(&artist)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "Artist")
	}
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find artist by id", err)
	}

	return &artist, nil
}

// GetAllArtistsByGenre - same signature, improved error handling
func (s *Service) GetAllArtistsByGenre(ctx context.Context, genre string) ([]ArtistDocument, *utils.AppError) {
	if !HasGenre(genre) {
		return nil, utils.ValidationErrorLog(ctx, "Invalid genre", "Genre '"+genre+"' is not valid")
	}

	cursor, err := s.artists.Find(ctx, bson.M{"genres": genre})
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find artists by genre", err)
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode artists by genre", err)
	}

	return results, nil
}

// GetArtistsByCity - same signature, improved error handling
func (s *Service) GetArtistsByCity(ctx context.Context, city string) ([]ArtistDocument, *utils.AppError) {
	if city == "" {
		return nil, utils.ValidationErrorLog(ctx, "City is required")
	}

	cursor, err := s.artists.Find(ctx, bson.M{"cities": city})
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find artists by city", err)
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode artists by city", err)
	}

	return results, nil
}

// CreateArtist - same signature, improved error handling
func (s *Service) CreateArtist(ctx context.Context, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	// Validate input
	if params.Name == "" {
		return nil, utils.ValidationErrorLog(ctx, "Artist name is required")
	}
	if len(params.Cities) == 0 {
		return nil, utils.ValidationErrorLog(ctx, "At least one city is required")
	}
	if err := ValidateGenres(ctx, params.Genres); err != nil {
		return nil, err
	}

	artist := ArtistDocument{
		ID:        primitive.NewObjectID(),
		Name:      params.Name,
		Genres:    params.Genres,
		Manager:   params.Manager,
		Cities:    params.Cities,
		SpotifyID: params.SpotifyID,
	}

	if _, err := s.artists.InsertOne(ctx, artist); err != nil {
		return nil, utils.Log(ctx,
			utils.DatabaseError("create artist", err),
			"Failed to create artist",
			"artist_name", params.Name,
		)
	}

	return &artist, nil
}

// UpdateArtist - same signature, improved error handling
func (s *Service) UpdateArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	if len(params.Genres) > 0 {
		if err := ValidateGenres(ctx, params.Genres); err != nil {
			return nil, err
		}
	}

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

	return &updatedArtist, nil
}

// UpdatePartialArtist - same signature, improved error handling
func (s *Service) UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	updateFields := bson.M{}

	if params.Name != "" {
		updateFields["name"] = params.Name
	}
	if len(params.Genres) > 0 {
		if err := ValidateGenres(ctx, params.Genres); err != nil {
			return nil, err
		}
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

	return &updatedArtist, nil
}

// DeleteArtist - same signature, improved error handling
func (s *Service) DeleteArtist(ctx context.Context, id primitive.ObjectID) *utils.AppError {
	result, err := s.artists.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return utils.DatabaseErrorLog(ctx, "delete artist", err)
	}

	if result.DeletedCount == 0 {
		return utils.NotFoundLog(ctx, "Artist")
	}

	return nil
}

// GetRecommendations - updated to use environment variable for limit
func (s *Service) GetRecommendations() ([]ArtistDocument, error) {
	ctx := context.Background()
	findOptions := options.Find().SetLimit(int64(getRecommendationLimit()))

	cursor, err := s.artists.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetRecommendationsByGenre - updated to use environment variable for limit
func (s *Service) GetRecommendationsByGenre(genre string) ([]ArtistDocument, error) {
	if !HasGenre(genre) {
		return nil, utils.ValidationError("Invalid genre", "Genre '"+genre+"' is not valid")
	}

	ctx := context.Background()
	filter := bson.M{"genres": genre}
	findOptions := options.Find().SetLimit(int64(getRecommendationLimit()))

	cursor, err := s.artists.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetRecommendationsByLocation - updated to use environment variable for limit
func (s *Service) GetRecommendationsByLocation(city string) ([]ArtistDocument, error) {
	if city == "" {
		return nil, utils.ValidationError("City is required")
	}

	ctx := context.Background()
	filter := bson.M{"cities": city}
	findOptions := options.Find().SetLimit(int64(getRecommendationLimit()))

	cursor, err := s.artists.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// getRecommendationLimit returns the recommendation limit from environment
func getRecommendationLimit() int {
	if limitStr := os.Getenv("RECOMMENDATION_LIMIT"); limitStr != "" {
		if limitVal, err := strconv.Atoi(limitStr); err == nil && limitVal > 0 {
			return limitVal
		}
	}
	return 10 // fallback default
}

// getDefaultSort returns the default sort configuration from environment
func getDefaultSort() bson.M {
	sortField := os.Getenv("DEFAULT_SORT_FIELD")
	if sortField == "" {
		sortField = "name" // fallback default
	}
	return bson.M{sortField: 1} // 1 for ascending order
}
