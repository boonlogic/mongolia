package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"testing"
)

type Tenant struct {
	mongolia.DefaultModel `bson:",inline"`
	TenantID              string `json:"tenantId" bson:"tenantId"`
	Name                  string `json:"name" bson:"name"`
}

func NewTenant(name string) *Tenant {
	return &Tenant{
		Name: name,
	}
}

func setup() {
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
			{"tenantName", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	coll.Indexes().CreateOne(context.Background(), model, options.CreateIndexes())
}

func TestSmoke(t *testing.T) {
	setup()

	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	err := mongolia.Connect("mongolia-local", opts)
	require.Nil(t, err)

	_ = NewTenant("dev-tenant")

	_, err = ioutil.ReadFile("mongolia/tenant.json")
	if err != nil {
		panic(err)
	}
}
