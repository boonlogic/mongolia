package mongolia

type Config struct {
	URI    *string
	DBName *string
}

func NewConfig() *Config {
	return new(Config)
}

func DefaultConfig() *Config {
	u := defaultURI
	n := defaultDBName
	return &Config{
		URI:    &u,
		DBName: &n,
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
	return c
}
