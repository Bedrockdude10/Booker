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
)

// Mock service
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateArtist(ctx context.Context, params CreateArtistParams) (*ArtistDocument, *utils.AppError) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ArtistDocument), args.Error(1)
}

func TestCreateArtistHandler(t *testing.T) {
	mockService := new(MockService)
	handler := &Handler{service: mockService}

	expectedArtist := &ArtistDocument{
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
	}

	mockService.On("CreateArtist", mock.Anything, mock.Anything).Return(expectedArtist, nil)

	reqBody := CreateArtistParams{
		Name:   "Test Artist",
		Genres: []string{"rock"},
		Cities: []string{"Nashville"},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/artists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateArtist(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}
