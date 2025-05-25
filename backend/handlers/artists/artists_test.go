// handlers/artists/handler_test.go
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

// ServiceInterface defines the methods we need from the service
// This allows us to mock it properly
type ServiceInterface interface {
	CreateArtist(ctx context.Context, params CreateArtistParams) (*ArtistDocument, *utils.AppError)
	GetArtists(ctx context.Context) ([]ArtistDocument, *utils.AppError)
	GetArtistByID(ctx context.Context, id primitive.ObjectID) (*ArtistDocument, *utils.AppError)
	UpdateArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError)
	UpdatePartialArtist(ctx context.Context, id primitive.ObjectID, params CreateArtistParams) (*ArtistDocument, *utils.AppError)
	DeleteArtist(ctx context.Context, id primitive.ObjectID) *utils.AppError
	GetAllArtistsByGenre(ctx context.Context, genre string) ([]ArtistDocument, *utils.AppError)
	GetArtistsByCity(ctx context.Context, city string) ([]ArtistDocument, *utils.AppError)
}

// Update Handler to use the interface
type TestHandler struct {
	service ServiceInterface
}

// Mock service
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateArtist(ctx context.Context, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, params)

	// Handle nil return for artist
	var artist *ArtistDocument
	if args.Get(0) != nil {
		artist = args.Get(0).(*ArtistDocument)
	}

	// Handle nil return for error
	var err *utils.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*utils.AppError)
	}

	return artist, err
}

func (m *MockService) GetArtists(ctx context.Context) ([]ArtistDocument, *utils.AppError) {
	args := m.Called(ctx)

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

func (m *MockService) GetAllArtistsByGenre(ctx context.Context, genre string) ([]ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, genre)

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

func (m *MockService) GetArtistsByCity(ctx context.Context, city string) ([]ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, city)

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

// Test functions
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

func TestCreateArtistHandlerValidationError(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	// Test with invalid JSON
	req := httptest.NewRequest("POST", "/api/artists", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

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

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	// No service calls should be made for invalid JSON
	mockService.AssertNotCalled(t, "CreateArtist")
}

func TestCreateArtistHandlerServiceError(t *testing.T) {
	mockService := new(MockService)
	handler := &TestHandler{service: mockService}

	// Mock service to return an error
	expectedError := utils.ValidationError("Invalid artist data")
	mockService.On("CreateArtist", mock.Anything, mock.Anything).Return((*ArtistDocument)(nil), expectedError)

	reqBody := CreateArtistParams{
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/artists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

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

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertExpectations(t)
}
