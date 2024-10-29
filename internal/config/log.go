package config

import (
	"io"
	"os"
)

type LogConfig struct {
	Level   string      `mapstructure:"level"`
	Outputs []LogOutput `mapstructure:"outputs"`
}

type LogOutput struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path,omitempty"`
}

// GetLogLevel returns the log level
func (c *LogConfig) GetLogLevel() string {
	return c.Level
}

// GetWriters returns the writers for the logger
func (c *LogConfig) GetWriters() ([]io.Writer, error) {
	var writers []io.Writer
	for _, output := range c.Outputs {
		switch output.Type {
		case "stdout":
			writers = append(writers, os.Stdout)
		case "file":
			file, err := os.OpenFile(output.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}
			writers = append(writers, file)
		}
	}
	return writers, nil
}
