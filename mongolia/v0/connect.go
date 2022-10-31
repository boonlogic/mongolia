package v0

import (
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia/v0/options"
)

// Connect sets up a connection to the global ODM instance.
func Connect(opts *options.ODMOptions) error {
	err := opts.Validate()
	if err != nil {
		return err
	}
	return connect(opts)
}

func connect(opts *options.ODMOptions) error {
	var err error
	if odm, err = NewODM(opts); err != nil {
		return err
	}
	return nil
}
