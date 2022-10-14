package odm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// Connect connects to the mongo database.
func Connect() error {
	opts := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(ctx(), opts)
	if err != nil {
		return err
	}
	db = client.Database(config.DBName)
	return nil
}

func Drop() error {
	return db.Drop(ctx())
}