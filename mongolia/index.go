package mongolia

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Index struct {
	Name   string
	Keys   []IndexKey
	Unique bool
}

type IndexKey struct {
	Field     string
	Ascending bool
}

func (idx *Index) toIndexDoc() bson.D {
	var doc = make(bson.D, len(idx.Keys))
	for i, key := range idx.Keys {
		dir := 1
		if !key.Ascending {
			dir = -1
		}
		doc[i] = bson.E{
			Key:   key.Field,
			Value: dir,
		}
	}
	return doc
}

func indexDocToIndex(doc bson.D) (Index, error) {
	var name string
	var keys bson.D
	var unique bool

	var ok bool
	var v any

	if v, ok = doc.Map()["name"]; !ok {
		return Index{}, errors.New("missing key \"name\"")
	}
	if name, ok = v.(string); !ok {
		return Index{}, errors.New(fmt.Sprintf("\"name\" must be type string, got %T"))
	}
}

func (idx *Index) MongoIndex() (MongoIndex, error) {
	var keys = make([]mongoIndexKey, len(idx.Keys))
	for i, key := range idx.Keys {
		k, err := key.IndexKey()
		if err != nil {
			return MongoIndex{}, err
		}
		keys[i] = k
	}
	index := MongoIndex{
		Name:   idx.Name,
		Keys:   keys,
		Unique: idx.Unique,
	}
	return index, nil
}

type MongoIndex struct {
	Name   string `bson:"name"`
	Keys   []mongoIndexKey
	Unique bool `bson:"omitempty"`
}

func (idx *MongoIndex) Index() (Index, error) {
	var keys = make([]IndexKey, len(idx.Keys))
	for i, key := range idx.Keys {
		k, err := key.IndexKey()
		if err != nil {
			return Index{}, err
		}
		keys[i] = k
	}
	index := Index{
		Name:   idx.Name,
		Keys:   keys,
		Unique: idx.Unique,
	}
	return index, nil
}

type mongoIndexKey struct {
	Field     string `bson:"key"`
	Direction int    `bson:"value"`
}

func (k *mongoIndexKey) IndexKey() (IndexKey, error) {
	var asc bool
	switch k.Direction {
	case 1:
		asc = true
	case -1:
		asc = false
	default:
		return IndexKey{}, errors.New(fmt.Sprintf("direction must be 1 or -1, got %d", k.Direction))
	}
	key := IndexKey{
		Field:     k.Field,
		Ascending: asc,
	}
	return key, nil
}

func listIndexes(coll *mongo.Collection) ([]Index, error) {
	curs, err := coll.Indexes().List(ctx())
	if err != nil {
		return nil, err
	}

	var indexes []MongoIndex
	if err = curs.All(ctx(), &indexes); err != nil {
		return nil, err
	}
}

func getIndex(coll *mongo.Collection, name string) (*Index, error) {
	collnames, err := db.ListCollectionNames(ctx(), bson.M{})
	if err != nil {
		return nil
	}

}
