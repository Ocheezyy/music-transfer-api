package producers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ocheezyy/music-transfer-api/helpers"
	"github.com/Ocheezyy/music-transfer-api/initializers"
	"github.com/Ocheezyy/music-transfer-api/models"
	"github.com/streadway/amqp"
)

var (
	exchangeName = os.Getenv("EXCHANGE_NAME")
	queueName    = os.Getenv("QUEUE_NAME")
)

func PublishSongTransfer(songMessage models.SongMessage) error {
	ch := initializers.Channel

	// Declare the queue to ensure it exists
	_, err := ch.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	songJSON, err := json.Marshal(songMessage)
	if err != nil {
		return err
	}

	// Publish the message to the queue
	err = ch.Publish(
		exchangeName, // Exchange
		queueName,    // Routing key (queue name)
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        songJSON,
		})
	if err != nil {
		return err
	}

	helpers.CoreLogInfo("PublishMessage", fmt.Sprintf("Message published to queue %s", queueName))
	return nil
}
