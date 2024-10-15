package main

import (
	"log"

	"github.com/Ocheezyy/music-transfer-api/initializers"
	"github.com/Ocheezyy/music-transfer-api/models"
)

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()

}

func main() {

	err := initializers.DB.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatalf("failed to migrate: %s", err)
	} else {
		log.Printf("migrated successfully")
	}

}
