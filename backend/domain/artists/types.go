// domain/artists/types.go
package artists

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ArtistDocument represents the core artist data structure
type ArtistDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Genres      []string           `bson:"genres" json:"genres"`
	Cities      []string           `bson:"cities" json:"cities"`
	ContactInfo ContactInfo        `bson:"contactInfo,omitempty" json:"contactInfo,omitempty"`
}

// CreateArtistParams for creating new artists
type CreateArtistParams struct {
	Name        string      `json:"name" validate:"required,min=1,max=100"`
	Genres      []string    `json:"genres" validate:"required,min=1,validgenres"`
	Cities      []string    `json:"cities" validate:"required,min=1"`
	ContactInfo ContactInfo `json:"contactInfo,omitempty"`
}

// ContactInfo represents all contact and social information for an artist
type ContactInfo struct {
	Social SocialMediaLinks `bson:"social,omitempty" json:"social,omitempty"` // All social media and streaming links

	// Professional Contact Information
	Manager     string `bson:"manager,omitempty" json:"manager,omitempty"`         // Manager name
	ManagerInfo string `bson:"managerInfo,omitempty" json:"managerInfo,omitempty"` // Manager contact details
	BookingInfo string `bson:"bookingInfo,omitempty" json:"bookingInfo,omitempty"` // Booking agent info
	LabelName   string `bson:"labelName,omitempty" json:"labelName,omitempty"`     // Record label
	LabelURL    string `bson:"labelURL,omitempty" json:"labelURL,omitempty"`       // Label website
}

// SocialMediaLinks represents all social media and streaming platform links
type SocialMediaLinks struct {
	// Streaming Services
	Spotify    string `bson:"spotify,omitempty" json:"spotify,omitempty"`       // Spotify artist URL
	AppleMusic string `bson:"appleMusic,omitempty" json:"appleMusic,omitempty"` // Apple Music artist URL
	Bandcamp   string `bson:"bandcamp,omitempty" json:"bandcamp,omitempty"`     // Bandcamp artist page URL

	// Social Media Platforms
	Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"` // Instagram profile URL
	YouTube   string `bson:"youtube,omitempty" json:"youtube,omitempty"`     // YouTube channel URL
	Facebook  string `bson:"facebook,omitempty" json:"facebook,omitempty"`   // Facebook page URL
	Twitter   string `bson:"twitter,omitempty" json:"twitter,omitempty"`     // Twitter/X profile URL
	TikTok    string `bson:"tiktok,omitempty" json:"tiktok,omitempty"`       // TikTok profile URL

	// Professional Platforms
	Website    string `bson:"website,omitempty" json:"website,omitempty"`       // Official website
	SoundCloud string `bson:"soundcloud,omitempty" json:"soundcloud,omitempty"` // SoundCloud profile
	Discogs    string `bson:"discogs,omitempty" json:"discogs,omitempty"`       // Discogs artist page
	Beatport   string `bson:"beatport,omitempty" json:"beatport,omitempty"`     // Beatport artist page
	Deezer     string `bson:"deezer,omitempty" json:"deezer,omitempty"`         // Deezer artist page
	Pandora    string `bson:"pandora,omitempty" json:"pandora,omitempty"`       // Pandora artist page

	// Messaging/Contact
	Email string `bson:"email,omitempty" json:"email,omitempty"` // Contact email
	Phone string `bson:"phone,omitempty" json:"phone,omitempty"` // Contact phone
}
