package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ocheezyy/music-transfer-api/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		db.Where("ID=?", claims["id"]).Find(&user)

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("currentUser", user)

		c.Next()
	}
}
