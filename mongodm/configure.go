package mongodm

import (
	"fmt"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/defaults"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

// Configure configures the global ODM instance.
func Configure(opts *options.ConfigureOptions) error {
	err := opts.Validate()
	if err != nil {
		return err
	}
	return configure(opts)
}

func configure(opts *options.ConfigureOptions) error {
	db, err := connectMongo(opts)
	if err != nil {
		return err
	}
	odm.db = db
	odm.ephemeral = ephemeral(opts)
	odm.schemas = make(map[string]*Schema)
	return nil
}

func ephemeral(opts *options.ConfigureOptions) bool {
	var environment = defaults.ENVIRONMENT
	if opts.Environment != nil {
		environment = *opts.Environment
	}
	var eph bool
	switch environment {
	case options.Development:
		eph = true
	case options.Testing:
		eph = true
	case options.Production:
		eph = false
	}
	return eph
}

func connectMongo(opts *options.ConfigureOptions) (*mongo.Database, error) {
	uri := mongoURI(opts)
	mopts := moptions.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx(), mopts)
	if err != nil {
		return nil, err
	}
	return client.Database(*opts.Database), nil
}

func mongoURI(opts *options.ConfigureOptions) string {
	var (
		cloud = defaults.ON_CLOUD
		port  = defaults.MONGO_PORT
	)
	if opts.OnCloud != nil && *opts.OnCloud {
		cloud = *opts.OnCloud
	}
	if opts.Port != nil {
		port = *opts.Port
	}
	return fmt.Sprintf("%s://%s:%s", protocol(cloud), *opts.Host, port)
}

func protocol(cloud bool) string {
	var proto = "mongodb"
	if cloud {
		proto += "[srv]"
	}
	return proto
}
