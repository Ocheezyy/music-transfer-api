package controllers

import (
	"bytes"
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

func TestCreateSong_Success(t *testing.T) {
	db, mock := test.NewMockDB(t)

	songTitle := "TEST_" + gofakeit.Color() + gofakeit.Name()
	artistName := fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.LastName())
	isrc := gofakeit.UUID()
	playlistId := uint(10)

	body := models.CreateSongInput{
		SongTitle:  songTitle,
		ArtistName: artistName,
		ISRC:       isrc,
		PlaylistID: playlistId,
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery(`INSERT INTO "songs"`).
		WithArgs(
			sqlmock.AnyArg(),
			body.SongTitle,
			body.ArtistName,
			body.ISRC,
			body.PlaylistID,
		).
		WillReturnRows(rows)

	controller := NewSongController(db)

	c, res := getTestContext()
	c.Set("currentUser", models.User{ID: 1})

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/song", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req
	controller.CreateSong(c)

	var songRes songResponse
	err := json.Unmarshal(res.Body.Bytes(), &songRes)
	if err != nil {
		log.Printf("Failed to unmarshal response: %s", err)
	}
	assert.Nil(t, err)
	song := songRes.Data

	assert.Equal(t, songTitle, song.SongTitle)
	assert.Equal(t, artistName, song.ArtistName)
	assert.Equal(t, isrc, song.ISRC)
	assert.Equal(t, playlistId, song.PlaylistID)

	assert.Equal(t, http.StatusCreated, res.Code)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
}

// func TestBulkInsert_Success(t *testing.T) {
// 	db, mock := test.NewMockDB(t)

// 	var body models.BulkCreateSongInput
// 	rows := sqlmock.NewRows([]string{"id", "song_title", "artist_name", "isrc", "playlist_id"})
// 	for i := 0; i < 4; i++ {
// 		songTitle := "TEST_" + gofakeit.Color() + gofakeit.Name()
// 		artistName := fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.LastName())
// 		isrc := gofakeit.UUID()
// 		playlistId := uint(10)

// 		newSong := models.Song{
// 			SongTitle:  songTitle,
// 			ArtistName: artistName,
// 			ISRC:       isrc,
// 			PlaylistID: playlistId,
// 		}
// 		body.Songs = append(body.Songs, newSong)
// 		rows.AddRow(uint(i), songTitle, artistName, isrc, playlistId)
// 	}

// 	songs := body.Songs

// 	mock.ExpectBegin()
// 	mock.ExpectQuery(`INSERT INTO "songs"`).WithArgs(
// 		sqlmock.AnyArg(), songs[0].SongTitle, songs[0].ArtistName, songs[0].ISRC, songs[0].PlaylistID,
// 		sqlmock.AnyArg(), songs[1].SongTitle, songs[1].ArtistName, songs[1].ISRC, songs[1].PlaylistID,
// 		sqlmock.AnyArg(), songs[2].SongTitle, songs[2].ArtistName, songs[2].ISRC, songs[2].PlaylistID,
// 		sqlmock.AnyArg(), songs[3].SongTitle, songs[3].ArtistName, songs[3].ISRC, songs[3].PlaylistID,
// 	).WillReturnRows(rows)
// 	mock.ExpectCommit()
// 	// .WillReturnResult(sqlmock.NewResult(4, 4))

// 	// mock.ExpectExec(
// 	//

// 	controller := NewSongController(db)

// 	jsonBody, err := json.Marshal(body)
// 	if err != nil {
// 		t.Logf("TestBulkInsert_Success: Failed to marshal request body, %s", err)
// 	}

// 	c, res := getTestContext()
// 	c.Set("currentUser", models.User{ID: 1})

// 	req, err := http.NewRequest("POST", "/songs", bytes.NewBuffer(jsonBody))
// 	if err != nil {
// 		t.Logf("TestBulkInsert_Success: Failed to create http request for /songs, %s", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	c.Request = req

// 	controller.BulkCreateSongs(c)

// 	bodyString := res.Body.String()
// 	assert.Equal(t, http.StatusCreated, res.Code)
// 	assert.Equal(t, bodyString, "{}")

// 	t.Cleanup(func() {
// 		sqlDB, _ := db.DB()
// 		sqlDB.Close()
// 	})
// }
