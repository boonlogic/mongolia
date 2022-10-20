package mongodm

import (
	"errors"
	"os"
)

// Config struct contains extra configuration properties for the mgm package.
type Config struct {
	// URI of MongoDB instance
	URI string

	// Name of database to use in mongo instance
	DBName string
}

func NewConfigFromEnvironment() *Config {
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

func defaultConfig() *Config {
	return &Config{
		URI:    "mongodb://localhost:27017",
		DBName: "mongodm-local",
	}
}
