package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ocheezyy/music-transfer-api/initializers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/Ocheezyy/music-transfer-api/test"
	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestGetPlaylist_Success(t *testing.T) {
	user := test.CreateTestUser()
	playlist := test.CreateTestPlaylist(user.ID)

	// Prepare a Gin context and add the mock user
	c, w := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", playlist.ID)}}
	c.Set("currentUser", user)

	// Call the handler function
	GetPlaylist(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)

	t.Cleanup(func() {
		initializers.DB.Delete(&playlist)
		initializers.DB.Delete(&user)
	})
}

func TestGetPlaylist_NotFound(t *testing.T) {
	c, res := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: "9000"}} // Non-existent playlist ID
	c.Set("currentUser", models.User{ID: 1})          // No need for full user object

	GetPlaylist(c)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Contains(t, res.Body.String(), `"error":"playlist not found"`)
}

func TestCreatePlaylist_Success(t *testing.T) {
	user := test.CreateTestUser()

	c, res := getTestContext()
	c.Set("currentUser", user)

	songCount := uint(gofakeit.Uint32())
	extPlaylistId := gofakeit.UUID()
	playlistPlatform := models.Spotify
	playlistName := "TEST_" + gofakeit.State() + gofakeit.Color()

	body := models.AddPlaylistInput{
		Name:          playlistName,
		Platform:      playlistPlatform,
		ExtPlaylistID: extPlaylistId,
		SongCount:     songCount,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	CreatePlaylist(c)

	var playlist models.Playlist
	err := json.Unmarshal(res.Body.Bytes(), &playlist)
	assert.Nil(t, err)

	assert.Equal(t, playlist.Name, playlistName)
	assert.Equal(t, playlist.ExtPlaylistID, extPlaylistId)
	assert.Equal(t, playlist.Platform, playlistPlatform)
	assert.Equal(t, playlist.SongCount, songCount)
	assert.Equal(t, playlist.UserID, user.ID)

	bodyString := res.Body.String()
	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, bodyString, `"data":"{"id":`)

	t.Cleanup(func() {
		initializers.DB.Delete(&playlist)
		initializers.DB.Delete(&user)
	})
}

func TestCreatePlaylist_Conflict(t *testing.T) {
	user := test.CreateTestUser()
	playlistOne := test.CreateTestPlaylist(user.ID)

	c, res := getTestContext()
	c.Set("currentUser", user)

	body := models.AddPlaylistInput{
		Name:          playlistOne.Name,
		Platform:      playlistOne.Platform,
		ExtPlaylistID: playlistOne.ExtPlaylistID,
		SongCount:     playlistOne.SongCount,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	CreatePlaylist(c)
	bodyString := res.Body.String()
	assert.Equal(t, http.StatusConflict, res.Code)
	assert.Contains(t, bodyString, `"error": "playlist already exists"`)

	t.Cleanup(func() {
		initializers.DB.Delete(&playlistOne)
		initializers.DB.Delete(&user)
	})
}
