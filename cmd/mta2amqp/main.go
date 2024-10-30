package main

import (
	"context"
	"mta2amqp/internal/config"
	"mta2amqp/internal/logger"
	"mta2amqp/internal/queues"
	"mta2amqp/internal/socket"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		println("Failed to load config: ", err)
	}
	logger.SetupLogger(&cfg.LogParams)

	log := &logger.ZeroLogger{}

	ctx := logger.WithLogger(context.Background(), log)
	ctx, cancel := context.WithCancel(ctx)

	queueManager, err := queues.SetupConsumer(&cfg.QueueParams)
	if err != nil {
		log.Fatalf("Failed to setup queue: %s ", err.Error())
	}

	if err = queueManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start queue: %s ", err.Error())
	}

	s := socket.NewSocket(&cfg.InputParams)
	err = s.Start(ctx, queueManager.Publish)
	if err != nil {
		log.Fatalf("Failed to start socket: %s ", err.Error())
	}

	// Work with signals to stop the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("Shutting down...")

	cancel()

	time.Sleep(2 * time.Second)
	os.Exit(0)
}
