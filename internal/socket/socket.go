package socket

import (
	"context"
	"io"
	"mta2amqp/internal/logger"
	"net"
	"os"
	"sync"
)

type InputConfigurator interface {
	GetType() string
	GetPath() string
}

type Socket struct {
	config InputConfigurator
	mu     *sync.Mutex
}

func NewSocket(config InputConfigurator) *Socket {
	return &Socket{
		config: config,
		mu:     &sync.Mutex{},
	}
}

func (s *Socket) Start(ctx context.Context, p func(msg []byte) error) error {
	log := logger.FromContext(ctx)
	l, err := net.Listen(s.config.GetType(), s.config.GetPath())
	if err != nil {
		return err
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalf("Failed to remove socket file: %v", err)
		}
	}(s.config.GetPath())
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("Failed to close listener: %v", err)
		}
	}(l)

	for {
		select {
		case <-ctx.Done():
			log.Info("Context canceled, stopping socket listener...")
			return nil
		default:
			conn, err := l.Accept()
			if err != nil {
				log.Errorf("Failed to accept connection: %v", err)
				continue
			}
			go s.handleConnection(ctx, conn, p)
		}
	}
}

func (s *Socket) handleConnection(ctx context.Context, conn net.Conn, f func(msg []byte) error) {
	log := logger.FromContext(ctx)
	email, err := io.ReadAll(conn)
	if err != nil {
		log.Errorf("Failed to read email: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := f(email); err != nil {
		log.Errorf("Failed to process email: %v", err)
	} else {
		log.Info("Email processed successfully")
	}
}
