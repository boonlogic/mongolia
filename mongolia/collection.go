package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection exposes CRUD operations for Model, while disallowing
// operations that would cause a document to violate the Schema.
type Collection struct {
	schema  *Schema
	coll    *mgm.Collection
	indexes map[string]Index
}

func (c *Collection) FindByID(id any, model Model) error {
	if err := c.coll.FindByID(id, model); err != nil {
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
	if err := c.coll.Create(model, opts); err != nil {
		return err
	}
	return nil
}

func (c *Collection) First(filter any, model Model, opts *options.FindOneOptions) error {
	if err := c.coll.First(filter, model, opts); err != nil {
		return err
	}
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Update(model Model, opts *options.UpdateOptions) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	if err := c.coll.Update(model, opts); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Delete(model Model) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	if err := c.coll.Delete(model); err != nil {
		return err
	}
	return nil
}

func (c *Collection) drop() error {
	return c.coll.Drop(ctx())
}
