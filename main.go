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
	songController := controllers.NewSongController(db)

	r.POST("/auth/signup", authController.CreateUser)
	r.POST("/auth/login", authController.Login)
	r.GET("/user/profile", authMiddleware, authController.GetUserProfile)

	r.POST("/playlist", authMiddleware, playlistController.CreatePlaylist)
	r.GET("/playlist/:id", authMiddleware, playlistController.GetPlaylist)

	r.GET("/song/:id", authMiddleware, songController.GetSong)
	r.POST("/song", authMiddleware, songController.CreateSong)
	r.POST("/songs", authMiddleware, songController.BulkCreateSongs)
	r.DELETE("/song", authMiddleware, songController.DeleteSong)

	// router.GET("/playlists", middlewares.CheckAuth, controllers.)
	r.Run()
}
