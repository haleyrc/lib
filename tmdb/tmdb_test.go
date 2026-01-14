package tmdb_test

import (
	"context"
	"os"
	"testing"

	"github.com/haleyrc/lib/assert"
	"github.com/haleyrc/lib/tmdb"
)

func TestClient_Authenticate(t *testing.T) {
	ctx := context.Background()

	c := newClient(t)
	resp, err := c.Authenticate(ctx)
	assert.OK(t, err).Fatal()

	assert.True(t, "success", resp.Success)
}

func TestClient_SearchMovie(t *testing.T) {
	ctx := context.Background()

	c := newClient(t)
	resp, err := c.SearchMovie(ctx, "Dogma")
	assert.OK(t, err).Fatal()

	first := resp.Results[0]
	assert.Equal(t, "id", 1832, first.ID)
	assert.Equal(t, "original title", "Dogma", first.OriginalTitle)
	assert.Equal(t, "release date", "1999-10-04", first.ReleaseDate)
}

func TestClient_GetMovieDetail(t *testing.T) {
	ctx := context.Background()

	c := newClient(t)
	resp, err := c.GetMovieDetail(ctx, 1832)
	assert.OK(t, err).Fatal()

	assert.Equal(t, "id", 1832, resp.ID)
	assert.Equal(t, "original title", "Dogma", resp.OriginalTitle)
	assert.Equal(t, "release date", "1999-10-04", resp.ReleaseDate)
}

func TestClient_GetMovieCredits(t *testing.T) {
	ctx := context.Background()

	c := newClient(t)
	resp, err := c.GetMovieCredits(ctx, 1832)
	assert.OK(t, err).Fatal()

	assert.Equal(t, "id", 1832, resp.ID)

	first := resp.Cast[0]
	assert.Equal(t, "id", 880, first.ID)
	assert.Equal(t, "name", "Ben Affleck", first.Name)
	assert.Equal(t, "character", "Bartleby", first.Character)
}

func newClient(t *testing.T) *tmdb.Client {
	token := os.Getenv("TMDB_ACCESS_TOKEN")
	if token == "" {
		t.Skip("set TMDB_ACCESS_TOKEN to run this test")
	}
	c, err := tmdb.NewClient(token)
	assert.OK(t, err).Fatal()
	return c
}
