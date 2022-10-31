package mongodm

import "gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"

type Model interface {
	Save(opts *options.SaveOptions) error
	Remove(opts *options.RemoveOptions) error
	Populate(opts *options.PopulateOptions) error
}
