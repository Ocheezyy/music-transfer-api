package models

type AddPlaylistInput struct {
	Name          string   `json:"name" binding:"required"`
	ExtPlaylistID string   `json:"extPlaylistId" binding:"required"`
	Platform      Platform `json:"platform" binding:"required"`
	SongCount     uint     `json:"songCount" binding:"required"`
}

type UpdateSongCountInput struct {
	PlaylistId uint `json:"playlistId" binding:"required"`
	SongCount  uint `json:"songCount" binding:"required"`
}
