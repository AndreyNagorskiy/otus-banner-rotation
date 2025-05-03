package amqp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type publisher struct {
	connManager *connManager
	channel     channel
	mu          sync.Mutex
	timeout     time.Duration
}

func newPublisher(connManager *connManager, timeout time.Duration) (*publisher, error) {
	p := &publisher{
		connManager: connManager,
		timeout:     timeout,
	}

	if err := p.setupChannel(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *publisher) setupChannel() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, err := p.connManager.getConnection()
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	p.channel = ch
	return nil
}

func (p *publisher) publish(ctx context.Context, exchange, routingKey string, message amqp091.Publishing) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel == nil {
		return ErrNotConnected
	}

	if p.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.timeout)
		defer cancel()
	}

	err := p.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		message,
	)
	if err != nil {
		// Попытка восстановить канал при ошибке
		if err = p.setupChannel(); err != nil {
			return fmt.Errorf("failed to recover channel: %w", err)
		}
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (p *publisher) close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}
