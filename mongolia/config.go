package mongolia

import (
	"time"
)

type Config struct {
	URI     *string
	DBName  *string
	Timeout *time.Duration
}

func NewConfig() *Config {
	return new(Config)
}

func DefaultConfig() *Config {
	u := defaultURI
	n := defaultDBName
	t := defaultTimeout
	return &Config{
		URI:     &u,
		DBName:  &n,
		Timeout: &t,
	}
}

func (c *Config) SetURI(uri string) *Config {
	c.URI = &uri
	return c
}

func (c *Config) SetDBName(dbname string) *Config {
	c.DBName = &dbname
	return c
}

func (c *Config) SetTimeout(timeout time.Duration) *Config {
	c.Timeout = &timeout
	return c
}

func (c *Config) Merge(updates *Config) *Config {
	if updates == nil {
		return c
	}
	if updates.URI != nil {
		c.URI = updates.URI
	}
	if updates.DBName != nil {
		c.DBName = updates.DBName
	}
	if updates.Timeout != nil {
		c.Timeout = updates.Timeout
	}
	return c
}
