package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type Client struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func New(
	host string,
	port string,
	user string,
	password string,
	exchange string,
	exchangeType string,
) (*Client, error) {

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		user,
		password,
		host,
		port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,  // durable
		false, // autoDelete
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	return &Client{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}, nil
}

// -------------------------------
// Publish Event
// -------------------------------
func (c *Client) Publish(routingKey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.channel.Publish(
		c.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// -------------------------------
// Consume Event
// -------------------------------
func (c *Client) Consume(
	queue string,
	routingKey string,
	handler func([]byte) error,
) error {

	_, err := c.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err := c.channel.QueueBind(
		queue,
		routingKey,
		c.exchange,
		false,
		nil,
	); err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}()

	return nil
}

// -------------------------------
// Close
// -------------------------------
func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
