package controllers

import (
	"log"
	"net/http"

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
