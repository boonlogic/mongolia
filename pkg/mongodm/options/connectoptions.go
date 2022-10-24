package options

import (
	"errors"
	"fmt"
	"regexp"
)

const dbRegex = "[a-z0-9-]+"

// ConnectOptions contains options to configure a mongodm.Client instance.
// Each option can be set through setter functions.
// See the documentation of each setter function for an explanation of the option.
type ConnectOptions struct {
	Cloud    *bool
	Port     *uint16
	Database *string
}

// Client creates a new ConnectOptions instance.
func Client() *ConnectOptions {
	return new(ConnectOptions)
}

// Validate validates the client options. This method will return the first error found.
func (c *ConnectOptions) Validate() error {
	return c.validate()
}

func (c *ConnectOptions) validate() error {
	if !regexp.MustCompile(dbRegex).Match([]byte(*c.Database)) {
		return errors.New(fmt.Sprintf("database must match regex '%s', got '%s'", dbRegex, *c.Database))
	}
	return nil
}

// SetCloud specifies whether the mongo instance to connect to is deployed in AtlasDB.
func (c *ConnectOptions) SetCloud(cloud bool) *ConnectOptions {
	c.Cloud = &cloud
	return c
}

// SetPort specifies the port of the mongo instance to connect to.
func (c *ConnectOptions) SetPort(port uint16) *ConnectOptions {
	c.Port = &port
	return c
}

// SetDatabase specifies the name of the mongo database to use.
// Setting a database name is required.
func (c *ConnectOptions) SetDatabase(database string) *ConnectOptions {
	c.Database = &database
	return c
}
