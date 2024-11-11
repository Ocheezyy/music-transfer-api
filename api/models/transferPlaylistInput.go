package models

type TransferPlaylistInput struct {
	PlaylistID uint `json:"playlistId"`
	UserID     uint `json:"userId"`
}
