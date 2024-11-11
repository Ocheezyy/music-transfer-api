package initializers

import (
	"fmt"
	"os"
	"time"

	"github.com/Ocheezyy/music-transfer-api/helpers"
	"github.com/streadway/amqp"
)

var amqpURI string
var exchangeName string
var exchangeType string

var RabbitMQ *amqp.Connection
var Channel *amqp.Channel

func init() {
	amqpURI = os.Getenv("AMQP_URI")
	exchangeName = os.Getenv("EXCHANGE_NAME")
	exchangeType = os.Getenv("EXCHANGE_TYPE")
}

func ConnectMQWithRetry(routingKey string) {
	var err error
	for i := 0; i < 5; i++ {
		err = ConnectMQ(routingKey)
		if err == nil {
			helpers.CoreLogInfo(
				"ConnectMQWithRetry",
				"Connected to RabbitMQ successfully",
			)
			break
		}

		helpers.CoreLogError(
			"ConnectMQWithRetry",
			fmt.Sprintf("Failed to connect to RabbitMQ, attempt %d/5: %v", i+1, err),
			false,
		)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		CloseMQ()
		helpers.CoreLogError(
			"ConnectMQWithRetry",
			"After 5 retries failed to connect to rabbitMQ",
			true,
		)
	}
}

func ConnectMQ(routingKey string) error {
	helpers.CoreLogInfo("ConnectMQ", fmt.Sprintf("dialing %q", amqpURI))

	RabbitMQ, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	helpers.CoreLogInfo("ConnectMQ", "Connection succeeded, getting channel")
	Channel, err := RabbitMQ.Channel()
	if err != nil {
		CloseMQ()
		return fmt.Errorf("channel: %s", err)
	}

	helpers.CoreLogInfo("ConnectMQ", fmt.Sprintf("Channel created, declaring %q Exchange (%q)", exchangeType, exchangeName))
	if err := Channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err)
	}

	return nil
}

func PublishMQ(routingKey string, body string) error {
	if err := Channel.Publish(
		exchangeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %s", err)
	}
	return nil
}

func CloseMQ() {
	if Channel != nil {
		Channel.Close()
	}
	if RabbitMQ != nil {
		RabbitMQ.Close()
	}
}
