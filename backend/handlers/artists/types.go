package artists

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateArtistParams struct {
	Name      string   `json:"name"`
	Genres    []string `json:"genres"`
	Manager   string   `json:"manager,omitempty"`
	Cities    []string `json:"cities"`
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
}
