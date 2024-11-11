package controllers

import (
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
	logMethod := "CreatePlaylist"
	var createPlaylistInput models.CreatePlaylistInput

	if err := c.ShouldBindJSON(&createPlaylistInput); err != nil {
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
	pc.DB.Where(
		"ext_playlist_id=? AND user_id=?", createPlaylistInput.ExtPlaylistID, user.ID,
	).Find(&playlistFound)

	if playlistFound.ID != 0 {
		helpers.HttpLogConflict(logMethod, "playlist already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "playlist already exists"})
		return
	}

	newPlaylist := models.Playlist{
		UserID:        user.ID,
		Name:          createPlaylistInput.Name,
		ExtPlaylistID: createPlaylistInput.ExtPlaylistID,
		Platform:      createPlaylistInput.Platform,
		SongCount:     createPlaylistInput.SongCount,
	}

	pc.DB.Create(&newPlaylist)
	c.JSON(http.StatusCreated, gin.H{"data": newPlaylist})
}

func (pc *PlaylistController) GetPlaylist(c *gin.Context) {
	logMethod := "GetPlaylist"

	user, ok := helpers.AssertUser(c)
	if !ok {
		helpers.HttpLogISR(logMethod, "failed to assert user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	playlistId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helpers.HttpLogBadRequest(logMethod, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var playlist models.Playlist
	pc.DB.Where(
		"id=? AND user_id=?", playlistId, user.ID,
	).Find(&playlist)

	if playlist.ID == 0 {
		helpers.HttpLogNotFound(logMethod, "playlist not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": playlist})
}
