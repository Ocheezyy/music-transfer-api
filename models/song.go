package models

import "time"

type Song struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	SongTitle  string `json:"songTitle" gorm:"uniqueIndex:idx_song_title_playlist_id"`
	ArtistName string `json:"artistName"`
	ISRC       string `json:"isrc"`
	PlaylistID uint   `json:"playlistId" gorm:"uniqueIndex:idx_song_title_playlist_id"`
	CreatedAt  time.Time
}
