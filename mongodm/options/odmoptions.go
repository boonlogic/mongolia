package options

import (
	"errors"
	"fmt"
	net "github.com/THREATINT/go-net"
	"regexp"
	"strings"
)

const (
	hostnameRegex = "[a-zA-Z0-9-_.]+"
	nameRegex     = "[a-z0-9-]+"
)

// ODMOptions holds options that configure the connection to the ODM.
// Each option is set through a setter function.
// See each setter doc string for an explanation of the option.
type ODMOptions struct {
	Name      *string
	Host      *string
	Port      *uint16
	Cloud     *bool
	Ephemeral *bool
}

// ODM creates a new ODMOptions instance.
func ODM() *ODMOptions {
	return new(ODMOptions)
}

// Validate validates the ODMOptions. This method will return the first error found.
func (o *ODMOptions) Validate() error {
	return o.validate()
}

func (o *ODMOptions) validate() error {
	if o.Name == nil {
		return errors.New("name is required")
	}
	if o.Host == nil {
		return errors.New("host is required")
	}

	if !regexp.MustCompile(nameRegex).Match([]byte(*o.Name)) {
		return errors.New(fmt.Sprintf("invalid name '%s', must match regex '%s'", *o.Name, nameRegex))
	}

	switch {
	case net.IsIPAddr(*o.Host):
	case regexp.MustCompile(hostnameRegex).Match([]byte(*o.Host)):
	default:
		return errors.New(fmt.Sprintf("host must be an IP address or valid hostname", *o.Host, hostnameRegex))
	}

	if *o.Ephemeral && strings.HasSuffix(*o.Name, "-tmp") {
		return errors.New(fmt.Sprintf(""))
	}

	return nil
}

// SetName specifies the name of the ODM instance.
// Name is required and may only contain lower case letters and hyphens.
func (o *ODMOptions) SetName(name string) *ODMOptions {
	o.Name = &name
	return o
}

// SetHost specifies mongoold instance to connect to.
// Host must be valid IP address or hostname.
func (o *ODMOptions) SetHost(host string) *ODMOptions {
	o.Host = &host
	return o
}

// SetPort specifies the port of the mongoold instance to connect to.
func (o *ODMOptions) SetPort(port uint16) *ODMOptions {
	o.Port = &port
	return o
}

// SetCloud specifies whether the mongoold instance is deployed in AtlasDB.
func (o *ODMOptions) SetCloud(cloud bool) *ODMOptions {
	o.Cloud = &cloud
	return o
}

// SetEphemeral determines whether the ODM will permit the Drop operation.
// Ephemeral may only be true if Name has the suffix "-tmp".
// Ephemeral should be used only for testing and development.
func (o *ODMOptions) SetEphemeral(ephemeral bool) *ODMOptions {
	o.Ephemeral = &ephemeral
	return o
}
