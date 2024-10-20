package controllers

import (
	"net/http"

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

	// Check that this object is correct
	var user models.User
	user, _ := c.Get("currentUser")

	var playlistFound models.Playlist
	initializers.DB.Where(
		"ext_playlist_id=? AND user_id=?", addPlaylistInput.ExtPlaylistID, addPlaylistInput.UserID
	).Find(&playlistExists)

	if playlistFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error", "playlist already exists"})
		return
	}

	// Now create the playlist
	newPlaylist := models.Playlist{
		UserID: user.ID,
		ExtPlaylistID: addPlaylistInput.ExtPlaylistID,
		Platform: addPlaylistInput.Platform,
		SongCount: addPlaylistInput.SongCount,
	}

	initializers.DB.Create(&newPlaylist)
	c.JSON(http.StatusOK, gin.H{"data": newPlaylist})
	return
}

func GetPlaylist(c *gin.Context) {
	var user models.User
	user, _ := c.Get("currentUser")

	playlistId := c.Param("id")

	var playlist models.Playlist
	initializers.DB.Where(
		"id=? AND user_id=?", playlistId, user.ID
	).Find(&playlist)

	if playlist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": playlist})
}
