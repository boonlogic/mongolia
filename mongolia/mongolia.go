package mongolia

import (
	"context"
	"time"
)

const (
	defaultURI     = "mongodb://localhost:27017"
	defaultDBName  = "mongolia-local"
	defaultTimeout = 10 * time.Second
)

var odm *ODM

func init() {
	odm = NewODM()
}

func ctx() context.Context {
	return context.Background()
}

func Connect(config *Config) error {
	return odm.Connect(config)
}

func AddSchema(name string, path string) error {
	return odm.AddSchema(name, path)
}

func GetCollection(name string) (*Collection, error) {
	return odm.GetCollection(name)
}
