package controllers

import (
	"fmt"
	"net/http"

	"github.com/Ocheezyy/music-transfer-api/helpers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/Ocheezyy/music-transfer-api/producers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlaylistTransferController struct {
	DB *gorm.DB
}

func NewPlaylistTransferController(db *gorm.DB) *PlaylistTransferController {
	return &PlaylistTransferController{DB: db}
}

func (ptc *PlaylistTransferController) TransferPlaylist(c *gin.Context) {
	logMethod := "TransferPlaylist"
	var transferPlaylistInput models.TransferPlaylistInput

	if err := c.ShouldBindJSON(&transferPlaylistInput); err != nil {
		errMsg := err.Error()
		helpers.HttpLogBadRequest(logMethod, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	user, ok := helpers.AssertUser(c)
	if !ok {
		helpers.HttpLogISR(logMethod, "failed to assert user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var playlistFound models.Playlist
	ptc.DB.Where(
		"id=? AND user_id=?", transferPlaylistInput.PlaylistID, user.ID,
	).Find(&playlistFound)

	if playlistFound.ID == 0 {
		helpers.HttpLogNotFound(logMethod, "playlist not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}

	var songs []models.Song
	ptc.DB.Where(
		"playlist_id=?", transferPlaylistInput.PlaylistID,
	).Find(&songs)

	for _, song := range songs {
		songMessage := models.SongMessage{
			SongID: song.ID,
			ISRC:   song.ISRC,
			UserID: user.ID,
		}

		err := producers.PublishSongTransfer(songMessage)
		if err != nil {
			// TODO: Handle retry
			helpers.CoreLogError(logMethod, fmt.Sprintf("failed to publish songMessage song id: %d", song.ID), false)
		}
	}
}
