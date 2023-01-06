package mongolia

import (
	"context"
	"fmt"
	"reflect"
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

func (c *Collection) CreateIndexes(indexes interface{}) *Error {
	if err := PopulateIndexes(c.ctx, c.coll, indexes); err != nil {
		return NewError(503, err)
	}
	return nil
}

func (c *Collection) Save(model Model) *Error {
	if err := model.ValidateRead(); err != nil {
		return NewError(406, err)
	}

	if err := beforeSaveHooks(model); err != nil {
		return err
	}

	return afterSaveHooks(model)
}

func (c *Collection) Create(model Model, opts *options.InsertOneOptions) *Error {
	if err := model.ValidateCreate(); err != nil {
		return NewError(406, err)
	}

	if err := beforeCreateHooks(model); err != nil {
		return err
	}

	res, err := c.coll.InsertOne(c.ctx, model, opts)
	if err != nil {
		return NewError(400, err)
	}

	// Set new id
	model.SetID(res.InsertedID)

	return afterCreateHooks(model)
}

func (c *Collection) FindByID(id any, model Model) *Error {
	idp, err := model.PrepareID(id)
	if err != nil {
		return NewError(406, err)
	}

	if err := c.coll.FindOne(c.ctx, bson.D{{"_id", idp}}).Decode(model); err != nil {
		return NewError(404, err)
	}
	if err := model.ValidateRead(); err != nil {
		return NewError(406, err)
	}
	return nil
}

func (c *Collection) FindOne(filter any, model Model, opts *options.FindOneOptions) *Error {
	if err := c.coll.FindOne(c.ctx, filter, opts).Decode(model); err != nil {
		return NewError(404, err)
	}
	if err := model.ValidateRead(); err != nil {
		return NewError(406, err)
	}
	return nil
}

// Validate read is not called here
func (c *Collection) Find(filter any, results interface{}, opts *options.FindOptions) (*FindResult, *Error) {

	//Verify Type
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr {
		return nil, NewErrorString(406, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}
	if resultsValue.Elem().Kind() != reflect.Slice {
		return nil, NewErrorString(406, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}

	//Get document counts
	var findResult FindResult
	countopts := options.Count().SetMaxTime(2 * time.Second)
	filtered, err := c.coll.CountDocuments(c.ctx, filter, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Filtered = filtered

	collection, err := c.coll.CountDocuments(c.ctx, bson.D{}, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Collection = collection

	lookup, err := c.coll.CountDocuments(c.ctx, filter, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Collection = lookup

	cursor, err := c.coll.Find(c.ctx, filter, opts)
	if err != nil {
		return nil, NewError(404, err)
	}

	err = cursor.All(c.ctx, results)
	if err != nil {
		return nil, NewError(404, err)
	}

	//Get Length By Reflecting
	modelType := reflect.ValueOf(results).Elem()
	findResult.Limit = int64(modelType.Len())

	// for _, temp := range modelType {
	// 	if err := temp.ValidateRead(); err != nil {
	// 		return nil, err
	// 	}
	// }

	return &findResult, nil
}

// Retrieve unique values of field in collection, filter to limit scope
func (c *Collection) Distinct(filter any, field string) (interface{}, *Error) {
	results, err := c.coll.Distinct(c.ctx, field, filter)
	if err != nil {
		return nil, NewError(404, err)
	}
	return results, nil
}

func (c *Collection) Update(model Model, opts *options.UpdateOptions) *Error {
	if err := model.ValidateUpdate(); err != nil {
		return NewError(406, err)
	}

	if err := beforeUpdateHooks(model); err != nil {
		return err
	}

	res, err := c.coll.UpdateOne(c.ctx, bson.M{"_id": model.GetID()}, bson.M{"$set": model}, opts)
	if err != nil {
		return NewError(400, err)
	}

	return afterUpdateHooks(res, model)
}

func (c *Collection) UpdateOne(filter any, model Model, opts *options.UpdateOptions) *Error {
	if err := model.ValidateUpdate(); err != nil {
		return NewError(406, err)
	}

	if err := beforeUpdateHooks(model); err != nil {
		return err
	}

	res, err := c.coll.UpdateOne(c.ctx, filter, bson.M{"$set": model}, opts)
	if err != nil {
		return NewError(400, err)
	}

	return afterUpdateHooks(res, model)
}

func (c *Collection) Delete(model Model) *Error {
	if err := beforeDeleteHooks(model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(c.ctx, bson.M{"_id": model.GetID()})
	if err != nil {
		return NewError(400, err)
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) DeleteByID(id any, model Model) *Error {
	idp, err := model.PrepareID(id)
	if err != nil {
		return NewError(406, err)
	}

	if err := beforeDeleteHooks(model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(c.ctx, bson.M{"_id": idp})
	if err != nil {
		return NewError(400, err)
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) DeleteOne(filter any, model Model) *Error {
	if err := beforeDeleteHooks(model); err != nil {
		return err
	}
	res, err := c.coll.DeleteOne(c.ctx, filter)
	if err != nil {
		return NewError(400, err)
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) Drop() *Error {
	if err := c.coll.Drop(c.ctx); err != nil {
		return NewError(404, err)
	}
	return nil
}

// Run a mongo.Pipeline on a specific colleciton
// Validate read is not called here
func (c *Collection) Aggregate(results interface{}, pipeline any) *Error {

	//Verify Type
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr {
		return NewErrorString(400, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}
	if resultsValue.Elem().Kind() != reflect.Slice {
		return NewErrorString(400, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}

	cursor, err := c.coll.Aggregate(c.ctx, pipeline)
	if err != nil {
		return NewError(404, err)
	}
	err = cursor.All(c.ctx, results)
	if err != nil {
		return NewError(404, err)
	}

	return nil
}
