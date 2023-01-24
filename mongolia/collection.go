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
	name    string
	coll    *mongo.Collection
	timeout time.Duration
}

// Validate read is not called here
func (c *Collection) CountDocuments(filter any, opts *options.CountOptions) (int64, *Error) {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)

	//Get document counts
	count, err := c.coll.CountDocuments(ctx, filter, opts)
	if err != nil {
		return int64(0), NewError(404, err)
	}
	return int64(count), nil
}

func (c *Collection) CreateIndexes(indexes interface{}) *Error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if err := PopulateIndexes(ctx, c.coll, indexes); err != nil {
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

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.InsertOne(ctx, model, opts)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return NewError(400, err)
		}
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

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if err := c.coll.FindOne(ctx, bson.D{{"_id", idp}}).Decode(model); err != nil {
		if err == mongo.ErrNoDocuments {
			return NewError(404, err)
		}
		return NewError(500, err)
	}
	if err := model.ValidateRead(); err != nil {
		return NewError(406, err)
	}
	return nil
}

func (c *Collection) FindOne(filter any, model Model, opts *options.FindOneOptions) *Error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if err := c.coll.FindOne(ctx, filter, opts).Decode(model); err != nil {
		if err == mongo.ErrNoDocuments {
			return NewError(404, err)
		}
		return NewError(500, err)
	}
	if err := model.ValidateRead(); err != nil {
		return NewError(406, err)
	}
	return nil
}

func (c *Collection) FindOneAndUpdate(filter any, update any, model Model, opts *options.FindOneAndUpdateOptions) *Error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if err := c.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(model); err != nil {
		if err == mongo.ErrNoDocuments {
			return NewError(404, err)
		}
		return NewError(500, err)
	}
	if err := model.ValidateRead(); err != nil {
		return NewError(406, err)
	}
	return nil
}

// Validate read is not called here
func (c *Collection) Find(filter any, results interface{}, opts *options.FindOptions) *Error {

	//Verify Type
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr {
		return NewErrorString(406, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}
	if resultsValue.Elem().Kind() != reflect.Slice {
		return NewErrorString(406, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)

	cursor, err := c.coll.Find(ctx, filter, opts)
	if err != nil {
		return NewError(404, err)
	}

	err = cursor.All(ctx, results)
	if err != nil {
		return NewError(404, err)
	}

	return nil
}

// Validate read is not called here
func (c *Collection) FindWithResults(filter any, results interface{}, opts *options.FindOptions) (*FindResult, *Error) {

	//Verify Type
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr {
		return nil, NewErrorString(406, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}
	if resultsValue.Elem().Kind() != reflect.Slice {
		return nil, NewErrorString(406, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)

	//Get document counts
	var findResult FindResult
	countopts := options.Count().SetMaxTime(2 * time.Second)
	filtered, err := c.coll.CountDocuments(ctx, filter, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Filtered = filtered

	collection, err := c.coll.CountDocuments(ctx, bson.D{}, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Collection = collection

	lookup, err := c.coll.CountDocuments(ctx, filter, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Collection = lookup

	cursor, err := c.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, NewError(404, err)
	}

	err = cursor.All(ctx, results)
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
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	results, err := c.coll.Distinct(ctx, field, filter)
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

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.UpdateOne(ctx, bson.M{"_id": model.GetID()}, bson.M{"$set": model}, opts)
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

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.UpdateOne(ctx, filter, bson.M{"$set": model}, opts)
	if err != nil {
		return NewError(400, err)
	}

	return afterUpdateHooks(res, model)
}

func (c *Collection) UpdateSet(filter any, update any, opts *options.UpdateOptions) *Error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	_, err := c.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return NewError(400, err)
	}

	return nil
}

func (c *Collection) Delete(model Model) *Error {
	if err := beforeDeleteHooks(model); err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.DeleteOne(ctx, bson.M{"_id": model.GetID()})
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

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.DeleteOne(ctx, bson.M{"_id": idp})
	if err != nil {
		return NewError(400, err)
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) DeleteOne(filter any, model Model) *Error {
	if err := beforeDeleteHooks(model); err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.DeleteOne(ctx, filter)
	if err != nil {
		return NewError(400, err)
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) DeleteMany(filter any) *Error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	_, err := c.coll.DeleteMany(ctx, filter)
	if err != nil {
		return NewError(400, err)
	}

	return nil
}

func (c *Collection) Drop() *Error {
	if err := c.coll.Drop(context.Background()); err != nil {
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

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	cursor, err := c.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return NewError(404, err)
	}
	err = cursor.All(ctx, results)
	if err != nil {
		return NewError(404, err)
	}

	return nil
}
