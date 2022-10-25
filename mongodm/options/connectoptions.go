package options

import (
	"errors"
	"fmt"
	net "github.com/THREATINT/go-net"
	"regexp"
)

const (
	hostnameRegex = "[a-zA-Z0-9-_.]+"
	dbRegex       = "[a-z0-9-]+"
)

// ConnectOptions holds options that configure the connection to the ODM.
// Each option is set through a setter function.
// See each setter doc string for an explanation of the option.
type ConnectOptions struct {
	Environment *Environment
	OnCloud     *bool
	Host        *string
	Port        *uint16
	Database    *string
}

// Connect creates a new ConnectOptions instance.
func Connect() *ConnectOptions {
	return new(ConnectOptions)
}

// Validate validates the connection options. This method will return the first error found.
func (c *ConnectOptions) Validate() error {
	return c.validate()
}

func (c *ConnectOptions) validate() error {
	if c.Host == nil {
		return errors.New("host is required")
	}
	if c.Database == nil {
		return errors.New("database is required")
	}
	if c.Environment == nil {
		return errors.New("environment is required")
	}

	switch {
	case net.IsIPAddr(*c.Host):
	case regexp.MustCompile(hostnameRegex).Match([]byte(*c.Host)):
	default:
		return errors.New(fmt.Sprintf("host must be an IP address or valid hostname", *c.Host, hostnameRegex))
	}

	if !regexp.MustCompile(dbRegex).Match([]byte(*c.Database)) {
		return errors.New(fmt.Sprintf("invalid database '%s', must match regex '%s'", *c.Database, dbRegex))
	}

	return nil
}

// SetEnvironment determines whether mongodm is running in a production, test, or development environment.
// Environment is Production by default.
func (c *ConnectOptions) SetEnvironment(environment Environment) *ConnectOptions {
	c.Environment = &environment
	return c
}

// SetCloud specifies whether the mongo instance is deployed in AtlasDB.
// OnCloud is false by default.
func (c *ConnectOptions) SetCloud(cloud bool) *ConnectOptions {
	c.OnCloud = &cloud
	return c
}

// SetHost specifies mongo instance to connect to.
// Host is required. It must be valid IP address or hostname.
func (c *ConnectOptions) SetHost(host string) *ConnectOptions {
	c.Host = &host
	return c
}

// SetPort specifies the port of the mongo instance to connect to.
// Port is 27017 by default.
func (c *ConnectOptions) SetPort(port uint16) *ConnectOptions {
	c.Port = &port
	return c
}

// SetDatabase specifies the name of the mongo database to use.
// Database is required. It may contain only lower case letters and hyphens.
func (c *ConnectOptions) SetDatabase(database string) *ConnectOptions {
	c.Database = &database
	return c
}
