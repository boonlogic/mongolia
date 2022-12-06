package mongolia

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection exposes CRUD operations for Model, while disallowing
// operations that would cause a document to violate the Schema.
type Collection struct {
	schema  *Schema
	name    string
	coll    *mongo.Collection
	indexes map[string]Index
}

func (c *Collection) Create(model Model, opts *options.InsertOneOptions) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}

	if err := beforeCreateHooks(ctx(), model); err != nil {
		return err
	}

	res, err := c.coll.InsertOne(ctx(), model, opts)
	if err != nil {
		return err
	}

	// Set new id
	model.SetID(res.InsertedID)

	return afterCreateHooks(ctx(), model)
}

func (c *Collection) FindByID(id any, model Model) error {
	idp, err := model.PrepareID(id)
	if err != nil {
		return err
	}

	if err := c.coll.FindOne(ctx(), bson.D{{"_id", idp}}).Decode(model); err != nil {
		return err
	}
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	return nil
}

func (c *Collection) FindOne(filter any, model Model, opts *options.FindOneOptions) error {
	if err := c.coll.FindOne(ctx(), filter, opts).Decode(model); err != nil {
		return err
	}
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Find(filter any, models []Model, opts *options.FindOptions) (*FindResult, error) {
	cursor, err := c.coll.Find(ctx(), filter, opts)
	if err != nil {
		return nil, err
	}

	var findResult FindResult
	countopts := options.Count().SetMaxTime(2 * time.Second)
	filtered, err := c.coll.CountDocuments(ctx(), filter, countopts)
	if err != nil {
		return nil, err
	}
	findResult.Filtered = filtered

	collection, err := c.coll.CountDocuments(ctx(), bson.D{}, countopts)
	if err != nil {
		return nil, err
	}
	findResult.Collection = collection

	err = cursor.All(ctx(), models)
	if err != nil {
		return nil, err
	}

	for _, model := range models {
		if err := c.schema.Validate(model); err != nil {
			return nil, err
		}
	}
	findResult.Limit = int64(len(models))

	return &findResult, nil
}

func (c *Collection) Update(model Model, opts *options.UpdateOptions) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}

	if err := beforeUpdateHooks(ctx(), model); err != nil {
		return err
	}

	res, err := c.coll.UpdateOne(ctx(), bson.M{"_id": model.GetID()}, bson.M{"$set": model}, opts)
	if err != nil {
		return err
	}

	return afterUpdateHooks(ctx(), res, model)
}

func (c *Collection) Delete(filter any, model Model) error {
	if err := beforeDeleteHooks(ctx(), model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(ctx(), filter)
	if err != nil {
		return err
	}

	return afterDeleteHooks(ctx(), res, model)
}

func (c *Collection) DeleteByID(model Model) error {
	if err := beforeDeleteHooks(ctx(), model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(ctx(), bson.M{"_id": model.GetID()})
	if err != nil {
		return err
	}

	return afterDeleteHooks(ctx(), res, model)
}

func (c *Collection) Drop() error {
	return c.coll.Drop(ctx())
}

func (c *Collection) Aggregate(models []Model, pipeline any) error {
	cursor, err := c.coll.Aggregate(ctx(), pipeline)
	if err != nil {
		return err
	}
	err = cursor.All(ctx(), &models)
	if err != nil {
		return err
	}

	for _, model := range models {
		if err := c.schema.Validate(model); err != nil {
			return err
		}
	}

	return nil
}
