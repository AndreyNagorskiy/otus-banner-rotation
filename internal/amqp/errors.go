package amqp

import "errors"

var (
	ErrNotConnected     = errors.New("not connected to RabbitMQ")
	ErrPublishTimeout   = errors.New("publish timeout")
	ErrConnectionClosed = errors.New("connection closed")
)
