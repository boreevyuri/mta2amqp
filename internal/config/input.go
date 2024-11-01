package config

type InputConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

// Parse returns the input configuration
func (c *InputConfig) Parse() map[string]string {
	return map[string]string{
		"type": c.Type,
		"path": c.Path,
	}
}
