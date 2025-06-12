// handlers/discovery/bandcamp.go
package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Bedrockdude10/Booker/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BandcampAPIResponse represents the response from Bandcamp's discover API
type BandcampAPIResponse struct {
	Results          []BandcampResult `json:"results"`
	ResultCount      int              `json:"result_count"`
	BatchResultCount int              `json:"batch_result_count"`
	Cursor           string           `json:"cursor"`
}

// BandcampResult represents a single artist/album from Bandcamp - FULL DATA
type BandcampResult struct {
	ID                     int64             `json:"id"`
	Title                  string            `json:"title"`
	IsAlbumPreorder        bool              `json:"is_album_preorder"`
	IsFreeDownload         bool              `json:"is_free_download"`
	IsPurchasable          bool              `json:"is_purchasable"`
	IsSetPrice             bool              `json:"is_set_price"`
	ItemURL                string            `json:"item_url"`
	ItemPrice              float64           `json:"item_price"`
	Price                  BandcampPrice     `json:"price"`
	ItemCurrency           string            `json:"item_currency"`
	ItemImageID            int64             `json:"item_image_id"`
	ResultType             string            `json:"result_type"` // "a" for album, "s" for single/merch
	BandID                 int64             `json:"band_id"`
	AlbumArtist            *string           `json:"album_artist"`
	BandName               string            `json:"band_name"`
	BandURL                string            `json:"band_url"`
	BandBioImageID         int64             `json:"band_bio_image_id"`
	BandLatestArtID        int64             `json:"band_latest_art_id"`
	BandGenreID            int               `json:"band_genre_id"`
	ReleaseDate            string            `json:"release_date"`
	TotalPackageCount      int               `json:"total_package_count"`
	PackageInfo            []BandcampPackage `json:"package_info,omitempty"`
	FeaturedTrack          *BandcampTrack    `json:"featured_track,omitempty"`
	BandLocation           string            `json:"band_location"`
	TrackCount             int               `json:"track_count"`
	ItemDuration           float64           `json:"item_duration"`
	ItemTags               interface{}       `json:"item_tags"` // Can be null or array
	IsFollowing            bool              `json:"is_following"`
	IsWishlisted           bool              `json:"is_wishlisted"`
	IsOwned                bool              `json:"is_owned"`
	LabelName              *string           `json:"label_name,omitempty"`
	LabelURL               *string           `json:"label_url,omitempty"`
	TshirtSecondaryImageID *int64            `json:"tshirt_secondary_image_id,omitempty"`
}

type BandcampPrice struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	IsMoney  bool   `json:"is_money"`
}

type BandcampPackage struct {
	ID                int64         `json:"id"`
	Title             string        `json:"title"`
	Format            string        `json:"format"`
	IsDigitalIncluded bool          `json:"is_digital_included"`
	IsSetPrice        bool          `json:"is_set_price"`
	IsPreorder        bool          `json:"is_preorder"`
	ImageID           int64         `json:"image_id"`
	Price             BandcampPrice `json:"price"`
	TypeID            int           `json:"type_id"`
}

type BandcampTrack struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	BandName  string  `json:"band_name"`
	BandID    int64   `json:"band_id"`
	Duration  float64 `json:"duration"`
	StreamURL string  `json:"stream_url"`
}

// ScrapedArtist represents a unique artist with their most recent release data
type ScrapedArtist struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	// Core Artist Info
	BandcampID      int64  `bson:"bandcamp_id" json:"bandcamp_id"`             // Unique artist identifier
	Name            string `bson:"name" json:"name"`                           // Artist name
	BandcampURL     string `bson:"bandcamp_url" json:"bandcamp_url"`           // Artist's Bandcamp page
	Location        string `bson:"location" json:"location"`                   // Artist location
	BandcampGenreID int    `bson:"bandcamp_genre_id" json:"bandcamp_genre_id"` // Primary genre

	// Artist Images & Branding
	BandBioImageID  int64 `bson:"band_bio_image_id" json:"band_bio_image_id"`   // Profile image
	BandLatestArtID int64 `bson:"band_latest_art_id" json:"band_latest_art_id"` // Latest artwork

	// Label & Professional Status
	LabelName *string `bson:"label_name,omitempty" json:"label_name,omitempty"` // Record label
	LabelURL  *string `bson:"label_url,omitempty" json:"label_url,omitempty"`   // Label's page

	// Most Recent Release Info (activity indicator)
	LatestReleaseID    int64     `bson:"latest_release_id" json:"latest_release_id"`
	LatestReleaseTitle string    `bson:"latest_release_title" json:"latest_release_title"`
	LatestReleaseURL   string    `bson:"latest_release_url" json:"latest_release_url"`
	LatestReleaseDate  time.Time `bson:"latest_release_date" json:"latest_release_date"`
	LatestReleaseType  string    `bson:"latest_release_type" json:"latest_release_type"` // "a" for album, "s" for single

	// Activity & Popularity Indicators
	IsFreeDownload        bool    `bson:"is_free_download" json:"is_free_download"`         // Accessibility
	LatestReleasePrice    float64 `bson:"latest_release_price" json:"latest_release_price"` // Pricing strategy
	LatestReleaseCurrency string  `bson:"latest_release_currency" json:"latest_release_currency"`
	TrackCount            int     `bson:"track_count" json:"track_count"`                 // Prolific indicator
	ItemDuration          float64 `bson:"item_duration" json:"item_duration"`             // Total music length
	TotalPackageCount     int     `bson:"total_package_count" json:"total_package_count"` // Physical releases

	// Sample Track (for preview/analysis)
	FeaturedTrack *BandcampTrack `bson:"featured_track,omitempty" json:"featured_track,omitempty"`

	// Additional Genre/Style Info
	ItemTags interface{} `bson:"item_tags,omitempty" json:"item_tags,omitempty"` // Additional tags

	// Physical Releases (indie cred indicator)
	HasPhysicalReleases bool              `bson:"has_physical_releases" json:"has_physical_releases"`
	PackageInfo         []BandcampPackage `bson:"package_info,omitempty" json:"package_info,omitempty"`

	// Metadata
	ScrapedAt time.Time `bson:"scraped_at" json:"scraped_at"`
	Source    string    `bson:"source" json:"source"` // "bandcamp"

	// Future Spotify Integration
	SpotifyID          string     `bson:"spotify_id,omitempty" json:"spotify_id,omitempty"`
	MonthlyListeners   int        `bson:"monthly_listeners,omitempty" json:"monthly_listeners,omitempty"`
	SpotifyProcessed   bool       `bson:"spotify_processed" json:"spotify_processed"`
	SpotifyProcessedAt *time.Time `bson:"spotify_processed_at,omitempty" json:"spotify_processed_at,omitempty"`
}

// BandcampService handles scraping and storing Bandcamp data
type BandcampService struct {
	client            *http.Client
	scrapedCollection *mongo.Collection
}

// NewBandcampService creates a new service
func NewBandcampService(scrapedCollection *mongo.Collection) *BandcampService {
	return &BandcampService{
		client:            &http.Client{Timeout: 30 * time.Second},
		scrapedCollection: scrapedCollection,
	}
}

// ScrapeBostonArtists fetches Boston artists from Bandcamp and stores them
func (bs *BandcampService) ScrapeBostonArtists(ctx context.Context, limit int) *utils.AppError {
	slog.InfoContext(ctx, "Starting Bandcamp scraping", "limit", limit)

	// Boston geoname_id from your working example
	geonameID := 4930956

	response, appErr := bs.fetchFromBandcamp(ctx, geonameID, limit)
	if appErr != nil {
		return appErr
	}

	slog.InfoContext(ctx, "Fetched results from Bandcamp",
		"results", len(response.Results),
		"total_available", response.ResultCount)

	// Process and store results
	artists := bs.processBandcampResults(response.Results)

	return bs.storeArtists(ctx, artists)
}

// fetchFromBandcamp calls the Bandcamp API
func (bs *BandcampService) fetchFromBandcamp(ctx context.Context, geonameID, limit int) (*BandcampAPIResponse, *utils.AppError) {
	url := "https://bandcamp.com/api/discover/1/discover_web"

	// JSON payload
	payload := fmt.Sprintf(`{
		"category_id": 0,
		"tag_norm_names": [],
		"geoname_id": %d,
		"slice": "new",
		"time_facet_id": null,
		"cursor": "*",
		"size": %d,
		"include_result_types": ["a", "s"]
	}`, geonameID, limit)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, utils.InternalErrorLog(ctx, "Failed to create Bandcamp request", err)
	}

	// Headers from your working example
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://bandcamp.com/discover")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := bs.client.Do(req)
	if err != nil {
		return nil, utils.InternalErrorLog(ctx, "Bandcamp request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, utils.InternalError(fmt.Sprintf("Bandcamp API returned status %d", resp.StatusCode), nil)
	}

	var apiResponse BandcampAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, utils.InternalErrorLog(ctx, "Failed to parse Bandcamp response", err)
	}

	return &apiResponse, nil
}

// processBandcampResults converts API results to unique artists with latest release data
func (bs *BandcampService) processBandcampResults(results []BandcampResult) []ScrapedArtist {
	var artists []ScrapedArtist
	seen := make(map[int64]*BandcampResult) // Track unique band IDs with most recent release

	for _, result := range results {
		// Skip merch items - focus on music
		if result.ResultType == "s" {
			continue
		}

		// If we've seen this artist, keep the most recent release
		if existing, exists := seen[result.BandID]; exists {
			existingDate, _ := time.Parse("2006-01-02 15:04:05 UTC", existing.ReleaseDate)
			currentDate, _ := time.Parse("2006-01-02 15:04:05 UTC", result.ReleaseDate)

			// Keep the more recent release
			if currentDate.After(existingDate) {
				seen[result.BandID] = &result
			}
		} else {
			seen[result.BandID] = &result
		}
	}

	// Convert to ScrapedArtist format
	for _, result := range seen {
		releaseDate, _ := time.Parse("2006-01-02 15:04:05 UTC", result.ReleaseDate)

		artist := ScrapedArtist{
			ID:              primitive.NewObjectID(),
			BandcampID:      result.BandID,
			Name:            result.BandName,
			BandcampURL:     result.BandURL,
			Location:        result.BandLocation,
			BandcampGenreID: result.BandGenreID,
			BandBioImageID:  result.BandBioImageID,
			BandLatestArtID: result.BandLatestArtID,

			// Label info
			LabelName: result.LabelName,
			LabelURL:  result.LabelURL,

			// Latest release info
			LatestReleaseID:       result.ID,
			LatestReleaseTitle:    result.Title,
			LatestReleaseURL:      result.ItemURL,
			LatestReleaseDate:     releaseDate,
			LatestReleaseType:     result.ResultType,
			IsFreeDownload:        result.IsFreeDownload,
			LatestReleasePrice:    result.ItemPrice,
			LatestReleaseCurrency: result.ItemCurrency,
			TrackCount:            result.TrackCount,
			ItemDuration:          result.ItemDuration,
			TotalPackageCount:     result.TotalPackageCount,

			// Additional data
			FeaturedTrack:       result.FeaturedTrack,
			ItemTags:            result.ItemTags,
			HasPhysicalReleases: len(result.PackageInfo) > 0,
			PackageInfo:         result.PackageInfo,

			// Metadata
			ScrapedAt:        time.Now(),
			Source:           "bandcamp",
			SpotifyProcessed: false,
		}

		artists = append(artists, artist)
	}

	return artists
}

// storeArtists saves artists to MongoDB with upsert logic
func (bs *BandcampService) storeArtists(ctx context.Context, artists []ScrapedArtist) *utils.AppError {
	if len(artists) == 0 {
		return nil
	}

	// Use bulk operations for efficiency
	var operations []mongo.WriteModel

	for _, artist := range artists {
		// Upsert based on bandcamp_id to avoid duplicates
		filter := bson.M{"bandcamp_id": artist.BandcampID}
		update := bson.M{
			"$set": bson.M{
				"name":                    artist.Name,
				"bandcamp_url":            artist.BandcampURL,
				"location":                artist.Location,
				"bandcamp_genre_id":       artist.BandcampGenreID,
				"band_bio_image_id":       artist.BandBioImageID,
				"band_latest_art_id":      artist.BandLatestArtID,
				"label_name":              artist.LabelName,
				"label_url":               artist.LabelURL,
				"latest_release_id":       artist.LatestReleaseID,
				"latest_release_title":    artist.LatestReleaseTitle,
				"latest_release_url":      artist.LatestReleaseURL,
				"latest_release_date":     artist.LatestReleaseDate,
				"latest_release_type":     artist.LatestReleaseType,
				"is_free_download":        artist.IsFreeDownload,
				"latest_release_price":    artist.LatestReleasePrice,
				"latest_release_currency": artist.LatestReleaseCurrency,
				"track_count":             artist.TrackCount,
				"item_duration":           artist.ItemDuration,
				"total_package_count":     artist.TotalPackageCount,
				"featured_track":          artist.FeaturedTrack,
				"item_tags":               artist.ItemTags,
				"has_physical_releases":   artist.HasPhysicalReleases,
				"package_info":            artist.PackageInfo,
				"source":                  artist.Source,
			},
			"$setOnInsert": bson.M{
				"_id":               artist.ID,
				"bandcamp_id":       artist.BandcampID,
				"scraped_at":        artist.ScrapedAt,
				"spotify_processed": false,
			},
		}

		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(filter)
		operation.SetUpdate(update)
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	// Execute bulk write
	opts := options.BulkWrite().SetOrdered(false)
	result, err := bs.scrapedCollection.BulkWrite(ctx, operations, opts)
	if err != nil {
		return utils.DatabaseErrorLog(ctx, "bulk upsert scraped artists", err)
	}

	slog.InfoContext(ctx, "Stored artists successfully",
		"inserted", result.InsertedCount,
		"modified", result.ModifiedCount,
		"upserted", result.UpsertedCount)

	return nil
}

// GetScrapedArtists retrieves artists from the collection
func (bs *BandcampService) GetScrapedArtists(ctx context.Context, limit int) ([]ScrapedArtist, *utils.AppError) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	opts.SetSort(bson.M{"scraped_at": -1}) // Most recent first

	cursor, err := bs.scrapedCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find scraped artists", err)
	}
	defer cursor.Close(ctx)

	var artists []ScrapedArtist
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode scraped artists", err)
	}

	return artists, nil
}

// GetArtistCount returns the total number of scraped artists
func (bs *BandcampService) GetArtistCount(ctx context.Context) (int64, *utils.AppError) {
	count, err := bs.scrapedCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, utils.DatabaseErrorLog(ctx, "count scraped artists", err)
	}
	return count, nil
}
