package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(config *Config) error {
	opts := options.Client().ApplyURI(*config.URI)
	client, err := mongo.Connect(ctx(), opts)
	if err != nil {
		return err
	}
	odm.db = client.Database(*config.DBName)
	return nil
}

func connectCollection(name string, schema *Schema) (*Collection, error) {
	coll := odm.db.Collection(name)
	if err := prepareIndexes(coll, schema); err != nil {
		return nil, err
	}
	c := &Collection{
		schema: schema,
		coll:   &mgm.Collection{Collection: coll},
	}
	return c, nil
}
