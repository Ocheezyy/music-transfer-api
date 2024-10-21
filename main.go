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

	playlistController := controllers.NewPlaylistController(db)
	authController := controllers.NewAuthController(db)

	r.POST("/auth/signup", authController.CreateUser)
	r.POST("/auth/login", authController.Login)
	r.GET("/user/profile", middlewares.CheckAuth, authController.GetUserProfile)
	r.POST("/playlist", middlewares.CheckAuth, playlistController.CreatePlaylist)
	r.GET("/playlist/:id", middlewares.CheckAuth, playlistController.GetPlaylist)
	// router.GET("/playlists", middlewares.CheckAuth, controllers.)
	r.Run()
}
