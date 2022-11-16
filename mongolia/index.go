package mongolia

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Index struct {
	Name   string
	Keys   []IndexKey
	Unique bool
}

func (idx *Index) Equals(other Index) bool {
	if idx.Name != other.Name {
		return false
	}
	for i, key := range idx.Keys {
		if !key.Equals(other.Keys[i]) {
			return false
		}
	}
	if idx.Unique != other.Unique {
		return false
	}
	return true
}

type IndexKey struct {
	Field string
	Type  IndexType
}

type IndexType string

const (
	Ascending  = IndexType("asc")
	Descending = IndexType("desc")
	Geospatial = IndexType("2dsphere")
	Text       = IndexType("text")
	Hashed     = IndexType("hashed")
)

func (k *IndexKey) Equals(other IndexKey) bool {
	if k.Field != other.Field {
		return false
	}
	return true
}

func isIndex(doc bson.D) bool {
	name, ok := doc.Map()["name"]
	if !ok {
		return false
	}
	if _, ok := name.(string); !ok {
		return false
	}

	keys, ok := doc.Map()["key"]
	if !ok {
		return false
	}
	keydoc, ok := keys.(bson.D)
	if !ok {
		return false
	}
	for _, elem := range keydoc {
		if elem.Key == "" {
			return false
		}

		switch v := elem.Value.(type) {
		case int32:
			switch v {
			case 1:
			case -1:
			default:
				return false
			}
		case string:
			t := IndexType(v)
			switch t {
			case Ascending:
			case Descending:
			case Geospatial:
			case Text:
			case Hashed:
			default:
				return false
			}
		default:
			return false
		}
	}

	if unique, ok := doc.Map()["unique"]; ok {
		if _, ok := unique.(bool); !ok {
			return false
		}
	}

	return true
}

func toIndex(doc bson.D) Index {
	keydocs := doc.Map()["key"].(bson.D)
	keys := make([]IndexKey, len(keydocs))
	for i, elem := range keydocs {
		var idxtype IndexType
		switch v := elem.Value.(type) {
		case int32:
			switch v {
			case 1:
				idxtype = Ascending
			case -1:
				idxtype = Descending
			}
		case string:
			idxtype = IndexType(v)
		}
		keys[i] = IndexKey{
			Field: elem.Key,
			Type:  idxtype,
		}
	}

	// if the index is unique, it will have a "unique" boolean attribute.
	unique := false
	if u, ok := doc.Map()["unique"]; ok {
		unique = u.(bool)
	}

	idx := Index{
		Name:   doc.Map()["name"].(string),
		Keys:   keys,
		Unique: unique,
	}

	// special case: the default index (name _id_ and an ascending index on the field _id)
	// is always unique, but the attribute "unique" is not returned in its index document.
	if idx.Name == "_id_" && len(idx.Keys) == 1 {
		key := idx.Keys[0]
		if key.Field == "_id" && key.Type == Ascending {
			idx.Unique = true
		}
	}

	return idx
}

func indexName(keys []IndexKey) string {
	name := ""
	for i, k := range keys {
		// mongo uses 1 for ascending and -1 for descending in index names
		var t string
		switch k.Type {
		case Ascending:
			t = "1"
		case Descending:
			t = "-1"
		default:
			t = string(k.Type)
		}
		name += fmt.Sprintf("%s_%d", k.Field, t)
		if i < len(keys)-1 {
			name += "_"
		}
	}
	return name
}

func listIndexes(coll *mongo.Collection) ([]Index, error) {
	curs, err := coll.Indexes().List(ctx())
	if err != nil {
		return nil, err
	}
	var docs []bson.D
	if err = curs.All(ctx(), &docs); err != nil {
		return nil, err
	}
	idxs := make([]Index, len(docs))
	for i, d := range docs {
		if !isIndex(d) {
			return nil, errors.New(fmt.Sprintf("cannot parse mongo document as an index: %+v", d))
		}
		idxs[i] = toIndex(d)
	}
	return idxs, nil
}

func getIndex(coll *mongo.Collection, name string) (*Index, error) {
	return nil, errors.New(fmt.Sprintf("no index named \"%s\"", name))
}

func addIndex(coll *mongo.Collection, index Index) error {
	idxs, err := listIndexes(coll)
	if err != nil {
		return err
	}
	for _, idx := range idxs {
		if index.Name == idx.Name {
			return errors.New(fmt.Sprintf("an index named \"%s\" already exists"))
		}
	}

	keys := make(bson.D, len(idxs))
	for i, k := range index.Keys {
		key := bson.E{
			Key:   k.Field,
			Value: string(k.Type),
		}
		keys[i] = key
	}

	opts := options.Index().
		SetName(index.Name).
		SetUnique(index.Unique)

	idxm := mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	}
	if _, err := coll.Indexes().CreateOne(ctx(), idxm, options.CreateIndexes()); err != nil {
		return err
	}
	return nil
}

func dropIndex(coll *mongo.Collection, name string) error {
	// todo: implement
	return nil
}
