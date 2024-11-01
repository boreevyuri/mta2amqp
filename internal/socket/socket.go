package socket

import (
	"context"
	"errors"
	"io"
	"mta2amqp/internal/logger"
	"net"
	"os"
)

// Configurator is an interface that provides a method to get the configuration
type Configurator interface {
	Parse() map[string]string
}

type Socket struct {
	cfg      map[string]string
	listener net.Listener
	// mu       *sync.Mutex
}

func NewSocket(config Configurator) *Socket {
	return &Socket{
		cfg: config.Parse(),
		// mu:  &sync.Mutex{},
	}
}

func (s *Socket) Start(ctx context.Context, p func(msg []byte) error) error {
	log := logger.FromContext(ctx)
	err := s.cleanupSocketFile(log)
	if err != nil {
		return err
	}

	err = s.makeSocket(ctx)
	if err != nil {
		return err
	}

	log.Infof("Listening on %s://%s", s.cfg["type"], s.cfg["path"])

	// Handle cancellation of the context
	go s.handleCancel(ctx, log)

	// Start the connection manager
	go s.manageConnection(ctx, log, p)

	return nil
}

func (s *Socket) cleanupSocketFile(log logger.Logger) error {
	if s.cfg["type"] != "unix" {
		return nil
	}

	_, err := os.Stat(s.cfg["path"])
	if os.IsNotExist(err) {
		return nil
	}

	conn, err := net.Dial(s.cfg["type"], s.cfg["path"])
	if err == nil {
		err = conn.Close()
		if err != nil {
			return err
		}
		return errors.New("socket file exists and already in use")
	}

	log.Infof("Socket file %s is not in use, removing it", s.cfg["path"])
	if err = os.Remove(s.cfg["path"]); err != nil {
		return err
	}

	return nil
}

func (s *Socket) makeSocket(ctx context.Context) error {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, s.cfg["type"], s.cfg["path"])
	if err != nil {
		return err
	}

	s.listener = l
	return nil
}

// manageConnection is a method that listens for incoming connections and processes them
func (s *Socket) manageConnection(ctx context.Context, log logger.Logger, f func(msg []byte) error) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			var opErr *net.OpError
			if errors.As(err, &opErr) && opErr.Err.Error() == "use of closed network connection" {
				log.Info("Listener closed, stopping connection manager")
				break
			}
			log.Errorf("Failed to accept connection: %v", err)
			continue
		} else {
			go s.handleConnection(ctx, conn, f)
		}
	}
}

func (s *Socket) handleConnection(ctx context.Context, conn net.Conn, f func(msg []byte) error) {
	log := logger.FromContext(ctx)
	log.Infof("Accepted connection from %s", conn.RemoteAddr())

	// Read the input data
	input, err := io.ReadAll(conn)
	if err != nil {
		log.Errorf("Failed to read input data: %v", err)
		return
	}

	// If we need to make queueing, we can use the following mutex
	// s.mu.Lock()
	// defer s.mu.Unlock()

	if err := f(input); err != nil {
		log.Errorf("Failed to process input data: %v", err)
	} else {
		log.Info("Data processed successfully")
	}
}

func (s *Socket) handleCancel(ctx context.Context, log logger.Logger) {
	<-ctx.Done()
	log.Info("Shutting down socket...")
	if err := s.Close(); err != nil {
		log.Errorf("Failed to close socket: %v", err)
	}
}

func (s *Socket) Close() error {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	// Remove the socket file
	if err := os.Remove(s.cfg["path"]); err != nil {
		return err
	}

	// Close the listener
	if err := s.listener.Close(); err != nil {
		return err
	}

	return nil
}
