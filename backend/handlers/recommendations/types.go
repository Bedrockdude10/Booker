// handlers/recommendations/types.go
package recommendations

import (
	"time"

	"github.com/Bedrockdude10/Booker/backend/domain/artists"
	artistsHandler "github.com/Bedrockdude10/Booker/backend/handlers/artists"
	"github.com/Bedrockdude10/Booker/backend/handlers/preferences"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RecommendationResult represents a recommended artist with scoring
type RecommendationResult struct {
	Artist artists.ArtistDocument `json:"artist"`
	Score  float64                `json:"score"`
	Reason string                 `json:"reason,omitempty"` // Why this artist was recommended
}

// RecommendationResponse for API responses
type RecommendationResponse struct {
	Data        []RecommendationResult `json:"data"`
	Total       int                    `json:"total"`
	RequestedBy string                 `json:"requestedBy,omitempty"` // "user", "genre", "city", "general"
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	HasMore     bool                   `json:"hasMore,omitempty"`
}

// Enhanced RecommendationParams with filtering
type EnhancedRecommendationParams struct {
	UserID  primitive.ObjectID   `json:"userId,omitempty"`
	Filters artists.FilterParams `json:"filters,omitempty"` // Use shared filtering
	Limit   int                  `json:"limit,omitempty"`
	Offset  int                  `json:"offset,omitempty"`
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

// TrackInteractionParams for logging user interactions
type TrackInteractionParams struct {
	UserID   primitive.ObjectID     `json:"userId" validate:"required"`
	ArtistID primitive.ObjectID     `json:"artistId" validate:"required"`
	Type     InteractionType        `json:"type" validate:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TrendingCache stores pre-computed trending data
type TrendingCache struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Type       string               `bson:"type"` // "global", "genre:rock", "city:nashville"
	ArtistIDs  []primitive.ObjectID `bson:"artistIds"`
	Scores     []float64            `bson:"scores,omitempty"`
	ComputedAt time.Time            `bson:"computedAt"`
	ExpiresAt  time.Time            `bson:"expiresAt"`
}

// Service struct for recommendations - uses composition
type Service struct {
	artistsService   *artistsHandler.Service // Compose artists service
	preferencesCol   *mongo.Collection       // Direct access to preferences collection
	interactionsCol  *mongo.Collection       // User interactions
	trendingCacheCol *mongo.Collection       // Trending cache
}

// UserPreferenceAlias for internal use (avoids import cycles)
type UserPreferenceAlias = preferences.UserPreference
