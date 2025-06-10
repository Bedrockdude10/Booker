// handlers/recommendations/types.go
package recommendations

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RecommendationParams for generating recommendations
type RecommendationParams struct {
	UserID  primitive.ObjectID   `json:"userId,omitempty"`
	Genres  []string             `json:"genres,omitempty"`
	Cities  []string             `json:"cities,omitempty"`
	Limit   int                  `json:"limit,omitempty"`
	Exclude []primitive.ObjectID `json:"exclude,omitempty"` // Artists to exclude
}

// RecommendationResult represents a recommended artist with scoring
type RecommendationResult struct {
	Artist ArtistDocument `json:"artist"`
	Score  float64        `json:"score"`
}

// RecommendationResponse for API responses
type RecommendationResponse struct {
	Data        []RecommendationResult `json:"data"`
	Total       int                    `json:"total"`
	RequestedBy string                 `json:"requestedBy,omitempty"` // "user", "genre", "city"
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UserInteraction tracks user behavior for better recommendations
type UserInteraction struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID    primitive.ObjectID     `bson:"userId" json:"userId"`
	ArtistID  primitive.ObjectID     `bson:"artistId" json:"artistId"`
	Type      InteractionType        `bson:"type" json:"type"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// InteractionType represents different user actions
type InteractionType string

const (
	InteractionView    InteractionType = "view"
	InteractionLike    InteractionType = "like"
	InteractionSave    InteractionType = "save"
	InteractionContact InteractionType = "contact"
	InteractionSkip    InteractionType = "skip"
)

// TrendingCache stores pre-computed trending data
type TrendingCache struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Type       string               `bson:"type"` // "global", "genre:rock", "city:nashville"
	ArtistIDs  []primitive.ObjectID `bson:"artistIds"`
	Scores     []float64            `bson:"scores,omitempty"`
	ComputedAt time.Time            `bson:"computedAt"`
	ExpiresAt  time.Time            `bson:"expiresAt"`
}

// ArtistDocument mirrors the artist structure (import from artists package in real implementation)
type ArtistDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Genres    []string           `bson:"genres" json:"genres"`
	Manager   string             `bson:"manager,omitempty" json:"manager,omitempty"`
	Cities    []string           `bson:"cities" json:"cities"`
	SpotifyID string             `bson:"spotifyId,omitempty" json:"spotifyId,omitempty"`
}

// UserPreference mirrors the preference structure
type UserPreference struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	AccountID       primitive.ObjectID   `bson:"accountId" json:"accountId"`
	PreferredGenres []string             `bson:"preferredGenres" json:"preferredGenres"`
	PreferredCities []string             `bson:"preferredCities" json:"preferredCities"`
	FavoriteArtists []primitive.ObjectID `bson:"favoriteArtists,omitempty" json:"favoriteArtists,omitempty"`
	CreatedAt       time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// Service struct for recommendations
type Service struct {
	artists       *mongo.Collection
	preferences   *mongo.Collection
	interactions  *mongo.Collection
	trendingCache *mongo.Collection
}

// TrackInteractionParams for logging user interactions
type TrackInteractionParams struct {
	UserID   primitive.ObjectID     `json:"userId" validate:"required"`
	ArtistID primitive.ObjectID     `json:"artistId" validate:"required"`
	Type     InteractionType        `json:"type" validate:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
