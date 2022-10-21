package mongodm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Connect(config *Config) error {
	opts := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(ctx(), opts)
	if err != nil {
		return err
	}
	db = client.Database(config.DBName)
	return nil
}

func ctx() context.Context {
	return context.Background()
}

func drop() {
	allowList := []string{
		"mongodm-local",
	}
	for _, s := range allowList {
		if db.Name() == s {
			if err := db.Drop(ctx()); err != nil {
				panic(err)
			}
			return
		}
	}
	panic(errors.New(fmt.Sprintf("dropping database '%s' is not allowed", db.Name())))
}
