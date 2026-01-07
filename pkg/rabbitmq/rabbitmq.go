package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Client struct {
	ch *amqp.Channel
}

func New(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Client{ch: ch}, nil
}

func (c *Client) Publish(queue string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return c.ch.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
