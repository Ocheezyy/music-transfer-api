package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/Ocheezyy/music-transfer-consumer/types"
	"github.com/streadway/amqp"
)

func TransferSong(msg amqp.Delivery) {
	var playlist types.PlaylistMessage
	if err := json.Unmarshal(msg.Body, &playlist); err != nil {
		log.Printf("Error decoding message: %v", err)
		return
	}

}

// NOTE: Old functions for transferring and entire playlist per message
// func TransferSongsFromPlaylist(msg amqp.Delivery) {
// 	var playlist models.PlaylistMessage
// 	if err := json.Unmarshal(msg.Body, &playlist); err != nil {
// 		log.Printf("Error decoding message: %v", err)
// 		return
// 	}

// 	log.Printf("Starting transfer for playlist %s owned by user %s", playlist.PlaylistID, playlist.UserID)

// 	// Example: Fetch songs (if not provided in the message)
// 	songs := playlist.Songs
// 	if len(songs) == 0 {
// 		// might need to query database or API to fetch songs
// 		// e.g. FetchSongsFromDB(playlist.PlaylistID)
// 	}

// 	// Iterate through the songs and transfer them
// 	for _, song := range songs {
// 		success := transferSong(song)
// 		if success {
// 			log.Printf("Transferred song: %s", song)
// 		} else {
// 			log.Printf("Failed to transfer song: %s", song)
// 			// Optionally: handle retries, logging, or DLQ (dead letter queue)
// 		}
// 	}

// 	log.Printf("Completed transfer for playlist %s", playlist.PlaylistID)
// }

// func transferSong(song string) bool {
// 	// This function handles the logic to transfer a single song
// 	// This could involve API calls to the target service (e.g., Spotify API)
// 	log.Printf("Transferring song: %s", song)

// 	// Mock transfer success (replace with actual API logic)
// 	return true
// }
