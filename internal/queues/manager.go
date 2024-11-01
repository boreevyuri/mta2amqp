package queues

import (
	"context"
	"errors"
	"io"
)

type Consumer interface {
	Start(ctx context.Context) error
	Deliveries() <-chan io.Reader
	Publish(msg []byte) error
	Close() error
}

func SetupConsumer(config QueueConfigurator) (Consumer, error) {
	cfg, err := config.Parse()
	if err != nil {
		return nil, err
	}
	switch cfg["type"] {
	case "rabbitmq":
		return NewRabbitMQ(cfg), nil
	case "kafka":
		return NewKafka(cfg), nil
	default:
		return nil, errors.New("unknown queue type")
	}
}

type QueueConfigurator interface {
	Parse() (map[string]string, error)
}
