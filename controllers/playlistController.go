package controllers

import (
	"net/http"

	"github.com/Ocheezyy/music-transfer-api/helpers"
	"github.com/Ocheezyy/music-transfer-api/initializers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
)

func CreatePlaylist(c *gin.Context) {
	var addPlaylistInput models.AddPlaylistInput

	if err := c.ShouldBindJSON(&addPlaylistInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := helpers.AssertUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var playlistFound models.Playlist
	initializers.DB.Where(
		"ext_playlist_id=? AND user_id=?", addPlaylistInput.ExtPlaylistID, user.ID,
	).Find(&playlistFound)

	if playlistFound.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "playlist already exists"})
		return
	}

	newPlaylist := models.Playlist{
		UserID:        user.ID,
		ExtPlaylistID: addPlaylistInput.ExtPlaylistID,
		Platform:      addPlaylistInput.Platform,
		SongCount:     addPlaylistInput.SongCount,
	}

	initializers.DB.Create(&newPlaylist)
	c.JSON(http.StatusCreated, gin.H{"data": newPlaylist})
}

func GetPlaylist(c *gin.Context) {
	user, ok := helpers.AssertUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	playlistId := c.Param("id")

	var playlist models.Playlist
	initializers.DB.Where(
		"id=? AND user_id=?", playlistId, user.ID,
	).Find(&playlist)

	if playlist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": playlist})
}
