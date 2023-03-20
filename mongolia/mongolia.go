package mongolia

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ODM struct {
	URI, DB  string
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
	ctx, _ := context.WithTimeout(context.Background(), odm.timeout)
	var err error

	// create mongo connection
	odm.client, err = mongo.Connect(ctx, options.Client().ApplyURI(odm.URI))
	if err != nil {
		return NewError(500, err)
	}

	// test connection
	err = odm.client.Ping(ctx, nil)
	if err != nil {
		errorString := fmt.Sprintf("Error unable to connect to Mongo: %v", err)
		return NewErrorString(500, errorString)
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
		name:    name,
		coll:    coll,
		timeout: odm.timeout,
	}
	odm.colls[name] = *c
	if indexes != nil {
		err := c.CreateIndexes(indexes)
		if err != nil {
			log.Printf("Error Creating Indexes: %v\n", err.ToString())
		}
	}
	return c
}

func (odm *ODM) CreateTimeSeriesCollection(name string, opts *options.TimeSeriesOptions, indexes interface{}) *Collection {
	//Specify a timeseries collection
	col_opts := options.CreateCollection().SetTimeSeriesOptions(opts)
	ctx, _ := context.WithTimeout(context.Background(), odm.timeout)
	err := odm.database.CreateCollection(ctx, name, col_opts)
	if err != nil {
		switch e := err.(type) {
		case mongo.CommandError: // raises a specific CommandError if collection already exists
			if e.Name != "NamespaceExists" {
				log.Printf("Error Creating TimeSeries %v\n", err.Error())
				return nil
			}
		default:
			log.Printf("Error Creating TimeSeries %v\n", err.Error())
			return nil
		}
	}

	c := &Collection{
		name:    name,
		coll:    odm.database.Collection(name),
		timeout: odm.timeout,
	}
	odm.colls[name] = *c
	if indexes != nil {
		err := c.CreateIndexes(indexes)
		if err != nil {
			log.Printf("Error Creating Indexes: %v\n", err.ToString())
		}
	}
	return c
}

func (odm *ODM) Disconnect() {
	odm.client.Disconnect(context.Background())
}

// Drop deletes all ODM data.
// It fails if ODM is not ephemeral.
func (odm *ODM) Drop() {
	odm.database.Drop(context.Background())
}
