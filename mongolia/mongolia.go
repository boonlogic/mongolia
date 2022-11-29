package mongolia

import (
	"context"
	"time"
)

const (
	defaultURI       = "mongodb://localhost:27017"
	defaultDBName    = "mongolia-local"
	defaultTimeout   = 10 * time.Second
	defaultEphemeral = false
)

var odm *ODM

func init() {
	odm = NewODM()
}

func ctx() context.Context {
	return context.Background()
}

// Connect establishes a connection with the underlying mongo instance.
func Connect(config *Config) error {
	return odm.Connect(config)
}

// AddSchema adds a Schema to ODM.
func AddSchema(name string, path string) (*Collection, error) {
	return odm.AddSchema(name, path)
}

// GetCollection returns a Collection by schema name.
func GetCollection(name string) (*Collection, error) {
	return odm.GetCollection(name)
}

// Drop deletes all ODM data.
// It fails if ODM is not ephemeral.
func Drop() {
	odm.Drop()
}
