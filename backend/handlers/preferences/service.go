// handlers/preferences/service.go
package preferences

import (
	"context"
	"fmt"
	"time"

	"github.com/Bedrockdude10/Booker/backend/cache"
	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/Bedrockdude10/Booker/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewService creates a new preferences service
func NewService(collections map[string]*mongo.Collection) *Service {
	return &Service{
		preferences: collections["preferences"],
	}
}

//==============================================================================
// Create User Preferences
//==============================================================================

func (s *Service) CreateUserPreference(ctx context.Context, params CreateUserPreferenceParams) (*UserPreference, *utils.AppError) {
	// Validate account ID
	if params.AccountID.IsZero() {
		return nil, utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	// Validate genres using domain package
	for _, genre := range params.PreferredGenres {
		if !domain.HasGenre(genre) {
			return nil, utils.ValidationErrorLog(ctx, "Invalid genre", fmt.Sprintf("Genre '%s' is not valid", genre))
		}
	}

	// Check if preferences already exist for this account
	existing, _ := s.GetUserPreferenceByAccountID(ctx, params.AccountID)
	if existing != nil {
		return nil, utils.ValidationErrorLog(ctx, "User preferences already exist for this account")
	}

	// Create preference document
	preference := UserPreference{
		ID:              primitive.NewObjectID(),
		AccountID:       params.AccountID,
		PreferredGenres: params.PreferredGenres,
		PreferredCities: params.PreferredCities,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Insert into database
	if _, err := s.preferences.InsertOne(ctx, preference); err != nil {
		return nil, utils.Log(ctx,
			utils.DatabaseError("create user preference", err),
			"Failed to create user preference",
			"account_id", params.AccountID.Hex(),
		)
	}

	// Invalidate cache for this account
	cache.Del(fmt.Sprintf("preferences:account:%s", params.AccountID.Hex()))

	return &preference, nil
}

//==============================================================================
// Read Operations (with caching)
//==============================================================================

// GetUserPreferenceByID retrieves user preference by ID
func (s *Service) GetUserPreferenceByID(ctx context.Context, id primitive.ObjectID) (*UserPreference, *utils.AppError) {
	if id.IsZero() {
		return nil, utils.ValidationErrorLog(ctx, "Invalid preference ID")
	}

	key := fmt.Sprintf("preferences:id:%s", id.Hex())

	// Try cache first
	if cached, found := cache.Get(key); found {
		if preference, ok := cached.(*UserPreference); ok {
			return preference, nil
		}
	}

	// Fetch from database
	var preference UserPreference
	err := s.preferences.FindOne(ctx, bson.M{"_id": id}).Decode(&preference)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "User preference")
	}
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find user preference by id", err)
	}

	// Cache for 30 minutes
	cache.Set(key, &preference, 30*time.Minute)

	return &preference, nil
}

// GetUserPreferenceByAccountID retrieves user preference by account ID
func (s *Service) GetUserPreferenceByAccountID(ctx context.Context, accountID primitive.ObjectID) (*UserPreference, *utils.AppError) {
	if accountID.IsZero() {
		return nil, utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	key := fmt.Sprintf("preferences:account:%s", accountID.Hex())

	// Try cache first
	if cached, found := cache.Get(key); found {
		if preference, ok := cached.(*UserPreference); ok {
			return preference, nil
		}
	}

	// Fetch from database
	var preference UserPreference
	err := s.preferences.FindOne(ctx, bson.M{"accountID": accountID}).Decode(&preference)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "User preference")
	}
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find user preference by account id", err)
	}

	// Cache for 30 minutes
	cache.Set(key, &preference, 30*time.Minute)

	return &preference, nil
}

// GetAllUserPreferences retrieves all user preferences with pagination
func (s *Service) GetAllUserPreferences(ctx context.Context, page, limit int) ([]UserPreference, *utils.AppError) {
	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Set up find options
	findOptions := options.Find().
		SetSort(bson.M{"createdAt": -1}). // Most recent first
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := s.preferences.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find user preferences", err)
	}
	defer cursor.Close(ctx)

	var results []UserPreference
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode user preferences", err)
	}

	return results, nil
}

//==============================================================================
// Update Operations
//==============================================================================

// UpdateUserPreference updates existing user preferences
func (s *Service) UpdateUserPreference(ctx context.Context, id primitive.ObjectID, params UpdateUserPreferenceParams) (*UserPreference, *utils.AppError) {
	if id.IsZero() {
		return nil, utils.ValidationErrorLog(ctx, "Invalid preference ID")
	}

	// Build update document dynamically
	updateFields := bson.M{
		"updatedAt": time.Now(),
	}

	// Validate and add genres if provided
	if len(params.PreferredGenres) > 0 {
		for _, genre := range params.PreferredGenres {
			if !domain.HasGenre(genre) {
				return nil, utils.ValidationErrorLog(ctx, "Invalid genre", fmt.Sprintf("Genre '%s' is not valid", genre))
			}
		}
		updateFields["preferredGenres"] = params.PreferredGenres
	}

	// Add cities if provided
	if len(params.PreferredCities) > 0 {
		updateFields["preferredCities"] = params.PreferredCities
	}

	// If no fields to update, just return the existing preference
	if len(updateFields) == 1 { // Only updatedAt
		return s.GetUserPreferenceByID(ctx, id)
	}

	// Use FindOneAndUpdate to get the updated document back
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedPreference UserPreference
	err := s.preferences.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
		opts,
	).Decode(&updatedPreference)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "User preference")
	}
	if err != nil {
		return nil, utils.Log(ctx,
			utils.DatabaseError("update user preference", err),
			"Failed to update user preference",
			"preference_id", id.Hex(),
		)
	}

	// Invalidate cache
	cache.Del(fmt.Sprintf("preferences:id:%s", id.Hex()))
	cache.Del(fmt.Sprintf("preferences:account:%s", updatedPreference.AccountID.Hex()))

	return &updatedPreference, nil
}

// UpdateUserPreferenceByAccountID updates user preferences by account ID
func (s *Service) UpdateUserPreferenceByAccountID(ctx context.Context, accountID primitive.ObjectID, params UpdateUserPreferenceParams) (*UserPreference, *utils.AppError) {
	if accountID.IsZero() {
		return nil, utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	// First get the existing preference to get its ID
	existing, appErr := s.GetUserPreferenceByAccountID(ctx, accountID)
	if appErr != nil {
		return nil, appErr
	}

	// Use the regular update method
	return s.UpdateUserPreference(ctx, existing.ID, params)
}

//==============================================================================
// Delete Operations
//==============================================================================

// DeleteUserPreference removes user preferences
func (s *Service) DeleteUserPreference(ctx context.Context, id primitive.ObjectID) *utils.AppError {
	if id.IsZero() {
		return utils.ValidationErrorLog(ctx, "Invalid preference ID")
	}

	// Get the preference first to know which account to invalidate cache for
	preference, appErr := s.GetUserPreferenceByID(ctx, id)
	if appErr != nil {
		return appErr
	}

	result, err := s.preferences.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return utils.DatabaseErrorLog(ctx, "delete user preference", err)
	}

	if result.DeletedCount == 0 {
		return utils.NotFoundLog(ctx, "User preference")
	}

	// Invalidate cache
	cache.Del(fmt.Sprintf("preferences:id:%s", id.Hex()))
	cache.Del(fmt.Sprintf("preferences:account:%s", preference.AccountID.Hex()))

	return nil
}

// DeleteUserPreferenceByAccountID removes preferences by account ID
func (s *Service) DeleteUserPreferenceByAccountID(ctx context.Context, accountID primitive.ObjectID) *utils.AppError {
	if accountID.IsZero() {
		return utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	result, err := s.preferences.DeleteOne(ctx, bson.M{"accountID": accountID})
	if err != nil {
		return utils.DatabaseErrorLog(ctx, "delete user preference by account", err)
	}

	if result.DeletedCount == 0 {
		return utils.NotFoundLog(ctx, "User preference")
	}

	// Invalidate cache
	cache.Del(fmt.Sprintf("preferences:account:%s", accountID.Hex()))

	return nil
}

//==============================================================================
// Analytics and Statistics
//==============================================================================

// GetPreferencesByGenre gets all users who prefer a specific genre
func (s *Service) GetPreferencesByGenre(ctx context.Context, genre string) ([]UserPreference, *utils.AppError) {
	// Validate genre
	if !domain.HasGenre(genre) {
		return nil, utils.ValidationErrorLog(ctx, "Invalid genre", fmt.Sprintf("Genre '%s' is not valid", genre))
	}

	key := fmt.Sprintf("preferences:genre:%s", genre)

	// Try cache first
	if cached, found := cache.Get(key); found {
		if preferences, ok := cached.([]UserPreference); ok {
			return preferences, nil
		}
	}

	// Fetch from database
	cursor, err := s.preferences.Find(ctx, bson.M{"preferredGenres": genre})
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find preferences by genre", err)
	}
	defer cursor.Close(ctx)

	var results []UserPreference
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode preferences by genre", err)
	}

	// Cache for 15 minutes
	cache.Set(key, results, 15*time.Minute)

	return results, nil
}

// GetPreferencesByCity gets all users who prefer a specific city
func (s *Service) GetPreferencesByCity(ctx context.Context, city string) ([]UserPreference, *utils.AppError) {
	if city == "" {
		return nil, utils.ValidationErrorLog(ctx, "City is required")
	}

	key := fmt.Sprintf("preferences:city:%s", city)

	// Try cache first
	if cached, found := cache.Get(key); found {
		if preferences, ok := cached.([]UserPreference); ok {
			return preferences, nil
		}
	}

	// Fetch from database
	cursor, err := s.preferences.Find(ctx, bson.M{"preferredCities": city})
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find preferences by city", err)
	}
	defer cursor.Close(ctx)

	var results []UserPreference
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode preferences by city", err)
	}

	// Cache for 15 minutes
	cache.Set(key, results, 15*time.Minute)

	return results, nil
}

// CountUserPreferences returns total count for pagination
func (s *Service) CountUserPreferences(ctx context.Context) (int64, *utils.AppError) {
	count, err := s.preferences.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, utils.DatabaseErrorLog(ctx, "count user preferences", err)
	}
	return count, nil
}

//==============================================================================
// Cache warming for performance
//==============================================================================

// WarmCache pre-loads popular queries into cache
func (s *Service) WarmCache(ctx context.Context) {
	// Popular genres to warm cache for
	popularGenres := []string{"rock", "pop", "hip-hop", "electronic", "jazz", "indie"}
	popularCities := []string{"Nashville", "Los Angeles", "New York", "Austin", "Chicago"}

	// Pre-load popular genre preferences
	for _, genre := range popularGenres {
		go s.GetPreferencesByGenre(ctx, genre) // Fire and forget
	}

	// Pre-load popular city preferences
	for _, city := range popularCities {
		go s.GetPreferencesByCity(ctx, city) // Fire and forget
	}
}
