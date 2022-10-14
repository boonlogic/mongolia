package odm

import (
	"errors"
	"os"
)

var config *Config

// Config struct contains extra configuration properties for the mgm package.
type Config struct {
	// URI of MongoDB instance
	URI string

	// Name of database to use in mongo instance
	DBName string
}

func Configure() error {
	uri := os.Getenv("AMBER_MONGO_URI")
	if uri == "" {
		return errors.New("missing environment variable: AMBER_MONGO_URI")
	}
	dbname := os.Getenv("AMBER_DB_NAME")
	if dbname == "" {
		return errors.New("missing environment variable: AMBER_DB_NAME")
	}
	config = new(Config)
	config.URI = uri
	config.DBName = dbname
	return nil
}
