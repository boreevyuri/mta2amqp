package logger

import (
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type ZeroLogger struct{}

func (l *ZeroLogger) Debugf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

func (l *ZeroLogger) Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func (l *ZeroLogger) Warnf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

func (l *ZeroLogger) Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

func (l *ZeroLogger) Fatalf(format string, args ...interface{}) {
	log.Fatal().Msgf(format, args...)
}

func (l *ZeroLogger) Debug(msg string) {
	log.Debug().Msg(msg)
}

func (l *ZeroLogger) Info(msg string) {
	log.Info().Msg(msg)
}

func (l *ZeroLogger) Warn(msg string) {
	log.Warn().Msg(msg)
}

func (l *ZeroLogger) Error(msg string) {
	log.Error().Msg(msg)
}

func (l *ZeroLogger) Fatal(msg string) {
	log.Fatal().Msg(msg)
}
