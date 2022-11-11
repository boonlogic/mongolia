package v0

import (
	"gitlab.boonlogic.com/development/expert/mongolia/v0/options"
)

type Model interface {
	Save(opts *options.SaveOptions) error
	Remove(opts *options.RemoveOptions) error
	Populate(opts *options.PopulateOptions) error
}
