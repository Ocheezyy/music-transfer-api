package test

import (
	"fmt"
	"os"

	"github.com/Ocheezyy/music-transfer-api/initializers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/brianvoe/gofakeit"
	"golang.org/x/crypto/bcrypt"
)

func CreateSamplePassword() string {
	plainPass := gofakeit.Password(true, true, true, true, false, 20)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Failed to hash password: %s", err)
		os.Exit(1)
	}
	return string(passwordHash)
}

func CreateTestUser() models.User {
	user := models.User{
		Email:    "TEST_" + gofakeit.Email(),
		Password: CreateSamplePassword(),
	}
	initializers.DB.Create(&user)
	return user
}

func CreateTestPlaylist(userID uint) models.Playlist {
	playlist := models.Playlist{
		Name:          "TEST_" + gofakeit.Username() + gofakeit.Color(),
		Platform:      models.Spotify,
		ExtPlaylistID: gofakeit.UUID(),
		SongCount:     uint(gofakeit.Uint32()),
		UserID:        userID,
	}
	initializers.DB.Create(&playlist)
	return playlist
}
