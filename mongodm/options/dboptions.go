package options

import (
	"errors"
	"fmt"
	"github.com/THREATINT/go-net"
	"regexp"
	"strings"
)

type DBOptions struct {
	Name *string
	URI  *string
}

// Validate validates the DBOptions. This method will return the first error found.
func (o *DBOptions) Validate() error {
	return o.validate()
}

func (o *DBOptions) validate() error {
	if o.Name == nil {
		return errors.New("name is required")
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

// SetURI specifies the mongoold connection string to use.
func (o *DBOptions) SetName(name string) *DBOptions {
	o.Name = &name
	return o
}

// SetURI specifies the mongoold connection string to use.
func (o *DBOptions) SetURI(uri string) *DBOptions {
	o.URI = &uri
	return o
}

// MongoURI returns a mongoold connection string corresponding to the current options.
func MongoURI(opts *DBOptions) string {
	var (
		protocol = "mongodb"
		port     = uint16(27017)
	)
	if opts.Cloud != nil && *opts.Cloud {
		protocol += "[srv]"
	}
	if opts.Port != nil {
		port = *opts.Port
	}
	return fmt.Sprintf("%s://%s:%d", protocol, *opts.Host, port)
}
