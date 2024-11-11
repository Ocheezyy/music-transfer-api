package controllers

import (
	"net/http"
	"strconv"

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

func (tc *TransferLogController) GetTransferLog(c *gin.Context) {
	logMethod := "GetTransferLog"

	transferLogId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errMsg := err.Error()
		helpers.HttpLogBadRequest(logMethod, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
	}

	var transferLog models.Song
	tc.DB.Where("id=?", transferLogId).Find(&transferLog)

	if transferLog.ID == 0 {
		helpers.HttpLogNotFound(logMethod, "transfer_log not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "transfer_log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transferLog})
}

func (tc *TransferLogController) CreateTransferLog(c *gin.Context) {
	logMethod := "CreateTransferLog"

	var createTransferLogInput models.CreateTransferLogInput

	if err := c.ShouldBindBodyWithJSON(&createTransferLogInput); err != nil {
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

	var playlist models.Playlist
	tc.DB.Where("id=?", createTransferLogInput.PlaylistID).Find(&playlist)
	if playlist.ID == 0 {
		helpers.HttpLogNotFound(logMethod, "playlist not found")
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

func (tc *TransferLogController) UpdateTransferLog(c *gin.Context) {
	logMethod := "UpdateTransferLog"

	var updateTransferLogInput models.UpdateTransferLogInput

	if err := c.ShouldBindJSON(&updateTransferLogInput); err != nil {
		helpers.HttpLogBadRequest(logMethod, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var transferLog models.TransferLog
	tc.DB.Where("id = ?", updateTransferLogInput.ID).Find(&transferLog)

	if transferLog.ID == 0 {
		helpers.HttpLogNotFound(logMethod, "transfer_log not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "transferLog not found"})
		return
	}

	transferLog.Status = updateTransferLogInput.Status
	transferLog.Message = updateTransferLogInput.Message
	tc.DB.Save(&transferLog)
	c.JSON(http.StatusOK, gin.H{})
}
