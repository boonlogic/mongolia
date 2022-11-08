package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(dbname string, opts *options.ClientOptions) error {
	client, err := mongo.Connect(ctx(), opts)
	if err != nil {
		return err
	}
	db := client.Database(dbname)

	collnames, err := db.ListCollectionNames(ctx(), bson.M{})
	if err != nil {
		return nil
	}

	for _, name := range collnames {
		coll := mgm.NewCollection(db, name, options.Collection())
		indexes, err := GetIndexes(coll.Collection)
		if err != nil {
			return err
		}
		c := &Collection{
			Collection: coll,
			indexes:    indexes,
		}
		odm.colls[name] = c
	}

	return nil
}
