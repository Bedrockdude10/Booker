package artists

import (
	"context"
	"strings"

	"github.com/Bedrockdude10/Booker/backend/utils"
)

// All genre-related constants and data
var ValidGenres = utils.NewSet(
	"acoustic",
	"afrobeat",
	"alt-rock",
	"alternative",
	"ambient",
	"blues",
	"bossanova",
	"brazil",
	"breakbeat",
	"british",
	"chill",
	"classical",
	"club",
	"country",
	"dance",
	"dancehall",
	"death-metal",
	"deep-house",
	"detroit-techno",
	"disco",
	"drum-and-bass",
	"dub",
	"dubstep",
	"edm",
	"electronic",
	"emo",
	"folk",
	"forro",
	"french",
	"funk",
	"garage",
	"gospel",
	"goth",
	"grunge",
	"hard-rock",
	"hardcore",
	"hardstyle",
	"heavy-metal",
	"hip-hop",
	"house",
	"idm",
	"indian",
	"indie",
	"indie-pop",
	"industrial",
	"iranian",
	"j-dance",
	"j-idol",
	"j-pop",
	"j-rock",
	"jazz",
	"k-pop",
	"kids",
	"latin",
	"latino",
	"malay",
	"mandopop",
	"metal",
	"metal-misc",
	"metalcore",
	"minimal-techno",
	"movies",
	"mpb",
	"new-age",
	"new-release",
	"opera",
	"pagode",
	"party",
	"philippines-opm",
	"piano",
	"pop",
	"pop-film",
	"post-dubstep",
	"power-pop",
	"progressive-house",
	"psych-rock",
	"punk",
	"punk-rock",
	"r-n-b",
	"rainy-day",
	"reggae",
	"reggaeton",
	"road-trip",
	"rock",
	"rock-n-roll",
	"rockabilly",
	"romance",
	"sad",
	"salsa",
	"samba",
	"sertanejo",
	"show-tunes",
	"singer-songwriter",
	"ska",
	"sleep",
	"songwriter",
	"soul",
	"soundtracks",
	"spanish",
	"study",
	"summer",
	"swedish",
	"synth-pop",
	"tango",
	"techno",
	"trance",
	"trip-hop",
	"turkish",
	"work-out",
	"world-music",
)

// ValidateGenres validates genres using the improved error handling
func ValidateGenres(ctx context.Context, genres []string) *utils.AppError {
	if len(genres) == 0 {
		return utils.ValidationErrorLog(ctx, "At least one genre is required")
	}

	var invalid []string
	for _, genre := range genres {
		if !ValidGenres.Has(genre) {
			invalid = append(invalid, genre)
		}
	}

	if len(invalid) > 0 {
		// Create helpful details with first few valid genres as examples
		allGenres := GetAllGenres()
		exampleCount := 5
		if len(allGenres) < exampleCount {
			exampleCount = len(allGenres)
		}

		details := "Invalid genres: " + strings.Join(invalid, ", ") +
			". Examples of valid genres: " + strings.Join(allGenres[:exampleCount], ", ") + "..."

		// Use Log function for custom attributes
		return utils.Log(ctx,
			utils.ValidationError("Invalid genres provided", details),
			"Genre validation failed",
			"invalid_genres", invalid,
			"invalid_count", len(invalid),
			"total_valid_genres", len(allGenres),
		)
	}

	return nil
}

// ValidateGenresSimple - for backward compatibility, returns standard error
func ValidateGenresSimple(genres []string) error {
	ctx := context.Background()
	if appErr := ValidateGenres(ctx, genres); appErr != nil {
		return appErr
	}
	return nil
}

// GetAllGenres returns all valid genre IDs
func GetAllGenres() []string {
	return ValidGenres.ToSlice()
}

// GetGenreCount returns total count of valid genres
func GetGenreCount() int {
	return ValidGenres.Size()
}

// HasGenre checks if a single genre is valid
func HasGenre(genre string) bool {
	return ValidGenres.Has(genre)
}
