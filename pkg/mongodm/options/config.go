package options

import (
	"errors"
	"os"
)

var config = new(Config)

// Config struct contains extra configuration properties for the mgm package.
type Config struct {
	// URI of MongoDB instance
	URI string

	// Name of database to use in mongo instance
	DBName string
}

func NewConfigFromEnvironment() *ConnectOptions {
}

func newConfigFromEnv() *ConnectOptions {
	config := defaultConfig()
	uri := os.Getenv("AMBER_MONGO_URI")
	if uri == "" {
		panic(errors.New("missing environment variable: AMBER_MONGO_URI"))
	}
	dbname := os.Getenv("AMBER_DB_NAME")
	if dbname == "" {
		panic(errors.New("missing environment variable: AMBER_DB_NAME"))
	}
	config.URI = uri
	config.DBName = dbname
	return config
}
