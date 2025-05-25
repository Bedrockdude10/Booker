// preferences/types.go
package preferences

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserPreference struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          primitive.ObjectID `bson:"userId"` // References Account.ID
	PreferredGenres []string           `bson:"preferredGenres"`
	PreferredCities []string           `bson:"preferredCities"`
	// ... rest of fields
}
