package mongodm

import (
	"context"
	"errors"
	"fmt"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/mongoold"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection is a handle to a mongoold.Collection.
type Collection struct {
	schema *Schema
	coll   *mongoold.Collection
}

func NewCollection(name string, schema *Schema) (*Collection, error) {
	c := &Collection{
		name:   name,
		schema: schema,
	}

	mcoll := odm.db.Collection(name)
	mongoold.MakeCollection(name, schema)
	// todo: take spec definition from schema
	// todo: convert spec into map[string]any (or something parsable)
	// todo: parse expected indexes from spec x-attrs
	// todo: list existing indexes
	// todo: create any expected indexes that are missing

	return c, nil
}

func (c *Collection) getIndexes() []Index {
	return nil
}

func (c *Collection) addIndex(index Index) {
}

func GetCollection(name string) (*Collection, error) {
	coll, ok := odm.colls[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no collection named '%s'", name))
	}
	return coll, nil
}

func (c *Collection) createOne(attrs *Attributes) (*Document, error) {
	doc := &Document{
		id:    primitive.ObjectID{},
		attrs: attrs,
	}
	fn := func(ctx context.Context) error {
		if err := c.insertOne(ctx, doc); err != nil {
			return err
		}
		return nil
	}
	c.schema.runWithHooks(ctx(), fn, doc)
	return doc, nil
}

func (c *Collection) insertOne(ctx context.Context, doc *Document) error {
	id, err := mongoold.InsertOne(c, ctx, bson.M(*doc.attrs))
	if err != nil {
		return err
	}
	doc.id = id
	return nil
}
