package mongodm

import (
	"fmt"
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

	var ephemeral = defaultEphemeral
	if opts.Ephemeral != nil && *opts.Ephemeral {
		ephemeral = *opts.Ephemeral
	}

	odm.db = db
	odm.colls = make(map[string]*mongo.Collection)
	odm.ephemeral = ephemeral

	return nil
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
		cloud = defaultCloud
		port  = defaultPort
	)
	if opts.Cloud != nil && *opts.Cloud {
		cloud = *opts.Cloud
	}
	if opts.Port != nil {
		port = *opts.Port
	}
	return fmt.Sprintf("%s://%s:%s", mongoProtocol(cloud), *opts.Host, port)
}

func mongoProtocol(cloud bool) string {
	var protocol = "mongodb"
	if cloud {
		protocol += "[srv]"
	}
	return protocol
}
