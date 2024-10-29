package config

type WebConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
