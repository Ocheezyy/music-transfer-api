package controllers

import (
	"log"
	"net/http"
	"strconv"

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
	songId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("GetSong: id argument is not a uint")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	var createSongInput models.CreateSongInput

	if err := c.ShouldBindJSON(&createSongInput); err != nil {
		log.Printf("CreateSong 400: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var playlistFound models.Playlist
	sc.DB.Where("id=?", createSongInput.PlaylistID).Find(&playlistFound)

	if playlistFound.ID != 0 {
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

func (sc *SongController) DeleteSong(c *gin.Context) {
	var deleteSongInput models.DeleteSongInput

	if err := c.ShouldBindJSON(&deleteSongInput); err != nil {
		log.Printf("DeleteSong 400: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sc.DB.Delete(&models.Song{}, deleteSongInput.ID)
	c.JSON(http.StatusNoContent, gin.H{})
}
