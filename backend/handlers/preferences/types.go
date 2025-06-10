// handlers/preferences/types.go
package preferences

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserPreference represents a user's music and location preferences
type UserPreference struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AccountID       primitive.ObjectID `bson:"accountID" json:"accountID" validate:"required"`
	PreferredGenres []string           `bson:"preferredGenres" json:"preferredGenres" validate:"required,min=1,validgenres"`
	PreferredCities []string           `bson:"preferredCities" json:"preferredCities" validate:"required,min=1"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// CreateUserPreferenceParams for creating new user preferences
type CreateUserPreferenceParams struct {
	AccountID       primitive.ObjectID `json:"accountID" validate:"required"`
	PreferredGenres []string           `json:"preferredGenres" validate:"required,min=1,validgenres"`
	PreferredCities []string           `json:"preferredCities" validate:"required,min=1"`
}

// UpdateUserPreferenceParams for updating user preferences
type UpdateUserPreferenceParams struct {
	PreferredGenres []string `json:"preferredGenres" validate:"omitempty,min=1,validgenres"`
	PreferredCities []string `json:"preferredCities" validate:"omitempty,min=1"`
}

// Service struct for user preferences operations
type Service struct {
	preferences *mongo.Collection
}
