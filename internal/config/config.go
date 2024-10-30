package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	QueueParams QueueConfig `mapstructure:"queue"`
	InputParams InputConfig `mapstructure:"input"`
	LogParams   LogConfig   `mapstructure:"log"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		var NotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &NotFoundErr) {
			fmt.Println("Config file not found: ", err)
		}
	}

	viper.SetEnvPrefix("mta2amqp")
	viper.AutomaticEnv()

	var cfg Config

	// Unmarshal the config into the struct.
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.setupLog()
	cfg.setupQueue()
	cfg.setupInput()

	return &cfg, nil
}

// setupLog method sets up the logger with the given configuration
func (c *Config) setupLog() {
	if len(c.LogParams.Outputs) == 0 {
		c.LogParams.Outputs = append(c.LogParams.Outputs, LogOutput{Type: "stdout"})
	}

	if logFile := viper.GetString("log_file"); logFile != "" {
		c.LogParams.Outputs = append(c.LogParams.Outputs, LogOutput{Type: "file", Path: logFile})
	}
}

// setupQueue method sets up the Queue with the given configuration
func (c *Config) setupQueue() {
	if c.QueueParams.Type == "" {
		if queueType := viper.GetString("queue_type"); queueType != "" {
			println("Queue Type: %s", queueType)
			c.QueueParams.Type = queueType
		} else {
			c.QueueParams.Type = "rabbitmq"
		}
	}
	if c.QueueParams.Uri == "" {
		if uri := viper.GetString("queue_uri"); uri != "" {
			println("Queue URL: %s", uri)
			c.QueueParams.Uri = uri
		} else {
			c.QueueParams.Uri = "amqp://guest:guest@localhost:5672/"
		}
	}

	if c.QueueParams.Queue == "" {
		if queueName := viper.GetString("rabbitmq_queue"); queueName != "" {
			println("Queue Queue: %s", queueName)
			c.QueueParams.Queue = queueName
		} else {
			c.QueueParams.Queue = "dsnparser"
		}
	}

	if c.QueueParams.Exchange == "" {
		if exchangeName := viper.GetString("rabbitmq_exchange"); exchangeName != "" {
			println("Queue Exchange: %s", exchangeName)
			c.QueueParams.Exchange = exchangeName
		} else {
			c.QueueParams.Exchange = "dsnparser"
		}
	}
}

func (c *Config) setupInput() {
	if c.InputParams.Type == "" {
		if inputType := viper.GetString("input_type"); inputType != "" {
			println("Input Type: %s", inputType)
			c.InputParams.Type = inputType
		} else {
			c.InputParams.Type = "unix"
		}
	}

	if c.InputParams.Path == "" {
		if path := viper.GetString("input_path"); path != "" {
			println("Input Path: ", path)
			c.InputParams.Path = path
		} else {
			c.InputParams.Path = "/var/run/mta2amqp.sock"
		}
	}
}
