package models

type Playlist struct {
	ID            uint     `json:"id" gorm:"primary_key"`
	Name          string   `json:"name"`
	Platform      Platform `json:"platform"`
	ExtPlaylistID string   `json:"extPlaylistId" gorm:"uniqueIndex:idx_user_id_playlist_id"`
	SongCount     uint     `json:"songCount"`
	UserID        uint     `json:"userId" gorm:"uniqueIndex:idx_user_id_playlist_id"`
}

type Platform string

const (
	Apple   Platform = "APPLE"
	Spotify Platform = "SPOTIFY"
	Youtube Platform = "YOUTUBE"
)

// func (p *Platform) Scan(value interface{}) error {
// 	*p = Platform(value.([]byte))
// 	return nil
// }

// func (p Platform) Value() (driver.Value, error) {
// 	return string(p), nil
// }
