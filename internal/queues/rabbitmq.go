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

type rmqConfig map[string]string

type RabbitMQ struct {
	cfg     rmqConfig
	conn    *amqp.Connection
	channel *amqp.Channel
	msgs    chan io.Reader
}

func NewRabbitMQ(config rmqConfig) Consumer {
	return &RabbitMQ{
		cfg:  config,
		msgs: make(chan io.Reader),
	}
}

// Start method starts the RabbitMQ consumer
func (r *RabbitMQ) Start(ctx context.Context) error {
	var err error
	log := logger.FromContext(ctx)
	go r.manageConnection(ctx, log)
	return err
}

func (r *RabbitMQ) manageConnection(ctx context.Context, log logger.Logger) {
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping RabbitMQ manager...")
			// r.handleClose()
			return
		default:
			err := r.connectAndSetup(ctx, log)
			if err != nil {
				log.Errorf("Error in connection or consuming: %v. Retrying in 5 seconds...", err)
				time.Sleep(5 * time.Second)
			} else {
				r.waitForConnectionLoss(ctx, log)
			}
		}
	}
}

func (r *RabbitMQ) connectAndSetup(_ context.Context, log logger.Logger) error {
	var err error

	r.conn, err = amqp.Dial(r.cfg["url"])
	if err != nil {
		return err
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		return err
	}

	log.Info("Connected to RabbitMQ")

	if err = r.setupExchangeAndQueues(); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) waitForConnectionLoss(ctx context.Context, log logger.Logger) {
	notifyClose := make(chan *amqp.Error)
	r.conn.NotifyClose(notifyClose)

	select {
	case <-ctx.Done():
		log.Info("Stopping RabbitMQ connection loss handler...")
		if err := r.Close(); err != nil {
			log.Errorf("Error closing RabbitMQ connection: %v", err)
		}
		return
	case err := <-notifyClose:
		if err != nil {
			log.Errorf("Connection lost: %v", err)
		}
	}
}

func (r *RabbitMQ) consume(ctx context.Context, log logger.Logger) error {
	msgs, err := r.channel.Consume(
		r.cfg["queue"],
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
			log.Info("Stopping RabbitMQ consumer...")
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
		r.cfg["exchange"],
		BounceRoutingKey,
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
		r.cfg["exchange"],
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
		r.cfg["queue"],
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
		r.cfg["exchange"],
		false,
		nil,
	)
}
