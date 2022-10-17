package odm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func ctx() context.Context {
	return context.Background()
}

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
	allowList := []string{
		"mongolia-dev",
		"mongolia-test",
	}
	for _, s := range allowList {
		if db.Name() == s {
			if err := db.Drop(ctx()); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New(fmt.Sprintf("dropping database '%s' is not allowed", db.Name()))
}
