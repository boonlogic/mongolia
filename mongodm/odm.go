package mongodm

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type ODM struct {
	db        *mongo.Database
	colls     map[string]*mongo.Collection
	ephemeral bool
}
