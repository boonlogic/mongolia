package mongodm

import (
	"fmt"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/options"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultProtocol = "mongodb"
	defaultPort     = uint16(27017)
	defaultDatabase = "mongodm-local"
)

var db *mongo.Database

func Connect(host string, opts *options.ConnectOptions) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	uri := mongoURI(host, opts)
	mongoOpts := moptions.Client().ApplyURI(uri)
	mongoClient, err := mongo.Connect(ctx(), mongoOpts)
	if err != nil {
		return err
	}

	dbname := defaultDatabase
	if opts.Database != nil {
		dbname = *opts.Database
	}
	db = mongoClient.Database(dbname)

	return nil
}

func mongoURI(host string, opts *options.ConnectOptions) string {
	var (
		protocol = defaultProtocol
		port     = defaultPort
	)
	if opts.Cloud != nil && *opts.Cloud {
		protocol = "mongodb[srv]"
	}
	if opts.Port != nil {
		port = *opts.Port
	}
	return fmt.Sprintf("%s://%s:%s", protocol, host, port)
}
