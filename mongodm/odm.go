package mongodm

import "go.mongodb.org/mongo-driver/mongo"

type ODM struct {
	db        *mongo.Database
	ephemeral bool
	schemas   map[string]*Schema
	colls     map[string]*mongo.Collection
}
