// handlers/artists/types.go
package artists

import (
	"github.com/dgraph-io/ristretto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateArtistParams struct {
	Name      string   `json:"name" validate:"required,min=1,max=100"`
	Genres    []string `json:"genres" validate:"required,min=1"`
	Manager   string   `json:"manager,omitempty" validate:"omitempty,min=1,max=100"`
	Cities    []string `json:"cities" validate:"required,min=1"`
	SpotifyID string   `json:"spotifyId,omitempty"`
}

type ArtistDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Genres    []string           `bson:"genres" json:"genres"`
	Manager   string             `bson:"manager,omitempty" json:"manager,omitempty"`
	Cities    []string           `bson:"cities" json:"cities"`
	SpotifyID string             `bson:"spotifyId,omitempty" json:"spotifyId,omitempty"`
}

// type UserPreference struct {
// 	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
// 	UserID   primitive.ObjectID   `bson:"userId" json:"userId"`
// 	Genres   []string             `bson:"genres" json:"genres"`
// 	Artists  []primitive.ObjectID `bson:"artists" json:"artists"`
// 	Location string               `bson:"location" json:"location"`
// }

/*
Artist Service to be used by Artist Handler to interact with the
Database layer of the application
*/
type Service struct {
	artists         *mongo.Collection
	userPreferences *mongo.Collection
	cache           *ristretto.Cache
}
