package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	*mgm.Collection
}

func (c *Collection) FindByID(id any, model Model) error {
	return c.Collection.FindByID(id, model)
}

func (c *Collection) First(filter any, model Model, opts *options.FindOneOptions) error {
	return c.Collection.First(filter, model, opts)
}

func (c *Collection) Create(model Model, opts *options.InsertOneOptions) error {
	return c.Collection.Create(model, opts)
}

func (c *Collection) Update(model Model, opts *options.UpdateOptions) error {
	return c.Collection.Update(model, opts)
}

func (c *Collection) Save(model Model) error {
	return c.Collection.Save(model)
}

func (c *Collection) Delete(model Model) error {
	return c.Collection.Delete(model)
}
