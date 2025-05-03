package amqp

import (
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// connManager управляет соединением с AMQP.
type connManager struct {
	connectionString string
	conn             connection
	mu               sync.Mutex
	isConnected      bool
	done             chan struct{}
}

func newConnManager(connectionString string) *connManager {
	return &connManager{
		connectionString: connectionString,
		done:             make(chan struct{}),
	}
}

func (cm *connManager) connect() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.isConnected {
		return nil
	}

	conn, err := amqp091.Dial(cm.connectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	cm.conn = conn
	cm.isConnected = true

	go cm.handleReconnect()

	return nil
}

func (cm *connManager) getConnection() (connection, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.isConnected {
		return nil, ErrNotConnected
	}

	return cm.conn, nil
}

// handleReconnect обрабатывает переподключение при разрыве соединения.
func (cm *connManager) handleReconnect() {
	closeChan := make(chan *amqp091.Error)
	cm.conn.NotifyClose(closeChan)

	select {
	case err := <-closeChan:
		if err != nil {
			cm.mu.Lock()
			cm.isConnected = false
			cm.mu.Unlock()

			// Пытаемся переподключиться с экспоненциальной задержкой
			retryDelay := time.Second
			maxRetryDelay := 30 * time.Second

			for {
				time.Sleep(retryDelay)

				if err := cm.connect(); err == nil {
					return
				}

				retryDelay *= 2
				if retryDelay > maxRetryDelay {
					retryDelay = maxRetryDelay
				}
			}
		}
	case <-cm.done:
		return
	}
}

func (cm *connManager) close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.isConnected {
		return nil
	}

	close(cm.done)
	cm.isConnected = false
	return cm.conn.Close()
}

func (cm *connManager) connected() bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.isConnected
}
