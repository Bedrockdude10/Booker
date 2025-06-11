// handlers/artists/handler_test.go - Updated for filtering support
package artists

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bedrockdude10/Booker/backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceInterface defines the methods we need from the service (updated for filtering)
type ServiceInterface interface {
	CreateArtist(ctx context.Context, params CreateArtistParams) (*ArtistDocument, *utils.AppError)
	GetArtists(ctx context.Context, filters FilterParams, limit, offset int) ([]ArtistDocument, *utils.AppError)
	GetArtistByID(ctx context.Context, id primitive.ObjectID) (*ArtistDocument, *utils.AppError)
	UpdateArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError)
	UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError)
	DeleteArtist(ctx context.Context, id primitive.ObjectID) *utils.AppError
	GetArtistsByGenre(ctx context.Context, genre string, additionalFilters FilterParams) ([]ArtistDocument, *utils.AppError)
	GetArtistsByCity(ctx context.Context, city string, additionalFilters FilterParams) ([]ArtistDocument, *utils.AppError)
}

// Update Handler to use the interface
type TestHandler struct {
	service ServiceInterface
}

// Mock service (updated for filtering)
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateArtist(ctx context.Context, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, params)

	var artist *ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) GetArtists(ctx context.Context, filters FilterParams, limit, offset int) ([]ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, filters, limit, offset)

	var artists []ArtistDocument
	if args.Get(0) != nil {
		artists = args.Get(0).([]ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artists, err
}

func (m *MockService) GetArtistByID(ctx context.Context, id primitive.ObjectID) (*ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, id)

	var artist *ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) UpdateArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, id, params)

	var artist *ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, id, params)

	var artist *ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*ArtistDocument)
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

func (m *MockService) GetArtistsByGenre(ctx context.Context, genre string, additionalFilters FilterParams) ([]ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, genre, additionalFilters)

	var artists []ArtistDocument
	if args.Get(0) != nil {
		artists = args.Get(0).([]ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artists, err
}

func (m *MockService) GetArtistsByCity(ctx context.Context, city string, additionalFilters FilterParams) ([]ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, city, additionalFilters)

	var artists []ArtistDocument
	if args.Get(0) != nil {
		artists = args.Get(0).([]ArtistDocument)
	}

	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artists, err
}

// Test functions (updated examples)
func TestCreateArtistHandler(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	expectedArtist := &ArtistDocument{
		ID:     primitive.NewObjectID(),
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
	}

	mockService.On("CreateArtist", mock.Anything, mock.Anything).Return(expectedArtist, (*utils.AppError)(nil))

	reqBody := CreateArtistParams{
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/artists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create the actual handler method for testing
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		var params CreateArtistParams

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

// Test filtering functionality
func TestGetArtistsWithFilters(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	expectedArtists := []ArtistDocument{
		{
			ID:     primitive.NewObjectID(),
			Name:   "Rock Artist",
			Genres: []string{"rock"},
			Cities: []string{"Nashville"},
			Rating: 4.5,
		},
	}

	filters := FilterParams{
		Genres:    []string{"rock"},
		MinRating: 4.0,
	}

	mockService.On("GetArtists", mock.Anything, filters, 10, 0).Return(expectedArtists, (*utils.AppError)(nil))

	req := httptest.NewRequest("GET", "/api/artists?genres=rock&minRating=4.0", nil)
	rr := httptest.NewRecorder()

	// Simplified test handler that mimics the actual GetArtists logic
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		filters := ParseFilterParams(r)

		if appErr := ValidateFilterParams(filters); appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		artists, appErr := handler.service.GetArtists(r.Context(), filters, 10, 0)
		if appErr != nil {
			utils.HandleError(w, appErr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": artists,
			"meta": map[string]interface{}{
				"count":   len(artists),
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
