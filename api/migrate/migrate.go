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

	err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.Playlist{},
		&models.Task{},
		&models.Song{},
		&models.TransferLog{},
	)

	if err != nil {
		log.Fatalf("failed to migrate: %s", err)
	} else {
		log.Printf("migrated successfully")
	}

}
