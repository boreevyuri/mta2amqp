package config

type DBConfig struct {
	Name string `mapstructure:"name"`
	Uri  string `mapstructure:"uri"`
	Type string `mapstructure:"type"`
	// Host     string `mapstructure:"host"`
	// Port     int    `mapstructure:"port"`
	// User     string `mapstructure:"user"`
	// Password string `mapstructure:"password"`
}

// GetDBName returns the database name
func (c *DBConfig) GetDBName() string {
	return c.Name
}

// GetDBUri returns the database uri
func (c *DBConfig) GetDBUri() string {
	return c.Uri
}

// GetDBType returns the database type
func (c *DBConfig) GetDBType() string {
	return c.Type
}
