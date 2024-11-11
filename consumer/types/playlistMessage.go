package types

type PlaylistMessage struct {
	PlaylistID string   `json:"playlistId"`
	UserID     string   `json:"userId"`
	Songs      []string `json:"songs"`
}
