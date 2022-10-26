package mongoold

import (
	"context"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"
	driver "go.mongodb.org/mongo-driver/mongo"
	driveropts "go.mongodb.org/mongo-driver/mongo/options"
)

// Connect sets up a connection to the global DB instance.
func Connect(ctx context.Context, opts *options.DBOptions) error {
	err := opts.Validate()
	if err != nil {
		return err
	}
	return connect(ctx, opts)
}

func connect(ctx context.Context, opts *options.DBOptions) error {
	dopts := driveropts.Client().ApplyURI(uri)

	client, err := driver.Connect(ctx, dopts)
	if err != nil {
		return err
	}
	if err = NewMongo(opts); err != nil {
		return err
	}
	return nil
}
