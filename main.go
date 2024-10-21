package main

import (
	"github.com/Ocheezyy/music-transfer-api/controllers"
	"github.com/Ocheezyy/music-transfer-api/initializers"
	"github.com/Ocheezyy/music-transfer-api/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()

	db := initializers.DB

	authMiddleware := middlewares.AuthMiddleware(db)

	playlistController := controllers.NewPlaylistController(db)
	authController := controllers.NewAuthController(db)

	r.POST("/auth/signup", authController.CreateUser)
	r.POST("/auth/login", authController.Login)
	r.GET("/user/profile", authMiddleware, authController.GetUserProfile)
	r.POST("/playlist", authMiddleware, playlistController.CreatePlaylist)
	r.GET("/playlist/:id", authMiddleware, playlistController.GetPlaylist)
	// router.GET("/playlists", middlewares.CheckAuth, controllers.)
	r.Run()
}
