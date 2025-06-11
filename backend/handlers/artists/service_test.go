// // handlers/artists/service_test.go
package artists

// import (
// 	"context"
// 	"encoding/json"
// 	"os"
// 	"testing"

// 	"github.com/Bedrockdude10/Booker/backend/domain"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
// )

// // ServiceTestSuite provides a clean test environment with setup/teardown
// type ServiceTestSuite struct {
// 	suite.Suite
// 	mt      *mtest.T
// 	service *Service
// 	ctx     context.Context
// }

// // SetupSuite runs once before all tests
// func (suite *ServiceTestSuite) SetupSuite() {
// 	suite.ctx = context.Background()
// }

// // SetupTest runs before each individual test
// func (suite *ServiceTestSuite) SetupTest() {
// 	// Create fresh mock database for each test
// 	mt := mtest.New(suite.T(), mtest.NewOptions().ClientType(mtest.Mock))
// 	suite.mt = mt

// 	collections := map[string]*mongo.Collection{
// 		"artists": mt.Coll,
// 	}
// 	suite.service = NewService(collections)
// }

// // TearDownTest runs after each individual test
// func (suite *ServiceTestSuite) TearDownTest() {
// 	// mtest.T doesn't have a Close method - it cleans up automatically
// 	// Just set to nil for clarity
// 	suite.mt = nil
// }

// // Helper method to load sample artists
// func (suite *ServiceTestSuite) loadSampleArtists() []ArtistDocument {
// 	data, err := os.ReadFile("../../testdata/artists.json")
// 	if err != nil {
// 		suite.T().Logf("No sample artists file found, using test data: %v", err)
// 		return suite.createTestArtists()
// 	}

// 	var artists []ArtistDocument
// 	if err := json.Unmarshal(data, &artists); err != nil {
// 		suite.T().Logf("Failed to parse JSON, using test data: %v", err)
// 		return suite.createTestArtists()
// 	}

// 	// Generate ObjectIDs if missing
// 	for i := range artists {
// 		if artists[i].ID.IsZero() {
// 			artists[i].ID = primitive.NewObjectID()
// 		}
// 	}

// 	return artists
// }

// // Helper method to create test artists
// func (suite *ServiceTestSuite) createTestArtists() []ArtistDocument {
// 	return []ArtistDocument{
// 		{
// 			ID:        primitive.NewObjectID(),
// 			Name:      "Test Rock Band",
// 			Genres:    []string{"rock", "hard-rock"},
// 			Cities:    []string{"Nashville", "Memphis"},
// 			Manager:   "Rock Manager",
// 			SpotifyID: "test-rock-band-123",
// 		},
// 		{
// 			ID:        primitive.NewObjectID(),
// 			Name:      "Test Pop Artist",
// 			Genres:    []string{"pop", "dance"},
// 			Cities:    []string{"Los Angeles", "Las Vegas"},
// 			Manager:   "Pop Manager",
// 			SpotifyID: "test-pop-artist-456",
// 		},
// 		{
// 			ID:        primitive.NewObjectID(),
// 			Name:      "Nashville Country Star",
// 			Genres:    []string{"country", "folk"},
// 			Cities:    []string{"Nashville", "Austin"},
// 			Manager:   "Country Manager",
// 			SpotifyID: "country-star-789",
// 		},
// 	}
// }

// // Helper method to create BSON document from artist
// func (suite *ServiceTestSuite) artistToBSON(artist ArtistDocument) bson.D {
// 	return bson.D{
// 		{Key: "_id", Value: artist.ID},
// 		{Key: "name", Value: artist.Name},
// 		{Key: "genres", Value: artist.Genres},
// 		{Key: "cities", Value: artist.Cities},
// 		{Key: "manager", Value: artist.Manager},
// 		{Key: "spotifyId", Value: artist.SpotifyID},
// 	}
// }

// // Helper method to mock cursor response for multiple artists
// func (suite *ServiceTestSuite) mockMultipleArtists(artists []ArtistDocument) {
// 	if len(artists) == 0 {
// 		suite.mt.AddMockResponses(
// 			mtest.CreateCursorResponse(0, "booker_test.artists", mtest.FirstBatch),
// 		)
// 		return
// 	}

// 	// Create first batch with first artist
// 	first := mtest.CreateCursorResponse(1, "booker_test.artists", mtest.FirstBatch, suite.artistToBSON(artists[0]))

// 	// Add remaining artists as next batches if needed
// 	var responses []bson.D
// 	responses = append(responses, first)

// 	// End cursor
// 	getNext := mtest.CreateCursorResponse(0, "booker_test.artists", mtest.NextBatch)
// 	killCursors := mtest.CreateCursorResponse(0, "booker_test.artists", mtest.NextBatch)

// 	suite.mt.AddMockResponses(first, getNext, killCursors)
// }

// // Test CreateArtist
// func (suite *ServiceTestSuite) TestCreateArtist() {
// 	suite.mt.AddMockResponses(mtest.CreateSuccessResponse())

// 	params := CreateArtistParams{
// 		Name:      "New Test Artist",
// 		Genres:    []string{"rock", "pop"},
// 		Cities:    []string{"Nashville"},
// 		Manager:   "Test Manager",
// 		SpotifyID: "new-test-123",
// 	}

// 	artist, err := suite.service.CreateArtist(suite.ctx, params)

// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), artist)
// 	assert.Equal(suite.T(), "New Test Artist", artist.Name)
// 	assert.Contains(suite.T(), artist.Genres, "rock")
// 	assert.Contains(suite.T(), artist.Genres, "pop")
// 	assert.False(suite.T(), artist.ID.IsZero())
// }

// func (suite *ServiceTestSuite) TestCreateArtistInvalidGenre() {
// 	// No mock response needed - should fail validation before DB call
// 	params := CreateArtistParams{
// 		Name:   "Test Artist",
// 		Genres: []string{"invalid-genre"},
// 		Cities: []string{"Nashville"},
// 	}

// 	artist, err := suite.service.CreateArtist(suite.ctx, params)

// 	assert.Error(suite.T(), err)
// 	assert.Nil(suite.T(), artist)
// }

// // Test GetArtists
// func (suite *ServiceTestSuite) TestGetArtists() {
// 	sampleArtists := suite.loadSampleArtists()
// 	suite.mockMultipleArtists(sampleArtists)

// 	artists, err := suite.service.GetArtists(suite.ctx)

// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), artists)

// 	if len(artists) > 0 {
// 		assert.NotEmpty(suite.T(), artists[0].Name)
// 		assert.NotEmpty(suite.T(), artists[0].Genres)
// 		assert.NotEmpty(suite.T(), artists[0].Cities)
// 	}
// }

// // Test GetArtistsByGenre
// func (suite *ServiceTestSuite) TestGetArtistsByGenreValid() {
// 	rockArtist := ArtistDocument{
// 		ID:     primitive.NewObjectID(),
// 		Name:   "Rock Band",
// 		Genres: []string{"rock", "hard-rock"},
// 		Cities: []string{"Nashville"},
// 	}

// 	suite.mockMultipleArtists([]ArtistDocument{rockArtist})

// 	artists, err := suite.service.GetAllArtistsByGenre(suite.ctx, "rock")

// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), artists)
// }

// func (suite *ServiceTestSuite) TestGetArtistsByGenreInvalid() {
// 	// No mock response needed - should fail validation
// 	artists, err := suite.service.GetAllArtistsByGenre(suite.ctx, "invalid-genre")

// 	assert.Error(suite.T(), err)
// 	assert.Nil(suite.T(), artists)
// 	assert.Contains(suite.T(), err.Error(), "Invalid genre")
// }

// // Test GetArtistsByCity
// func (suite *ServiceTestSuite) TestGetArtistsByCity() {
// 	nashvilleArtist := ArtistDocument{
// 		ID:     primitive.NewObjectID(),
// 		Name:   "Nashville Artist",
// 		Genres: []string{"country"},
// 		Cities: []string{"Nashville", "Memphis"},
// 	}

// 	suite.mockMultipleArtists([]ArtistDocument{nashvilleArtist})

// 	artists, err := suite.service.GetArtistsByCity(suite.ctx, "Nashville")

// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), artists)
// }

// func (suite *ServiceTestSuite) TestGetArtistsByCityEmpty() {
// 	artists, err := suite.service.GetArtistsByCity(suite.ctx, "")

// 	assert.Error(suite.T(), err)
// 	assert.Nil(suite.T(), artists)
// 	assert.Contains(suite.T(), err.Error(), "City is required")
// }

// // Test GetArtistByID
// func (suite *ServiceTestSuite) TestGetArtistByIDFound() {
// 	artistID := primitive.NewObjectID()
// 	artist := ArtistDocument{
// 		ID:     artistID,
// 		Name:   "Found Artist",
// 		Genres: []string{"rock"},
// 		Cities: []string{"Nashville"},
// 	}

// 	suite.mt.AddMockResponses(
// 		mtest.CreateCursorResponse(1, "booker_test.artists", mtest.FirstBatch, suite.artistToBSON(artist)),
// 	)

// 	foundArtist, err := suite.service.GetArtistByID(suite.ctx, artistID)

// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), foundArtist)
// 	assert.Equal(suite.T(), artistID, foundArtist.ID)
// 	assert.Equal(suite.T(), "Found Artist", foundArtist.Name)
// }

// func (suite *ServiceTestSuite) TestGetArtistByIDNotFound() {
// 	// Mock empty response (not found)
// 	suite.mt.AddMockResponses(
// 		mtest.CreateCursorResponse(0, "booker_test.artists", mtest.FirstBatch),
// 	)

// 	nonExistentID := primitive.NewObjectID()
// 	artist, err := suite.service.GetArtistByID(suite.ctx, nonExistentID)

// 	assert.Error(suite.T(), err)
// 	assert.Nil(suite.T(), artist)
// 	assert.Contains(suite.T(), err.Error(), "Artist not found")
// }

// // Run the test suite
// func TestServiceSuite(t *testing.T) {
// 	suite.Run(t, new(ServiceTestSuite))
// }

// // Additional unit tests that don't need mocking
// func TestDomainValidation(t *testing.T) {
// 	t.Run("valid genres", func(t *testing.T) {
// 		validGenres := []string{"rock", "pop", "jazz", "country", "electronic"}
// 		for _, genre := range validGenres {
// 			assert.True(t, domain.HasGenre(genre), "Genre %s should be valid", genre)
// 		}
// 	})

// 	t.Run("invalid genres", func(t *testing.T) {
// 		invalidGenres := []string{"invalid-genre", "", "not-a-real-genre"}
// 		for _, genre := range invalidGenres {
// 			assert.False(t, domain.HasGenre(genre), "Genre %s should be invalid", genre)
// 		}
// 	})

// 	t.Run("get all genres", func(t *testing.T) {
// 		genres := domain.GetAllGenres()
// 		assert.Greater(t, len(genres), 0)
// 		assert.Contains(t, genres, "rock")
// 		assert.Contains(t, genres, "pop")
// 	})
// }
