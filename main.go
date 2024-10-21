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
	router := gin.Default()

	router.POST("/auth/signup", controllers.CreateUser)
	router.POST("/auth/login", controllers.Login)
	router.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	router.POST("/playlist", middlewares.CheckAuth, controllers.CreatePlaylist)
	router.GET("/playlist/:id", middlewares.CheckAuth, controllers.GetPlaylist)
	// router.GET("/playlists", middlewares.CheckAuth, controllers.)
	router.Run()
}
