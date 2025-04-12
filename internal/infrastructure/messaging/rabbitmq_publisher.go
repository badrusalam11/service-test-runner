package messaging

import (
	"github.com/streadway/amqp"
)

// ConnectToRabbitMQ connects to the RabbitMQ server using the provided URL and returns the connection and channel.
func ConnectToRabbitMQ(amqpURL string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	return conn, ch, nil
}

// Publisher defines the interface for a message publisher.
type Publisher interface {
	Publish(message []byte) error
}

// RabbitMQPublisher is a concrete implementation that publishes messages to a specific exchange.
type RabbitMQPublisher struct {
	Channel      *amqp.Channel
	ExchangeName string
}

// NewRabbitMQPublisher creates a new RabbitMQPublisher for an exchange.
func NewRabbitMQPublisher(channel *amqp.Channel, exchangeName string) *RabbitMQPublisher {
	err := channel.ExchangeDeclare(
		exchangeName,
		"fanout", // or "direct"/"topic" as needed
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		panic(err) // handle error appropriately in production
	}

	return &RabbitMQPublisher{
		Channel:      channel,
		ExchangeName: exchangeName,
	}
}

// Publish sends the message to the exchange.
func (p *RabbitMQPublisher) Publish(message []byte) error {
	return p.Channel.Publish(
		p.ExchangeName, // exchange
		"",             // routing key (not used for fanout)
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
