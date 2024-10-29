package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/Ocheezyy/music-transfer-api/test"
	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type songResponse struct {
	Data models.Song `json:"data"`
}

func TestGetSong_Success(t *testing.T) {
	db, mock := test.NewMockDB(t)

	songTitle := "TEST_" + gofakeit.Color() + gofakeit.Name()
	artistName := fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.LastName())
	isrc := gofakeit.UUID()
	playlistId := uint(10)

	rows := sqlmock.NewRows([]string{"id", "song_title", "artist_name", "isrc", "playlist_id"}).
		AddRow(1, songTitle, artistName, isrc, playlistId)

	mock.ExpectQuery(`SELECT \* FROM "songs"`).
		WithArgs(1).
		WillReturnRows(rows)

	controller := NewSongController(db)

	c, res := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", 1)}}
	c.Set("currentUser", models.User{ID: 1})

	controller.GetSong(c)

	var songRes songResponse
	err := json.Unmarshal(res.Body.Bytes(), &songRes)
	if err != nil {
		log.Printf("Failed to unmarshal response %s", err)
	}
	assert.Nil(t, err)
	song := songRes.Data

	assert.Equal(t, songTitle, song.SongTitle)
	assert.Equal(t, artistName, song.ArtistName)
	assert.Equal(t, isrc, song.ISRC)
	assert.Equal(t, playlistId, song.PlaylistID)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

func TestGetSong_NotFound(t *testing.T) {
	db, mock := test.NewMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "songs"`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	controller := NewSongController(db)

	c, res := getTestContext()
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("currentUser", models.User{ID: 1})

	controller.GetSong(c)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Contains(t, res.Body.String(), `"error":"song not found"`)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}
