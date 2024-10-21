package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Ocheezyy/music-transfer-api/helpers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlaylistController struct {
	DB *gorm.DB
}

func NewPlaylistController(db *gorm.DB) *PlaylistController {
	return &PlaylistController{DB: db}
}

func (pc *PlaylistController) CreatePlaylist(c *gin.Context) {
	var addPlaylistInput models.AddPlaylistInput

	if err := c.ShouldBindJSON(&addPlaylistInput); err != nil {
		log.Printf("CreatePlaylist 400: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := helpers.AssertUser(c)
	if !ok {
		log.Print("GetPlaylist: failed to assert user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var playlistFound models.Playlist
	pc.DB.Where(
		"ext_playlist_id=? AND user_id=?", addPlaylistInput.ExtPlaylistID, user.ID,
	).Find(&playlistFound)

	if playlistFound.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "playlist already exists"})
		return
	}

	newPlaylist := models.Playlist{
		UserID:        user.ID,
		Name:          addPlaylistInput.Name,
		ExtPlaylistID: addPlaylistInput.ExtPlaylistID,
		Platform:      addPlaylistInput.Platform,
		SongCount:     addPlaylistInput.SongCount,
	}

	pc.DB.Create(&newPlaylist)
	c.JSON(http.StatusCreated, gin.H{"data": newPlaylist})
}

func (pc *PlaylistController) GetPlaylist(c *gin.Context) {
	user, ok := helpers.AssertUser(c)
	if !ok {
		log.Print("GetPlaylist: failed to assert user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	playlistId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("GetPlaylist: id argument is not a uint")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var playlist models.Playlist
	pc.DB.Where(
		"id=? AND user_id=?", playlistId, user.ID,
	).Find(&playlist)

	if playlist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": playlist})
}
