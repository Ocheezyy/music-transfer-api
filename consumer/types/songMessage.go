package types

type SongMessage struct {
	SongID uint   `json:"songId"`
	ISRC   string `json:"isrc"`
	UserID uint   `json:"userId"`
}
