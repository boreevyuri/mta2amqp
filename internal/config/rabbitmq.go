package config

import "errors"

type QueueConfig struct {
	Type     string `mapstructure:"type"`
	Uri      string `mapstructure:"url"`
	Exchange string `mapstructure:"exchange"`
	Queue    string `mapstructure:"queue"`
}

func (R *QueueConfig) Validate() error {
	if R.Type == "" {
		return errors.New("queue type is required")
	}

	if R.Uri == "" {
		return errors.New("queue uri is required")
	}

	if R.Queue == "" {
		return errors.New("queue name is required")
	}

	if R.Exchange == "" {
		return errors.New("queue exchange is required")
	}

	return nil
}

func (R *QueueConfig) QueueType() string {
	return R.Type
}

func (R *QueueConfig) AccessUri() string {
	return R.Uri
}

func (R *QueueConfig) ExchangeName() string {
	return R.Exchange
}

func (R *QueueConfig) QueueName() string {
	return R.Queue
}
