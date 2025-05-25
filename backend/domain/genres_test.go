// domain/genres_test.go
package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasGenre(t *testing.T) {
	assert.True(t, HasGenre("rock"))
	assert.True(t, HasGenre("jazz"))
	assert.False(t, HasGenre("invalid-genre"))
	assert.False(t, HasGenre(""))
}

func TestGetAllGenres(t *testing.T) {
	genres := GetAllGenres()
	assert.Greater(t, len(genres), 0)
	assert.Contains(t, genres, "rock")
	assert.Contains(t, genres, "jazz")
}
