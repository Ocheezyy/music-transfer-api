package controllers

import (
	"net/http"
	"strconv"

	"github.com/Ocheezyy/music-transfer-api/helpers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SongController struct {
	DB *gorm.DB
}

func NewSongController(db *gorm.DB) *SongController {
	return &SongController{DB: db}
}

func (sc *SongController) GetSong(c *gin.Context) {
	logMethod := "GetSong"

	songId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errMsg := err.Error()
		helpers.HttpLogBadRequest(logMethod, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
	}

	var song models.Song
	sc.DB.Where("id=?", songId).Find(&song)

	if song.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": song})
}

func (sc *SongController) CreateSong(c *gin.Context) {
	logMethod := "CreateSong"

	var createSongInput models.CreateSongInput

	if err := c.ShouldBindJSON(&createSongInput); err != nil {
		errMsg := err.Error()
		helpers.HttpLogBadRequest(logMethod, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	var playlistFound models.Playlist
	sc.DB.Where("id=?", createSongInput.PlaylistID).Find(&playlistFound)

	if playlistFound.ID != 0 {
		helpers.HttpLogNotFound(logMethod, "playlist not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "playlist not found"})
		return
	}

	newSong := models.Song{
		SongTitle:  createSongInput.SongTitle,
		ArtistName: createSongInput.ArtistName,
		ISRC:       createSongInput.ISRC,
	}
	sc.DB.Create(&newSong)

	c.JSON(http.StatusCreated, gin.H{"data": newSong})
}

func (sc *SongController) BulkCreateSongs(c *gin.Context) {
	logMethod := "BulkCreateSongs"

	var createSongsInput models.BulkCreateSongInput

	if err := c.ShouldBindJSON(&createSongsInput); err != nil {
		errMsg := err.Error()
		helpers.HttpLogBadRequest(logMethod, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	songsToInsert := createSongsInput.Songs

	if err := helpers.BulkInsertSongs(sc.DB, songsToInsert); err != nil {
		helpers.HttpLogISR(logMethod, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (sc *SongController) DeleteSong(c *gin.Context) {
	logMethod := "DeleteSong"

	var deleteSongInput models.DeleteSongInput

	if err := c.ShouldBindJSON(&deleteSongInput); err != nil {
		errMsg := err.Error()
		helpers.HttpLogBadRequest(logMethod, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	sc.DB.Delete(&models.Song{}, deleteSongInput.ID)
	c.JSON(http.StatusNoContent, gin.H{})
}
