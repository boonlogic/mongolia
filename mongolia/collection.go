package mongolia

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection exposes CRUD operations for Model, while disallowing
// operations that would cause a document to violate the Schema.
type Collection struct {
	name string
	coll *mongo.Collection
	ctx  context.Context
}

func (c *Collection) CreateIndexes(indexes interface{}) error {
	return PopulateIndexes(c.ctx, c.coll, indexes)
}

func (c *Collection) Save(model Model) error {
	if err := model.ValidateRead(); err != nil {
		return err
	}

	if err := beforeSaveHooks(model); err != nil {
		return err
	}

	return afterSaveHooks(model)
}

func (c *Collection) Create(model Model, opts *options.InsertOneOptions) error {
	if err := model.ValidateCreate(); err != nil {
		return err
	}

	if err := beforeCreateHooks(model); err != nil {
		return err
	}

	res, err := c.coll.InsertOne(c.ctx, model, opts)
	if err != nil {
		return err
	}

	// Set new id
	model.SetID(res.InsertedID)

	return afterCreateHooks(model)
}

func (c *Collection) FindByID(id any, model Model) error {
	idp, err := model.PrepareID(id)
	if err != nil {
		return err
	}

	if err := c.coll.FindOne(c.ctx, bson.D{{"_id", idp}}).Decode(model); err != nil {
		return err
	}
	if err := model.ValidateRead(); err != nil {
		return err
	}
	return nil
}

func (c *Collection) FindOne(filter any, model Model, opts *options.FindOneOptions) error {
	if err := c.coll.FindOne(c.ctx, filter, opts).Decode(model); err != nil {
		return err
	}
	if err := model.ValidateRead(); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Find(filter any, models []Model, opts *options.FindOptions) (*FindResult, error) {
	cursor, err := c.coll.Find(c.ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var findResult FindResult
	countopts := options.Count().SetMaxTime(2 * time.Second)
	filtered, err := c.coll.CountDocuments(c.ctx, filter, countopts)
	if err != nil {
		return nil, err
	}
	findResult.Filtered = filtered

	collection, err := c.coll.CountDocuments(c.ctx, bson.D{}, countopts)
	if err != nil {
		return nil, err
	}
	findResult.Collection = collection

	err = cursor.All(c.ctx, models)
	if err != nil {
		return nil, err
	}

	for _, model := range models {
		if err := model.ValidateRead(); err != nil {
			return nil, err
		}
	}
	findResult.Limit = int64(len(models))

	return &findResult, nil
}

func (c *Collection) Update(model Model, opts *options.UpdateOptions) error {
	if err := model.ValidateUpdate(); err != nil {
		return err
	}

	if err := beforeUpdateHooks(model); err != nil {
		return err
	}

	res, err := c.coll.UpdateOne(c.ctx, bson.M{"_id": model.GetID()}, bson.M{"$set": model}, opts)
	if err != nil {
		return err
	}

	return afterUpdateHooks(res, model)
}

func (c *Collection) Delete(filter any, model Model) error {
	if err := beforeDeleteHooks(model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(c.ctx, filter)
	if err != nil {
		return err
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) DeleteByID(id any, model Model) error {
	idp, err := model.PrepareID(id)
	if err != nil {
		return err
	}

	if err := beforeDeleteHooks(model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(c.ctx, bson.M{"_id": idp})
	if err != nil {
		return err
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) DeleteModel(model Model) error {
	if err := beforeDeleteHooks(model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(c.ctx, bson.M{"_id": model.GetID()})
	if err != nil {
		return err
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) Drop() error {
	return c.coll.Drop(c.ctx)
}

func (c *Collection) Aggregate(models []Model, pipeline any) error {
	cursor, err := c.coll.Aggregate(c.ctx, pipeline)
	if err != nil {
		return err
	}
	err = cursor.All(c.ctx, &models)
	if err != nil {
		return err
	}

	for _, model := range models {
		if err := model.ValidateRead(); err != nil {
			return err
		}
	}

	return nil
}
