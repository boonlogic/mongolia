package v0

import (
	options2 "gitlab.boonlogic.com/development/expert/mongolia/mongodm/v0/options"
)

type Model interface {
	Save(opts *options2.SaveOptions) error
	Remove(opts *options2.RemoveOptions) error
	Populate(opts *options2.PopulateOptions) error
}
