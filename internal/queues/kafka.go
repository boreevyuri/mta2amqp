package queues

import (
	"context"
	"io"
)

// Kafka example
type Kafka struct {
	config QueueConfigProvider
}

func NewKafka(_ QueueConfigProvider) *Kafka {
	return &Kafka{}
}

func (k *Kafka) Start(_ context.Context) error {
	return nil
}

func (k *Kafka) Deliveries() <-chan io.Reader {
	return nil
}

// Publish sends a message to the queue
func (k *Kafka) Publish(_ []byte) error {
	return nil
}

func (k *Kafka) Close() error {
	return nil
}
