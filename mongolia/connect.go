package mongolia

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var g_Context context.Context

func ctx() context.Context {
	return g_Context
}

// connect connects ODM to the underlying mongo instance.
func connect(config *Config) error {
	g_Context, _ := context.WithTimeout(context.Background(), *config.Timeout)
	client, err := mongo.Connect(g_Context, options.Client().ApplyURI(*config.URI))
	if err != nil {
		return err
	}

	odm.client = client
	odm.db = client.Database(*config.DBName)
	return nil
}

// connectCollection connects to a mgm.Collection with the given name and
// does all setup needed to make it a Collection enforcing the Schema.
func connectCollection(name string, schema *Schema) (*Collection, error) {
	coll := odm.db.Collection(name)
	indexes, err := prepareIndexes(coll, schema)
	if err != nil {
		return nil, err
	}
	indexmap := make(map[string]Index)
	for _, idx := range indexes {
		indexmap[idx.Name] = idx
	}
	c := &Collection{
		schema:  schema,
		name:    name,
		coll:    coll,
		indexes: indexmap,
	}
	return c, nil
}
