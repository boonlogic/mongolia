package mongolia

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ODM struct {
	URI, DB  string
	ctx      context.Context
	client   *mongo.Client
	database *mongo.Database
	colls    map[string]Collection
	timeout  time.Duration
}

func NewODM() *ODM {
	return &ODM{
		URI:     "mongodb://localhost:27017",
		DB:      "mongolia-local",
		timeout: 10 * time.Second,
		colls:   make(map[string]Collection),
	}
}

func (odm *ODM) SetURI(uri string) *ODM {
	odm.URI = uri
	return odm
}

func (odm *ODM) SetDBName(db string) *ODM {
	odm.DB = db
	return odm
}

func (odm *ODM) SetTimeout(timeout time.Duration) *ODM {
	odm.timeout = timeout
	return odm
}

func (odm *ODM) Connect() *Error {
	odm.ctx, _ = context.WithTimeout(context.Background(), odm.timeout)
	var err error
	odm.client, err = mongo.Connect(odm.ctx, options.Client().ApplyURI(odm.URI))
	if err != nil {
		return NewError(500, err)
	}

	odm.database = odm.client.Database(odm.DB)
	return nil
}

func (odm *ODM) GetCollection(name string) *Collection {
	coll, ok := odm.colls[name]
	if !ok {
		return odm.CreateCollection(name, nil)
	}
	return &coll
}

func (odm *ODM) CreateCollection(name string, indexes interface{}) *Collection {
	coll := odm.database.Collection(name)
	c := &Collection{
		name: name,
		coll: coll,
		ctx:  odm.ctx,
	}
	odm.colls[name] = *c
	if indexes != nil {
		err := c.CreateIndexes(indexes)
		log.Printf("Error Creating Indexes: %v\n", err.ToString())
	}
	log.Printf("added collection '%s'", name)
	return c
}

func (odm *ODM) Disconnect() {
	odm.client.Disconnect(odm.ctx)
}

// Drop deletes all ODM data.
// It fails if ODM is not ephemeral.
func (odm *ODM) Drop() {
	odm.database.Drop(odm.ctx)
}
