package config

type LogConfig struct {
	Level   string      `mapstructure:"level"`
	Outputs []LogOutput `mapstructure:"outputs"`
}

type LogOutput struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path,omitempty"`
}

func (c *LogConfig) Parse() []map[string]string {
	config := make([]map[string]string, 0)
	for _, output := range c.Outputs {
		logConf := make(map[string]string)
		logConf["type"] = output.Type
		logConf["path"] = output.Path
		logConf["level"] = c.Level
		config = append(config, logConf)
	}

	return config
}
