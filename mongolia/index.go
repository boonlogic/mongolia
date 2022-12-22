package mongolia

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func insertIndexes(ctx context.Context, coll *mongo.Collection, indexes []mongo.IndexModel) error {
	_, err := coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	return nil
}

func listIndexes(ctx context.Context, coll *mongo.Collection) ([]bson.D, error) {
	curs, err := coll.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	var docs []bson.D
	if err = curs.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func dropIndex(coll *mongo.Collection, name string) error {
	if _, err := coll.Indexes().DropOne(context.Background(), name); err != nil {
		return err
	}
	return nil
}

func matchSingleIndex(current bson.D, requested bsonx.Doc) error {
	for _, requested_index := range requested {
		if requested_index.Value.String() == "text" {
			//mongo doesn't store text index fields, so we can't compare
			continue
		}

		for _, current_index := range current {
			if current_index.Key == requested_index.Key {
				if current_index.Value != requested_index.Value {
					return errors.New(fmt.Sprintf("Requested Index Field %v has a value (%v )that does not match existing value (%v)", requested_index.Key, requested_index.Value, current_index.Value))
				}
			}
		}
	}
	return nil
}

func validateIndexes(current []bson.D, requested []mongo.IndexModel) error {
	//loop through requested ones
	for _, requested_index := range requested {
		requested_options := *requested_index.Options
		requested_keys := requested_index.Keys.(bsonx.Doc)

		//Requested Name
		if requested_options.Name == nil {
			continue
		}
		requested_name := *requested_options.Name

		//Unique?
		var requested_unique bool = false
		if requested_options.Unique != nil {
			requested_unique = *requested_options.Unique
		}

		match := false
		for _, current_index := range current {
			current_map := current_index.Map()

			//Current name
			current_name_intf, exists := current_map["name"]
			if !exists {
				//No name so we can't compare
				continue
			}
			current_name := current_name_intf.(string)

			//get keys
			current_keys_intf, exists := current_map["key"]
			var current_keys bson.D
			if exists {
				current_keys = current_keys_intf.(bson.D)
			}

			//get unique
			current_unique_intf, exists := current_map["unique"]
			var current_unique bool = false
			if exists {
				current_unique = current_unique_intf.(bool)
			}

			if current_name == requested_name {
				//compare these two and make sure there is no conflict
				if requested_unique != current_unique {
					return errors.New(fmt.Sprintf("Requested Index %v unique property (%v) does not match existing (%v) ", requested_name, requested_unique, current_unique))
				}

				//Compare fields and values
				err := matchSingleIndex(current_keys, requested_keys)
				if err != nil {
					return err
				}

				//Match was found, break now
				match = true
				break
			}

		}
		if !match {
			log.Printf("Created new index %v \n", requested_name)
		}
	}
	return nil
}

func PopulateIndexes(ctx context.Context, coll *mongo.Collection, indexes interface{}) error {
	var indexModel []mongo.IndexModel
	var err error
	switch v := indexes.(type) {
	case []mongo.IndexModel:
		indexModel = indexes.([]mongo.IndexModel)
	default:
		return errors.New(fmt.Sprintf("Unknown index type %v", v))
	}

	//Get current indexes
	current, err := listIndexes(ctx, coll)
	if err != nil {
		return err
	}

	//Compare with existing
	if err := validateIndexes(current, indexModel); err != nil {
		return err
	}

	//Insert indexes
	if err := insertIndexes(ctx, coll, indexModel); err != nil {
		return err
	}

	return nil
}
