package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

//type Tenant struct {
//	mongolia.DefaultModel `bson:",inline"`
//	TenantID              string `json:"tenantId" bson:"tenantId"`
//	Name                  string `json:"name" bson:"name"`
//}
//
//func NewTenant(name string) *Tenant {
//	return &Tenant{
//		Name: name,
//	}
//}

func setup() *mongo.Collection {
	uri := "mongodb://localhost:27017"
	copts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), copts)
	if err != nil {
		panic(err)
	}
	coll := client.Database("mongolia-local").Collection("existing-1")

	coll.Drop(context.Background())

	model := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
			{"tenantId", 1},
		},
		Options: options.Index(),
	}
	coll.Indexes().CreateOne(context.Background(), model, options.CreateIndexes())

	model = mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
			{"name", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	coll.Indexes().CreateOne(context.Background(), model, options.CreateIndexes())
	return coll
}

func TestSmoke(t *testing.T) {
	setup()

	uri := "mongodb://localhost:27017"
	dbname := "mongolia-local"
	timeout := 10 * time.Second

	err := mongolia.Connect(uri, dbname, timeout)
	require.Nil(t, err)

	err = mongolia.AddSchema("tenant", "mongolia/tenant.json")
	require.Nil(t, err)

	//coll, err := mongolia.GetCollection("tenant")
	//require.NotNil(t, coll)
	//require.Nil(t, err)

	//coll, err = mongolia.GetCollection("missing coll")
	//require.Nil(t, coll)
	//require.NotNil(t, err)

	//_ = NewTenant("dev-tenant")
}
