package initializers

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

var RabbitMQ *amqp.Connection

func ConnectRabbitMQ(amqpURL string) {
	var err error

	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		RabbitMQ, err = amqp.Dial(amqpURL)
		if err == nil {
			log.Println("Connected to RabbitMQ successfully.")
			break
		}

		log.Printf("Failed to connect to RabbitMQ, attempt %d/5: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("could not connect to RabbitMQ: %w", err)
}
