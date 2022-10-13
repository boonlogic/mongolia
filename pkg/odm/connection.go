package odm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// Connect connects to the mongo database.
func Connect() error {
	opts := options.Client().ApplyURI(config.URI)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return err
	}
	db = client.Database(config.DBName)
	return nil
}

// Schema adds a collection with the correct indexes for the given document.
// Creates validators for putting documents in/out of the collections.
func AddSchema(schema *Schema) error {
	return nil
}
