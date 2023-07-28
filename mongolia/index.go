package mongolia

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func Int32OK(value interface{}) (int32, bool) {
	switch v := value.(type) {
	case int:
		return int32(v), true
	case int32:
		return v, true
	case int64:
		return int32(v), true
	case uint:
		return int32(v), true
	case uint32:
		return int32(v), true
	case uint64:
		return int32(v), true
	case float32:
		return int32(v), true
	case float64:
		return int32(v), true
	default:
		// nah
		return int32(0), false
	}
}

func matchSingleIndex(current bson.D, requested bson.D) error {
	for _, requested_index := range requested {
		if strvalue, ok := requested_index.Value.(string); ok {
			if strvalue == "text" {
				//mongo doesn't store text index fields, so we can't compare
				continue
			}
		}

		for _, current_index := range current {
			if current_index.Key == requested_index.Key {
				request_val, int32_ok := Int32OK(requested_index.Value) //check if requested can be converted to int32 first
				if reflect.TypeOf(current_index.Value).Kind() == reflect.Int32 && int32_ok {
					// 1 or -1 are int32 and the requested bson needs to be converted to compare with current
					if current_index.Value != request_val {
						return errors.New(fmt.Sprintf("Requested Index Field %v has a value (%v) that does not match existing value (%v)", requested_index.Key, request_val, current_index.Value))
					}
				} else if current_index.Value != requested_index.Value {
					return errors.New(fmt.Sprintf("Requested Index Field %v has a value (%v) that does not match existing value (%v)", requested_index.Key, requested_index.Value, current_index.Value))
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
		requested_keys, ok := requested_index.Keys.(bson.D)
		if !ok {
			return errors.New("Requested Index keys must be bson.D")
		}

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
				break
			}

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
	case mongo.IndexModel:
		indexModel = []mongo.IndexModel{indexes.(mongo.IndexModel)}
	default:
		return errors.New(fmt.Sprintf("Unknown index type %v", v))
	}

	//Get current indexes
	current, err := listIndexes(ctx, coll)
	if err != nil {
		return err
	}

	//Compare with existing
	if err = validateIndexes(current, indexModel); err != nil {
		return err
	}

	//Insert indexes
	if err = insertIndexes(ctx, coll, indexModel); err != nil {
		return err
	}

	return nil
}
