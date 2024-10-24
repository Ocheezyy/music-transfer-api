package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

type playlistResponse struct {
	Data models.Playlist `json:"data"`
}

func TestGetPlaylist_Success(t *testing.T) {
	db, mock := test.NewMockDB(t)

	playlistName := "TEST_" + gofakeit.State() + gofakeit.Color()
	extPlaylistID := gofakeit.UUID()

	rows := sqlmock.NewRows([]string{"id", "name", "ext_playlist_id", "platform", "song_count", "user_id"}).
		AddRow(1, playlistName, extPlaylistID, models.Spotify, 10, 1)

	mock.ExpectQuery(`SELECT \* FROM "playlists"`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	controller := NewPlaylistController(db)

	// Prepare a Gin context and add the mock user
	c, res := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", 1)}}
	c.Set("currentUser", models.User{ID: 1})

	// Call the handler function
	controller.GetPlaylist(c)

	var playlistRes playlistResponse
	err := json.Unmarshal(res.Body.Bytes(), &playlistRes)
	if err != nil {
		log.Printf("Failed to unmarshal response: %s", err)
	}
	assert.Nil(t, err)
	playlist := playlistRes.Data

	// Assert the response
	assert.Equal(t, playlistName, playlist.Name)
	assert.Equal(t, extPlaylistID, playlist.ExtPlaylistID)
	assert.Equal(t, uint(1), playlist.ID)
	assert.Equal(t, uint(1), playlist.UserID)
	assert.Equal(t, models.Spotify, playlist.Platform)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), `"id":1`)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

func TestGetPlaylist_NotFound(t *testing.T) {
	db, mock := test.NewMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "playlists"`).
		WithArgs(uint(1), uint(1)).
		WillReturnError(sql.ErrNoRows)

	controller := NewPlaylistController(db)

	c, res := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: "1"}} // Non-existent playlist ID
	c.Set("currentUser", models.User{ID: 1})       // No need for full user object

	controller.GetPlaylist(c)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Contains(t, res.Body.String(), `"error":"playlist not found"`)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

func TestCreatePlaylist_Success(t *testing.T) {
	db, mock := test.NewMockDB(t)

	songCount := uint(gofakeit.Uint32())
	extPlaylistId := gofakeit.UUID()
	playlistPlatform := models.Spotify
	playlistName := "TEST_" + gofakeit.State() + gofakeit.Color()

	body := models.CreatePlaylistInput{
		Name:          playlistName,
		Platform:      playlistPlatform,
		ExtPlaylistID: extPlaylistId,
		SongCount:     songCount,
	}

	mock.ExpectQuery(`SELECT \* FROM "playlists"`).
		WithArgs(extPlaylistId, 1).
		WillReturnError(sql.ErrNoRows)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery(`INSERT INTO "playlists"`).
		WithArgs(sqlmock.AnyArg(), uint(1), body.Name, body.ExtPlaylistID, body.Platform, body.SongCount).
		WillReturnRows(rows)

	controller := NewPlaylistController(db)

	c, res := getTestContext()
	c.Set("currentUser", models.User{ID: 1})

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	controller.CreatePlaylist(c)

	var playlistRes playlistResponse
	err := json.Unmarshal(res.Body.Bytes(), &playlistRes)
	if err != nil {
		log.Printf("Failed to unmarshal response: %s", err)
	}
	assert.Nil(t, err)
	playlist := playlistRes.Data

	assert.Equal(t, playlistName, playlist.Name)
	assert.Equal(t, extPlaylistId, playlist.ExtPlaylistID)
	assert.Equal(t, playlistPlatform, playlist.Platform)
	assert.Equal(t, songCount, playlist.SongCount)
	assert.Equal(t, uint(1), playlist.UserID)

	bodyString := res.Body.String()
	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, bodyString, `{"data":{"id":`)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

func TestCreatePlaylist_Conflict(t *testing.T) {
	db, mock := test.NewMockDB(t)

	extPlaylistId := gofakeit.UUID()
	body := models.CreatePlaylistInput{
		Name:          "TEST_" + gofakeit.State() + gofakeit.Color(),
		Platform:      models.Spotify,
		ExtPlaylistID: extPlaylistId,
		SongCount:     uint(gofakeit.Uint32()),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "ext_playlist_id", "platform", "song_count"}).
		AddRow(1, 1, "TEST_"+gofakeit.State()+gofakeit.Color(), extPlaylistId, body.Platform, uint(gofakeit.Uint32()))

	mock.ExpectQuery(`SELECT \* FROM "playlists"`).
		WithArgs(extPlaylistId, uint(1)).
		WillReturnRows(rows)

	controller := NewPlaylistController(db)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Logf("TestCreatePlaylist_Conflict: Failed to marshal request body, %s", err)
	}

	c, res := getTestContext()
	c.Set("currentUser", models.User{ID: 1})

	req, err := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Logf("TestCreatePlaylist_Conflict: Failed to create http request for /playlists, %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreatePlaylist(c)

	bodyString := res.Body.String()
	assert.Equal(t, http.StatusConflict, res.Code)
	assert.Equal(t, bodyString, `{"error":"playlist already exists"}`)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}
