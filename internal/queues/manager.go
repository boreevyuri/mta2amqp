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

func SetupConsumer(config QueueConfigProvider) (Consumer, error) {
	switch config.QueueType() {
	case "rabbitmq":
		return NewRabbitMQ(config), nil
	case "kafka":
		return NewKafka(config), nil
	default:
		return nil, errors.New("unknown queue type")
	}
}

type Config struct {
	Queue QueueConfigProvider `mapstructure:"queue"`
}

type QueueConfigProvider interface {
	Validate() error
	QueueType() string
	AccessUri() string
	ExchangeName() string
	QueueName() string
}
