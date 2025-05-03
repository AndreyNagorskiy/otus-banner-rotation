package amqp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Client основной клиент для работы с AMQP.
type Client struct {
	connManager *connManager
	publisher   *publisher
}

// NewClient создает новый клиент AMQP.
func NewClient(connectionString string) (*Client, error) {
	cm := newConnManager(connectionString)
	if err := cm.connect(); err != nil {
		return nil, err
	}

	p, err := newPublisher(cm, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return &Client{
		connManager: cm,
		publisher:   p,
	}, nil
}

// Publish отправляет сообщение в AMQP.
func (c *Client) Publish(ctx context.Context, exchange, routingKey string, message amqp091.Publishing) error {
	return c.publisher.publish(ctx, exchange, routingKey, message)
}

// PublishJSON отправляет сообщение в AMQP как JSON.
func (c *Client) PublishJSON(ctx context.Context, exchange, routingKey string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.Publish(ctx, exchange, routingKey, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
	})
}

// Close закрывает соединение с AMQP.
func (c *Client) Close() error {
	if err := c.publisher.close(); err != nil {
		return err
	}
	return c.connManager.close()
}

// IsConnected проверяет соединение с AMQP.
func (c *Client) IsConnected() bool {
	return c.connManager.connected()
}
