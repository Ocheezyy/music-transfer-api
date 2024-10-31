package producers

import (
	"log"

	"github.com/streadway/amqp"
)

func PublishMessage(conn *amqp.Connection, queueName string, messageBody []byte) error {
	// Open a channel for the connection
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the queue to ensure it exists
	queue, err := ch.QueueDeclare(
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

	// Publish the message to the queue
	err = ch.Publish(
		"",         // Exchange
		queue.Name, // Routing key (queue name)
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		})
	if err != nil {
		return err
	}

	log.Printf("Message published to queue %s", queueName)
	return nil
}
