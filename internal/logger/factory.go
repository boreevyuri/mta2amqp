package logger

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
)

type LogConfigurator interface {
	GetLogLevel() string
	GetWriters() ([]io.Writer, error)
}

type ctxKey struct{}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	return ctx.Value(ctxKey{}).(Logger)
}

func SetupLogger(config LogConfigurator) {
	// Set default log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Set log level from config
	// logLevel := viper.GetString("log_level")
	switch config.GetLogLevel() {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	writers, err := config.GetWriters()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}

	// If there are any writers, create a multiwriter
	if len(writers) == 1 {
		log.Logger = log.Output(writers[0])
	} else if len(writers) > 1 {
		multi := io.MultiWriter(writers...)
		log.Logger = log.Output(multi)
	}

	log.Info().Msg("Logger initialized")
}
