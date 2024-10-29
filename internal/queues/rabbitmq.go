package queues

import (
	"context"
	"io"
	"time"

	"mta2amqp/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	BounceRoutingKey = "email.bounce"
	// ComplaintRoutingKey = "email.complaint"
	// UnsentRoutingKey    = "email.unsent"
	// SuccessRoutingKey   = "email.success"
)

type RabbitMQ struct {
	config  QueueConfigProvider
	conn    *amqp.Connection
	channel *amqp.Channel
	msgs    chan io.Reader
	log     logger.Logger
}

func NewRabbitMQ(config QueueConfigProvider) Consumer {
	return &RabbitMQ{
		config: config,
		msgs:   make(chan io.Reader),
	}
}

// Start method starts the RabbitMQ consumer
func (r *RabbitMQ) Start(ctx context.Context) error {
	var err error
	r.log = logger.FromContext(ctx)
	go r.manageConnection(ctx)
	return err
}

func (r *RabbitMQ) manageConnection(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			r.log.Error("Context canceled, stopping RabbitMQ manager...")
			r.handleClose()
			return
		default:
			err := r.connectAndSetup(ctx)
			if err != nil {
				r.log.Errorf("Error in connection or consuming: %v. Retrying in 5 seconds...", err)
				time.Sleep(5 * time.Second)
			} else {
				go r.waitForConnectionLoss(ctx)
			}
		}
	}
}

func (r *RabbitMQ) connectAndSetup(ctx context.Context) error {
	var err error

	r.conn, err = amqp.Dial(r.config.AccessUri())
	if err != nil {
		return err
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		return err
	}

	r.log.Info("Connected to RabbitMQ")

	if err = r.setupExchangeAndQueues(); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) waitForConnectionLoss(ctx context.Context) {
	notifyClose := make(chan *amqp.Error)
	r.conn.NotifyClose(notifyClose)

	select {
	case <-ctx.Done():
		r.log.Info("Context canceled, stopping RabbitMQ connection loss handler...")
		r.handleClose()
		return
	case err := <-notifyClose:
		if err != nil {
			r.log.Errorf("Connection lost: %v", err)
		}
	}
}

func (r *RabbitMQ) consume(ctx context.Context) error {
	msgs, err := r.channel.Consume(
		r.config.QueueName(),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			r.log.Info("Context canceled, stopping RabbitMQ consumer...")
			return nil
		case msg := <-msgs:
			r.msgs <- &Message{Body: msg.Body}
		}
	}
}

// Deliveries return the channel where the messages are sent
func (r *RabbitMQ) Deliveries() <-chan io.Reader {
	return r.msgs
}

// Publish method publishes a message to the RabbitMQ
func (r *RabbitMQ) Publish(msg []byte) error {
	return r.channel.Publish(
		r.config.ExchangeName(),
		r.config.QueueName(),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
}

func (r *RabbitMQ) Close() error {
	var err error
	if r.channel != nil {
		err = r.channel.Close()
		if err != nil {
			return err
		}
	}
	if r.conn != nil {
		err = r.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMQ) setupExchangeAndQueues() error {
	// Declare the exchange
	err := r.channel.ExchangeDeclare(
		r.config.ExchangeName(),
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return r.declareQueue()
}

func (r *RabbitMQ) declareQueue() error {
	// Declare the queues
	q, err := r.channel.QueueDeclare(
		r.config.QueueName(),
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return r.channel.QueueBind(
		q.Name,
		BounceRoutingKey,
		r.config.ExchangeName(),
		false,
		nil,
	)
}

func (r *RabbitMQ) handleClose() {
	if err := r.Close(); err != nil {
		r.log.Errorf("Error closing RabbitMQ connection: %v", err)
	}
}
