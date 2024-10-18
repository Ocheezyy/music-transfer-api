package controllers

import (
	"net/http"

	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
)

func CreatePlaylist(c *gin.Context) {
	var addPlaylistInput models.AddPlaylistInput

	if err := c.ShouldBindJSON(&addPlaylistInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
