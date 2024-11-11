package helpers

import (
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/gin-gonic/gin"
)

// Get user from gin.Context with type assertion
func AssertUser(c *gin.Context) (models.User, bool) {
	cUser, _ := c.Get("currentUser")
	user, ok := cUser.(models.User)
	return user, ok
}
