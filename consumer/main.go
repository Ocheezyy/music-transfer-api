package main

import (
	"log"
	"os"
	"time"

	"github.com/Ocheezyy/music-transfer-consumer/rabbitmq"
)

var queueName string
var amqpURI string
var exchangeName string
var exchangeType string
var bindingKey string

func init() {
	amqpURI = os.Getenv("AMQP_URI")
	exchangeName = os.Getenv("EXCHANGE_NAME")
	exchangeType = os.Getenv("EXCHANGE_TYPE")
	queueName = os.Getenv("QUEUE_NAME")
	bindingKey = ""
}

func main() {
	consumerTag := "music-transfer-consumer-1"
	lifetime := 0

	c, err := rabbitmq.NewConsumer(amqpURI, exchangeName, exchangeType, queueName, bindingKey, consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if lifetime > 0 {
		log.Printf("running for %d", lifetime)
		time.Sleep(time.Duration(lifetime))
	} else {
		log.Printf("running forever")
		select {}
	}

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
}
