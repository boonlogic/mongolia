package mongodm

import (
	"fmt"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/defaults"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

// Connect connects to the global ODM instance.
func Connect(opts *options.ConnectOptions) error {
	err := opts.Validate()
	if err != nil {
		return err
	}
	return connect(opts)
}

func connect(opts *options.ConnectOptions) error {
	db, err := connectMongo(opts)
	if err != nil {
		return err
	}
	odm.db = db
	odm.ephemeral = ephemeral(opts)
	odm.schemas = make(map[string]*Schema)
	return nil
}

func ephemeral(opts *options.ConnectOptions) bool {
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

func connectMongo(opts *options.ConnectOptions) (*mongo.Database, error) {
	uri := mongoURI(opts)
	mopts := moptions.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx(), mopts)
	if err != nil {
		return nil, err
	}
	return client.Database(*opts.Database), nil
}

func mongoURI(opts *options.ConnectOptions) string {
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
	return fmt.Sprintf("%s://%s:%d", protocol(cloud), *opts.Host, port)
}

func protocol(cloud bool) string {
	var proto = "mongodb"
	if cloud {
		proto += "[srv]"
	}
	return proto
}
