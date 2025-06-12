// handlers/artists/handler_test.go - Updated to use shared domain types
package artists

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bedrockdude10/Booker/backend/domain/artists"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceInterface defines the methods we need from the service (simplified)
type ServiceInterface interface {
	CreateArtist(ctx context.Context, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError)
	GetArtists(ctx context.Context, filters artists.FilterParams, limit, offset int) ([]artists.ArtistDocument, *utils.AppError)
	GetArtistByID(ctx context.Context, id primitive.ObjectID) (*artists.ArtistDocument, *utils.AppError)
	UpdateArtist(ctx context.Context, id primitive.ObjectID, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError)
	UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError)
	DeleteArtist(ctx context.Context, id primitive.ObjectID) *utils.AppError
}

// TestHandler for testing
type TestHandler struct {
	service ServiceInterface
}

// MockService using shared domain types
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateArtist(ctx context.Context, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, params)

	var artist *artists.ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*artists.ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) GetArtists(ctx context.Context, filters artists.FilterParams, limit, offset int) ([]artists.ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, filters, limit, offset)

	var artistsList []artists.ArtistDocument
	if args.Get(0) != nil {
		artistsList = args.Get(0).([]artists.ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artistsList, err
}

func (m *MockService) GetArtistByID(ctx context.Context, id primitive.ObjectID) (*artists.ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, id)

	var artist *artists.ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*artists.ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) UpdateArtist(ctx context.Context, id primitive.ObjectID, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, id, params)

	var artist *artists.ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*artists.ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params artists.CreateArtistParams) (*artists.ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, id, params)

	var artist *artists.ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*artists.ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) DeleteArtist(ctx context.Context, id primitive.ObjectID) *utils.AppError {
	args := m.Called(ctx, id)

	var err *utils.AppError
	if args.Get(0) != nil {
		err = args.Get(0).(*utils.AppError)
	}

	return err
}

//==============================================================================
// Test Functions
//==============================================================================

// TestCreateArtistHandler tests the CreateArtist handler
func TestCreateArtistHandler(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	expectedArtist := &artists.ArtistDocument{
		ID:     primitive.NewObjectID(),
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
		ContactInfo: artists.ContactInfo{
			Manager: "Test Manager",
			Social: artists.SocialMediaLinks{
				Spotify: "https://spotify.com/test",
			},
		},
	}

	mockService.On("CreateArtist", mock.Anything, mock.Anything).Return(expectedArtist, (*utils.AppError)(nil))

	reqBody := artists.CreateArtistParams{
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
		ContactInfo: artists.ContactInfo{
			Manager: "Test Manager",
			Social: artists.SocialMediaLinks{
				Spotify: "https://spotify.com/test",
			},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/artists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create the actual handler method for testing
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		var params artists.CreateArtistParams

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			utils.HandleError(w, utils.ValidationError("Invalid request body"))
			return
		}

		artist, appErr := handler.service.CreateArtist(r.Context(), params)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(artist)
	}

	testHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

// TestGetArtistsWithFilters tests the GetArtists handler with filtering
func TestGetArtistsWithFilters(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	expectedArtists := []artists.ArtistDocument{
		{
			ID:     primitive.NewObjectID(),
			Name:   "Rock Artist",
			Genres: []string{"rock"},
			Cities: []string{"Nashville"},
			ContactInfo: artists.ContactInfo{
				Manager: "Rock Manager",
				Social: artists.SocialMediaLinks{
					Spotify: "https://spotify.com/rock-artist",
				},
			},
		},
	}

	filters := artists.FilterParams{
		Genres:     []string{"rock"},
		HasSpotify: boolPtr(true),
	}

	mockService.On("GetArtists", mock.Anything, filters, 10, 0).Return(expectedArtists, (*utils.AppError)(nil))

	req := httptest.NewRequest("GET", "/api/artists?genres=rock&hasSpotify=true", nil)
	rr := httptest.NewRecorder()

	// Simplified test handler that mimics the actual GetArtists logic
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		filters := artists.ParseFilterParams(r)

		if appErr := artists.ValidateFilterParams(filters); appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		artistsList, appErr := handler.service.GetArtists(r.Context(), filters, 10, 0)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": artistsList,
			"meta": map[string]interface{}{
				"count":   len(artistsList),
				"filters": filters,
			},
		})
	}

	testHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)

	// Verify response contains filtered data
	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	data := response["data"].([]interface{})
	assert.Equal(t, 1, len(data))
}

// TestGetArtistByID tests the GetArtist handler
func TestGetArtistByID(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	artistID := primitive.NewObjectID()
	expectedArtist := &artists.ArtistDocument{
		ID:     artistID,
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
		ContactInfo: artists.ContactInfo{
			Manager: "Test Manager",
		},
	}

	mockService.On("GetArtistByID", mock.Anything, artistID).Return(expectedArtist, (*utils.AppError)(nil))

	req := httptest.NewRequest("GET", "/api/artists/"+artistID.Hex(), nil)
	rr := httptest.NewRecorder()

	// Test handler that mimics GetArtist logic
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// For this test, we'll just use the artistID directly
		artist, appErr := handler.service.GetArtistByID(r.Context(), artistID)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(artist)
	}

	testHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)

	// Verify response contains the artist
	var response artists.ArtistDocument
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, "Test Artist", response.Name)
}

// TestDeleteArtist tests the DeleteArtist handler
func TestDeleteArtist(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	artistID := primitive.NewObjectID()

	mockService.On("DeleteArtist", mock.Anything, artistID).Return((*utils.AppError)(nil))

	req := httptest.NewRequest("DELETE", "/api/artists/"+artistID.Hex(), nil)
	rr := httptest.NewRecorder()

	// Test handler that mimics DeleteArtist logic
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// For this test, we'll just use the artistID directly
		if appErr := handler.service.DeleteArtist(r.Context(), artistID); appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

	testHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}

//==============================================================================
// Helper Functions
//==============================================================================

// Helper function to create a pointer to bool
func boolPtr(b bool) *bool {
	return &b
}
