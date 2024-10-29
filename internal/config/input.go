package config

type InputConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"listen"`
}

// GetPath returns the listen address
func (c *InputConfig) GetPath() string {
	return c.Path
}

// GetType returns the type of input
func (c *InputConfig) GetType() string {
	return c.Type
}
