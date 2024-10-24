package models

type CreateSongInput struct {
	SongTitle  string `json:"songTitle" binding:"required"`
	ArtistName string `json:"artistName" binding:"required"`
	PlaylistID uint   `json:"playlistId" binding:"required"`
	ISRC       string `json:"isrc" binding:"required"`
}

type DeleteSongInput struct {
	ID uint `json:"id" binding:"required"`
}
