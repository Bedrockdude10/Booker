package artists

import (
	"context"

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

// GetArtists fetches all artist documents from MongoDB with pagination
func (s *Service) GetArtists(page int, limit int) ([]ArtistDocument, int64, error) {
	ctx := context.Background()

	skip := (page - 1) * limit

	totalCount, err := s.artists.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	findOptions.SetSort(bson.M{"name": 1}) // Sort by name alphabetically

	cursor, err := s.artists.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []ArtistDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// GetArtistByID returns a single artist document by its ObjectID
func (s *Service) GetArtistByID(id primitive.ObjectID) (*ArtistDocument, error) {
	ctx := context.Background()
	filter := bson.M{"_id": id}

	var artist ArtistDocument
	err := s.artists.FindOne(ctx, filter).Decode(&artist)

	if err == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	} else if err != nil {
		return nil, err
	}

	return &artist, nil
}

// GetAllArtistsByGenre retrieves artists by genre
func (s *Service) GetAllArtistsByGenre(genre string) ([]ArtistDocument, error) {
	ctx := context.Background()

	filter := bson.M{"genres": genre}

	cursor, err := s.artists.Find(ctx, filter)
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

// GetArtistsByCity retrieves artists by city
func (s *Service) GetArtistsByCity(city string) ([]ArtistDocument, error) {
	ctx := context.Background()

	filter := bson.M{"city": city}

	cursor, err := s.artists.Find(ctx, filter)
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

// CreateArtist adds a new artist document
func (s *Service) CreateArtist(params CreateArtistParams) (*ArtistDocument, error) {
	ctx := context.Background()

	artist := ArtistDocument{
		ID:        primitive.NewObjectID(),
		Name:      params.Name,
		Genres:    params.Genres,
		Manager:   params.Manager,
		Cities:    params.Cities,
		SpotifyID: params.SpotifyID,
	}

	_, err := s.artists.InsertOne(ctx, artist)
	if err != nil {
		return nil, err
	}

	return &artist, nil
}

// UpdateArtist updates an existing artist document
func (s *Service) UpdateArtist(id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, error) {
	ctx := context.Background()
	filter := bson.M{"_id": id}

	updateFields := bson.M{
		"name":      params.Name,
		"genres":    params.Genres,
		"manager":   params.Manager,
		"city":      params.Cities,
		"streaming": params.SpotifyID,
	}

	update := bson.M{"$set": updateFields}

	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	var updatedArtist ArtistDocument
	err := s.artists.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&updatedArtist)
	if err != nil {
		return nil, err
	}

	return &updatedArtist, nil
}

// UpdatePartialArtist partially updates an artist document
func (s *Service) UpdatePartialArtist(id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, error) {
	ctx := context.Background()
	filter := bson.M{"_id": id}

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
		updateFields["city"] = params.Cities
	}

	if params.SpotifyID != "" {
		updateFields["streaming"] = params.SpotifyID
	}

	// If no fields to update, return the existing document
	if len(updateFields) == 0 {
		return s.GetArtistByID(id)
	}

	update := bson.M{"$set": updateFields}

	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	var updatedArtist ArtistDocument
	err := s.artists.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&updatedArtist)
	if err != nil {
		return nil, err
	}

	return &updatedArtist, nil
}

// DeleteArtist removes an artist document
func (s *Service) DeleteArtist(id primitive.ObjectID) error {
	ctx := context.Background()
	filter := bson.M{"_id": id}

	_, err := s.artists.DeleteOne(ctx, filter)
	return err
}

// SaveUserPreferences saves a user's listening preferences
// func (s *Service) SaveUserPreferences(userID primitive.ObjectID, genres []string, artists []primitive.ObjectID, location string) (*UserPreference, error) {
// 	ctx := context.Background()

// 	// Check if user preferences already exist
// 	filter := bson.M{"userId": userID}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"genres":   genres,
// 			"artists":  artists,
// 			"location": location,
// 		},
// 	}

// 	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

// 	var userPref UserPreference
// 	err := s.userPreferences.FindOneAndUpdate(ctx, filter, update, opts).Decode(&userPref)
// 	if err != nil && err != mongo.ErrNoDocuments {
// 		return nil, err
// 	}

// 	// If it was an upsert and no document was returned, get the new document
// 	if err == mongo.ErrNoDocuments {
// 		err = s.userPreferences.FindOne(ctx, filter).Decode(&userPref)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &userPref, nil
// }

// // GetUserPreferences retrieves a user's listening preferences
// func (s *Service) GetUserPreferences(userID primitive.ObjectID) (*UserPreference, error) {
// 	ctx := context.Background()
// 	filter := bson.M{"userId": userID}

// 	var userPref UserPreference
// 	err := s.userPreferences.FindOne(ctx, filter).Decode(&userPref)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &userPref, nil
// }

// GetRecommendations returns general artist recommendations
func (s *Service) GetRecommendations() ([]ArtistDocument, error) {
	// This would be replaced with your recommendation algorithm
	// For now, just returning a limited set of artists
	ctx := context.Background()

	findOptions := options.Find()
	findOptions.SetLimit(10)

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

// // GetPersonalizedRecommendations returns artist recommendations for a specific user
// func (s *Service) GetPersonalizedRecommendations(userID primitive.ObjectID) ([]ArtistDocument, error) {
// 	// Get user preferences
// 	userPref, err := s.GetUserPreferences(userID)
// 	if err != nil {
// 		// If no preferences exist, return general recommendations
// 		if err == mongo.ErrNoDocuments {
// 			return s.GetRecommendations()
// 		}
// 		return nil, err
// 	}

// 	// Use preferences to create a more targeted query
// 	ctx := context.Background()

// 	// Find artists that match the user's preferred genres
// 	filter := bson.M{
// 		"genres": bson.M{
// 			"$in": userPref.Genres,
// 		},
// 	}

// 	// If user has a location preference, consider that too
// 	if userPref.Location != "" {
// 		filter["city"] = userPref.Location
// 	}

// 	// Don't recommend artists they already follow
// 	if len(userPref.Artists) > 0 {
// 		filter["_id"] = bson.M{
// 			"$nin": userPref.Artists,
// 		}
// 	}

// 	findOptions := options.Find()
// 	findOptions.SetLimit(10)

// 	cursor, err := s.artists.Find(ctx, filter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var results []ArtistDocument
// 	if err := cursor.All(ctx, &results); err != nil {
// 		return nil, err
// 	}

// 	return results, nil
// }

// GetRecommendationsByGenre returns artist recommendations for a specific genre
func (s *Service) GetRecommendationsByGenre(genre string) ([]ArtistDocument, error) {
	ctx := context.Background()

	filter := bson.M{"genres": genre}

	findOptions := options.Find()
	findOptions.SetLimit(10)

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

// GetRecommendationsByLocation returns artist recommendations for a specific location
func (s *Service) GetRecommendationsByLocation(city string) ([]ArtistDocument, error) {
	ctx := context.Background()

	filter := bson.M{"city": city}

	findOptions := options.Find()
	findOptions.SetLimit(10)

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
