package queues

import (
	"context"
	"io"
)

type kafkaConfig map[string]string

// Kafka example
type Kafka struct {
	cfg kafkaConfig
}

func NewKafka(_ kafkaConfig) *Kafka {
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
