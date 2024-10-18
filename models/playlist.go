package models

type Playlist struct {
	ID            uint     `json:"id" gorm:"primary_key"`
	Platform      Platform `json:"platform" gorm:"type:enum('APPLE', 'SPOTIFY', 'YOUTUBE')"`
	ExtPlaylistID string   `json:"extPlaylistId" gorm:"unique"`
	SongCount     uint     `json:"songCount"`
	UserID        uint     `json:"userId"`
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
