package logger

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Configurator interface {
	Parse() []map[string]string
}

type ctxKey struct{}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	return ctx.Value(ctxKey{}).(Logger)
}

func SetupLogger(config Configurator) {
	cfg := config.Parse()
	// Set default log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	writers, err := makeWriters(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}

	// If there are any writers, create a multiwriter
	if len(writers) == 1 {
		log.Logger = log.Level(setLevel(cfg[0]["level"])).Output(writers[0])
	} else if len(writers) > 1 {
		multi := io.MultiWriter(writers...)
		// Пока у всех один уровень логгирования
		log.Logger = log.Level(setLevel(cfg[0]["level"])).Output(multi)
	}

	log.Info().Msg("Logger initialized")
}

func makeWriters(outputs []map[string]string) ([]io.Writer, error) {
	var writers []io.Writer
	for _, output := range outputs {
		switch output["type"] {
		case "stdout":
			writers = append(writers, os.Stdout)
		case "file":
			file, err := os.OpenFile(output["path"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}
			writers = append(writers, file)
		}
	}
	return writers, nil
}

// setLevel sets the log level
func setLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
