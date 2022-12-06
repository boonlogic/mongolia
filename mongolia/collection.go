package mongolia

import (
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

func (c *Collection) FindByID(id any, model Model) error {
	if err := c.coll.FindOne(ctx(), bson.D{{"_id", id}}).Decode(&model); err != nil {
		return err
	}
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Create(model Model, opts *options.InsertOneOptions) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	_, err := c.coll.InsertOne(ctx(), model, opts)
	if err != nil {
		return err
	}
	return nil
}

func (c *Collection) FindOne(filter any, model Model, opts *options.FindOneOptions) error {
	if err := c.coll.FindOne(ctx(), filter, opts).Decode(&model); err != nil {
		return err
	}
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	return nil
}

func (c *Collection) FindMany(filter any, models []Model, opts *options.FindOptions) error {
	cursor, err := c.coll.Find(ctx(), filter, opts)
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

func (c *Collection) Update(model Model, opts *options.UpdateOptions) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	_, err := c.coll.UpdateOne(ctx(), bson.M{"_id": model.GetID()}, bson.M{"$set": model}, opts)
	if err != nil {
		return err
	}
	return nil
}

func (c *Collection) Delete(model Model) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}

	_, err := c.coll.DeleteOne(ctx(), bson.M{"_id": model.GetID()})
	if err != nil {
		return err
	}

	return nil
}

func (c *Collection) Drop() error {
	return c.coll.Drop(ctx())
}

func (c *Collection) Aggregate(models []Model, pipeline interface{}) error {
	cursor, err := c.coll.Aggregate(ctx(), pipeline)
	if err != nil {
		return err
	}
	err = cursor.All(ctx(), &models)
	if err != nil {
		return err
	}

	return nil
}
