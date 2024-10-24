package controllers

import (
	"log"
	"net/http"

	"github.com/Ocheezyy/music-transfer-api/helpers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransferLogController struct {
	DB *gorm.DB
}

func NewTransferLogController(db *gorm.DB) *TransferLogController {
	return &TransferLogController{DB: db}
}

func (tc *TransferLogController) CreateTransferLog(c *gin.Context) {
	var createTransferLogInput models.CreateTransferLogInput

	if err := c.ShouldBindBodyWithJSON(&createTransferLogInput); err != nil {
		log.Printf("CreateTransferLog 400: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := helpers.AssertUser(c)
	if !ok {
		log.Printf("CreateTransferLog: Failed to assert user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var playlist models.Playlist
	tc.DB.Where("id=?", createTransferLogInput.PlaylistID).Find(&playlist)
	if playlist.ID == 0 {
		log.Printf("CreateTransferLog: playlist not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "playlist not found"})
		return
	}

	newTransferLog := models.TransferLog{
		UserID:     user.ID,
		PlaylistID: playlist.ID,
		Status:     models.Queued,
	}

	tc.DB.Create(&newTransferLog)
	c.JSON(http.StatusCreated, gin.H{"data": newTransferLog})
}
