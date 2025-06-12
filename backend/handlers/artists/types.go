// handlers/artists/types.go
package artists

import (
	"github.com/Bedrockdude10/Booker/backend/domain/artists"
	"go.mongodb.org/mongo-driver/mongo"
)

// Type aliases to shared domain types (for convenience)
type ArtistDocument = artists.ArtistDocument
type CreateArtistParams = artists.CreateArtistParams
type ContactInfo = artists.ContactInfo
type SocialMediaLinks = artists.SocialMediaLinks
type FilterParams = artists.FilterParams

// Service struct for artists operations
type Service struct {
	artists         *mongo.Collection
	userPreferences *mongo.Collection
}
