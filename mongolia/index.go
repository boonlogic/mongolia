package mongolia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func convertToBsonx(index_type interface{}) interface{} {
	if index_int, ok := index_type.(int); ok {
		return bsonx.Int32(int32(index_int)) // -1 or 1
	} else if index_float, ok := index_type.(float64); ok {
		return bsonx.Int32(int32(index_float)) // -1 or 1
	} else if index_string, ok := index_type.(string); ok {
		return bsonx.String(index_string) // "text" or "2dsphere"
	}
	return nil
}

func convertToInt(setting interface{}) (int32, error) {
	if opt_int, ok := setting.(int); ok {
		return int32(opt_int), nil
	} else if opt_float, ok := setting.(float64); ok {
		return int32(opt_float), nil
	} else if opt_string, ok := setting.(string); ok {
		i, err := strconv.ParseInt(opt_string, 10, 32)
		if err != nil {
			return 0, err
		}
		return int32(i), nil
	}
	return 0, errors.New(fmt.Sprintf("Could not convert to int %v", setting))
}

func convertToBool(setting interface{}) (bool, error) {
	if opt_bool, ok := setting.(bool); ok {
		return opt_bool, nil
	} else if opt_string, ok := setting.(string); ok {
		boolValue, err := strconv.ParseBool(opt_string)
		if err != nil {
			return false, err
		}
		return boolValue, nil
	} else if opt_int, ok := setting.(int); ok {
		return (opt_int != 0), nil
	}
	return false, errors.New(fmt.Sprintf("Could not convert to bool %v", setting))
}

func parseIndex(indexmap map[string]interface{}) bson.D {
	//parse keys
	keys := make(bson.D, 0)
	for field, value := range indexmap {
		value_parsed := convertToBsonx(value)
		if value_parsed == nil {
			log.Printf("Invalid index value: %v \n", value)
		} else {
			key := bson.E{
				Key:   field,
				Value: value_parsed,
			}
			keys = append(keys, key)
		}
	}
	return keys
}

func parseOptions(name string, optionsmap map[string]interface{}) *options.IndexOptions {
	//parse keys
	opts := options.Index().SetName(name)
	for option, setting := range optionsmap {
		if strings.Contains(strings.ToLower(option), "unique") {
			csetting, err := convertToBool(setting)
			if err == nil {
				opts = opts.SetUnique(csetting)
			}
		} else if strings.Contains(strings.ToLower(option), "ttl") {
			csetting, err := convertToInt(setting)
			if err == nil {
				opts = opts.SetExpireAfterSeconds(csetting)
			}
		}
	}
	return opts
}

func convertJsonToIndex(jsonString []byte) ([]mongo.IndexModel, error) {
	indexModel := make([]mongo.IndexModel, 0)
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	if err != nil {
		return nil, err
	}
	for name, jsonbody := range jsonMap {
		//convert to map
		body, ok := jsonbody.(map[string]interface{})
		if !ok {
			return nil, errors.New(fmt.Sprintf("Invalid JSON Body %v ", jsonbody))
		}
		var opts *options.IndexOptions
		var keys bson.D
		for label, jsoncontent := range body {
			//Convert to map
			content, ok := jsoncontent.(map[string]interface{})
			if !ok {
				return nil, errors.New(fmt.Sprintf("Invalid JSON Content %v ", jsoncontent))
			}

			//parse keys and options
			if strings.Contains(strings.ToLower(label), "key") {
				keys = parseIndex(content)
			} else if strings.Contains(strings.ToLower(label), "option") {
				opts = parseOptions(name, content)
			}
		}
		idxm := mongo.IndexModel{
			Keys:    keys,
			Options: opts,
		}
		indexModel = append(indexModel, idxm)
	}
	return indexModel, nil
}

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

func validateIndexes(have []bson.D, want []mongo.IndexModel) error {
	//Todo
	return nil
}

func PopulateIndexes(ctx context.Context, coll *mongo.Collection, indexes interface{}) error {
	fmt.Printf("PopulateIndexes... \n")
	var indexmodel []mongo.IndexModel
	switch v := indexes.(type) {
	case []mongo.IndexModel:
		indexmodel = indexes
	case string:
		if indexModel, err := convertJsonToIndex([]byte(indexes)); err != nil {
			return err
		}
	case []byte:
		if indexModel, err := convertJsonToIndex(indexes); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Unknown index type %v", v))
	}

	fmt.Printf("Requested Indexes: %v \n", indexmodel)

	//Get current indexes
	if current, err := listIndexes(ctx, coll); err != nil {
		return err
	}
	fmt.Printf("Current Indexes: %v \n", current)

	//Compare with existing
	if err := validateIndexes(current, indexmodel); err != nil {
		return err
	}

	//Insert indexes
	if err := insertIndexes(ctx, coll, indexmodel); err != nil {
		return err
	}
	return nil
}
