package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	schema *Schema
	coll   *mgm.Collection
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

func (c *Collection) First(filter any, model Model, opts *options.FindOneOptions) error {
	if err := c.coll.First(filter, model, opts); err != nil {
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

func (c *Collection) Update(model Model, opts *options.UpdateOptions) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	if err := c.coll.Update(model, opts); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Save(model Model) error {
	if err := c.schema.Validate(model); err != nil {
		return err
	}
	if err := c.coll.Save(model); err != nil {
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
