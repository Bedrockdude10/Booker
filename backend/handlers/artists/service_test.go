// handlers/artists/service_test.go
package artists

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceTestSuite struct {
	suite.Suite
	service   *Service
	container *mongodb.MongoDBContainer
	client    *mongo.Client
}

func (suite *ServiceTestSuite) SetupSuite() {
	ctx := context.Background()

	// Start MongoDB container
	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:6"),
	)
	suite.Require().NoError(err)
	suite.container = mongoContainer

	// Connect to MongoDB
	uri, err := mongoContainer.ConnectionString(ctx)
	suite.Require().NoError(err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	suite.Require().NoError(err)
	suite.client = client

	// Set up collections
	db := client.Database("booker_test")
	collections := map[string]*mongo.Collection{
		"artists": db.Collection("artists"),
	}

	// Initialize service
	suite.service = NewService(collections)

	// Load sample data
	suite.loadSampleArtists(ctx, collections["artists"])
}

func (suite *ServiceTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.client != nil {
		suite.client.Disconnect(ctx)
	}
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ServiceTestSuite) loadSampleArtists(ctx context.Context, collection *mongo.Collection) {
	// Read the JSON file
	data, err := os.ReadFile("../../testdata/artists.json")
	if err != nil {
		suite.T().Logf("No sample artists file found, skipping: %v", err)
		return
	}

	// Parse JSON
	var artists []ArtistDocument
	err = json.Unmarshal(data, &artists)
	suite.Require().NoError(err)

	// Generate ObjectIDs if missing
	for i := range artists {
		if artists[i].ID.IsZero() {
			artists[i].ID = primitive.NewObjectID()
		}
	}

	// Insert into database
	if len(artists) > 0 {
		docs := make([]interface{}, len(artists))
		for i, artist := range artists {
			docs[i] = artist
		}
		_, err = collection.InsertMany(ctx, docs)
		suite.Require().NoError(err)
	}
}

func (suite *ServiceTestSuite) TestCreateArtist() {
	ctx := context.Background()

	params := CreateArtistParams{
		Name:   "Test Artist",
		Genres: []string{"rock", "pop"},
		Cities: []string{"Nashville"},
	}

	artist, err := suite.service.CreateArtist(ctx, params)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Artist", artist.Name)
	assert.Contains(suite.T(), artist.Genres, "rock")
}

func (suite *ServiceTestSuite) TestGetArtists() {
	ctx := context.Background()

	artists, err := suite.service.GetArtists(ctx)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), artists)
}

func (suite *ServiceTestSuite) TestGetArtistsByGenre() {
	ctx := context.Background()

	artists, err := suite.service.GetAllArtistsByGenre(ctx, "rock")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), artists)

	for _, artist := range artists {
		assert.Contains(suite.T(), artist.Genres, "rock")
	}
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
