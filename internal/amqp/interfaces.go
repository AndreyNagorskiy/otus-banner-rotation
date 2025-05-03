package amqp

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

// connection интерфейс для абстракции AMQP соединения.
type connection interface {
	Channel() (*amqp091.Channel, error)
	Close() error
	NotifyClose(receiver chan *amqp091.Error) chan *amqp091.Error
}

// channel интерфейс для абстракции AMQP канала.
type channel interface {
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error
	QueueBind(name, key, exchange string, noWait bool, args amqp091.Table) error
	Close() error
	NotifyClose(c chan *amqp091.Error) chan *amqp091.Error
}

type PublisherClient interface {
	Publish(ctx context.Context, exchange, routingKey string, message amqp091.Publishing) error
	PublishJSON(ctx context.Context, exchange, routingKey string, data any) error
	Close() error
}
