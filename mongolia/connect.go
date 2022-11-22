package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(config *Config) error {
	opts := options.Client().ApplyURI(*config.URI)
	mgmconf := &mgm.Config{
		CtxTimeout: *config.Timeout,
	}
	if err := mgm.SetDefaultConfig(mgmconf, *config.DBName, opts); err != nil {
		return err
	}
	return nil
}

func connectCollection(name string, schema *Schema) (*Collection, error) {
	coll := mgm.CollectionByName(name)
	indexes, err := prepareIndexes(coll.Collection, schema)
	if err != nil {
		return nil, err
	}
	indexmap := make(map[string]Index)
	for _, idx := range indexes {
		indexmap[idx.Name] = idx
	}
	c := &Collection{
		schema:  schema,
		coll:    coll,
		indexes: indexmap,
	}
	return c, nil
}
