package v0

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

// Collection wraps client.Collection while enforcing the Schema.
type Collection struct {
	schema *Schema
	coll   *mongo.Collection
}

func GetCollection(name string) (*Collection, error) {
	coll, ok := odm.colls[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no collection named '%s'", name))
	}
	return coll, nil
}

func newCollection(name string, schema *Schema) (*Collection, error) {
	coll := odm.db.Collection(name)

	// todo: initialize indexes
	listIndexes(ctx(), coll)

	c := &Collection{
		schema: schema,
		coll:   coll,
	}
	return c, nil
}

func initIndexes() error {
	// todo: parse expected indexes from spec x-attrs
	// todo: list existing indexes
	// todo: create any expected indexes that are missing
	return nil
}

//// CreateOne instantiates a new object in the database.
//// Model is a handle to the newly created object.
//func (c *Collection) CreateOne(m) (*Model, error) {
//	return c.createOne(doc)
//}
//
//func (c *Collection) createOne(document Document) (*Model, error) {
//	model := &Model{
//		Document: document,
//		coll:     c,
//	}
//	op := func(ctx context.Context) error {
//		id, err := insertOne(ctx, c.coll, model.Document)
//		if err != nil {
//			return err
//		}
//		model.Document["id"] = id
//		return nil
//	}
//	if err := c.schema.runWithHooks(ctx(), op, model); err != nil {
//		return nil, err
//	}
//	return model, nil
//}
