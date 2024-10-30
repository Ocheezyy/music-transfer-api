package helpers

import (
	"github.com/Ocheezyy/music-transfer-api/models"
	"gorm.io/gorm"
)

func BulkInsertSongs(db *gorm.DB, songs []models.Song) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&songs).Error; err != nil {
			return err
		}

		return nil
	})
}
