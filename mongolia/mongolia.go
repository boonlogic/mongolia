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

// Connect establishes ODM's connection with the underlying mongo instance.
func Connect(config *Config) error {
	return odm.Connect(config)
}

// AddSchema adds a Schema to ODM.
//
// When a new Schema is added, a Collection enforcing that schema is added to the ODM.
func AddSchema(name string, path string) (*Collection, error) {
	return odm.AddSchema(name, path)
}

// GetCollection returns a handle to Collection given its name.
func GetCollection(name string) (*Collection, error) {
	return odm.GetCollection(name)
}
