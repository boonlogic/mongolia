package mongolia

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return NewError(400, err)
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
	return afterReadHooks(model)
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
	return afterReadHooks(model)
}

func (c *Collection) FindOneAndUpdate(filter any, update any, model Model, opts *options.FindOneAndUpdateOptions) *Error {
	//Handle map based updates
	switch update.(type) {
	case map[string]any:
		bson_update := CastMapToDB(update.(map[string]any))
		return c.FindOneAndUpdate(filter, bson_update, model, opts)
	}

	if err := model.ValidateUpdate(update); err != nil {
		return NewError(406, err)
	}

	if err := beforeUpdateHooks(update, model); err != nil {
		return err
	}

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
	return afterReadHooks(model)
}

// Validate/hooks read is not called here
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

	//Get Underlying struct
	pointerSliceModel := resultsValue.Elem()
	sliceLength := pointerSliceModel.Len()

	//Attempt to cast each element of the slice back to model so we can call the validate function
	for i := 0; i < sliceLength; i++ {
		modelVar, ok := pointerSliceModel.Index(i).Interface().(Model)
		if ok {
			//If this pointer can be cast to model, we can call hooks
			if err := modelVar.ValidateRead(); err != nil {
				return NewError(406, err)
			}
			if merr := afterReadHooks(modelVar); merr != nil {
				return merr
			}
		}
	}

	return nil
}

// Validate/hooks read is not called here
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

	//cursor find
	cursor, err := c.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, NewError(404, err)
	}

	//Parse results
	cerr := cursor.All(ctx, results)
	if cerr != nil {
		return nil, NewError(404, cerr)
	}

	//Get Underlying struct
	pointerSliceModel := resultsValue.Elem()
	sliceLength := pointerSliceModel.Len()

	//Get document counts
	var findResult FindResult
	findResult.Filtered = int64(sliceLength)

	countopts := options.Count().SetMaxTime(2 * time.Second).SetHint("_id_")
	collection, err := c.coll.CountDocuments(ctx, bson.D{}, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Collection = collection

	findResult.Skip = 0
	if opts != nil && opts.Skip != nil {
		findResult.Skip = *opts.Skip
	}

	findResult.Limit = collection
	if opts != nil && opts.Limit != nil {
		findResult.Limit = *opts.Limit
	}

	//Attempt to cast each element of the slice back to model so we can call the validate function
	for i := 0; i < sliceLength; i++ {
		modelVar, ok := pointerSliceModel.Index(i).Interface().(Model)
		if ok {
			//If this pointer can be cast to model, we can call hooks
			if err := modelVar.ValidateRead(); err != nil {
				return nil, NewError(406, err)
			}
			if merr := afterReadHooks(modelVar); merr != nil {
				return nil, merr
			}
		}
	}

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

// Partial update
func (c *Collection) Update(filter any, update any, model Model, opts *options.UpdateOptions) *Error {
	//Handle map based updates
	switch update.(type) {
	case map[string]any:
		bson_update := CastMapToDB(update.(map[string]any))
		return c.Update(filter, bson_update, model, opts)
	}

	if err := model.ValidateUpdate(update); err != nil {
		return NewError(406, err)
	}

	if err := beforeUpdateHooks(update, model); err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return NewError(400, err)
	}

	if res.MatchedCount == 0 && res.UpsertedCount == 0 {
		return NewErrorString(404, "mongo: no documents in result")
	}

	return afterUpdateHooks(res, update, model)
}

// Update entire model
func (c *Collection) UpdateModel(model Model, opts *options.UpdateOptions) *Error {
	filter := bson.D{{"_id", model.GetID().(primitive.ObjectID)}}

	return c.UpdateModelQuery(filter, model, opts)
}

// Update entire model from filter
func (c *Collection) UpdateModelQuery(filter any, model Model, opts *options.UpdateOptions) *Error {
	if err := model.ValidateUpdateModel(); err != nil {
		return NewError(406, err)
	}

	if err := beforeUpdateHooks(nil, model); err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	res, err := c.coll.UpdateOne(ctx, filter, bson.M{"$set": model}, opts)
	if err != nil {
		return NewError(400, err)
	}

	if res.MatchedCount == 0 && res.UpsertedCount == 0 {
		return NewErrorString(404, "mongo: no documents in result")
	}

	return afterUpdateHooks(res, nil, model)

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

	//throw error if no documents found
	if res.DeletedCount == 0 {
		return NewErrorString(404, "mongo: no documents in result")
	}

	return afterDeleteHooks(res, model)
}

func (c *Collection) DeleteByID(id any, model Model) *Error {
	idp, err := model.PrepareID(id)
	if err != nil {
		return NewError(400, err)
	}

	return c.DeleteOne(bson.M{"_id": idp}, model)
}

// delete by query, populate model with result
func (c *Collection) DeleteOne(filter any, model Model) *Error {
	if err := beforeDeleteHooks(model); err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	singleResult := c.coll.FindOneAndDelete(ctx, filter)
	if singleResult.Err() != nil {
		//throw error if no documents found
		return NewErrorString(404, "mongo: no documents in result")
	}

	//populate deleted document
	derr := singleResult.Decode(model)
	if derr != nil {
		return NewError(500, derr)
	}

	//If we got here, the deleted count is one
	res := mongo.DeleteResult{
		DeletedCount: 1,
	}

	return afterDeleteHooks(&res, model)
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

	//Get Underlying struct
	pointerSliceModel := resultsValue.Elem()
	sliceLength := pointerSliceModel.Len()

	//Attempt to cast each element of the slice back to model so we can call the validate function
	for i := 0; i < sliceLength; i++ {
		modelVar, ok := pointerSliceModel.Index(i).Interface().(Model)
		if ok {
			//If this pointer can be cast to model, we can call hooks
			if err := modelVar.ValidateRead(); err != nil {
				return NewError(406, err)
			}
			if merr := afterReadHooks(modelVar); merr != nil {
				return merr
			}
		}
	}

	return nil
}

// Run a mongo.Pipeline on a specific colleciton with find results
func (c *Collection) AggregateWithResults(results interface{}, pipeline any, skip *int64, limit *int64) (*FindResult, *Error) {

	//Verify Type
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr {
		return nil, NewErrorString(400, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}
	if resultsValue.Elem().Kind() != reflect.Slice {
		return nil, NewErrorString(400, fmt.Sprintf("Expecting results to be a pointer to slice, instead got %v", resultsValue.Kind()))
	}

	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	cursor, err := c.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, NewError(404, err)
	}

	//Parse all results
	err = cursor.All(ctx, results)
	if err != nil {
		return nil, NewError(404, err)
	}

	//Get Underlying struct
	pointerSliceModel := resultsValue.Elem()
	sliceLength := pointerSliceModel.Len()

	//Attempt to cast each element of the slice back to model so we can call the validate function
	for i := 0; i < sliceLength; i++ {
		modelVar, ok := pointerSliceModel.Index(i).Interface().(Model)
		if ok {
			//If this pointer can be cast to model, we can call hooks
			if err := modelVar.ValidateRead(); err != nil {
				return nil, NewError(406, err)
			}
			if merr := afterReadHooks(modelVar); merr != nil {
				return nil, merr
			}
		}
	}

	//Get document counts
	var findResult FindResult
	findResult.Filtered = int64(sliceLength)

	countopts := options.Count().SetMaxTime(2 * time.Second).SetHint("_id_")
	fullcount, err := c.coll.CountDocuments(ctx, bson.D{}, countopts)
	if err != nil {
		return nil, NewError(404, err)
	}
	findResult.Collection = fullcount

	findResult.Skip = 0
	if skip != nil {
		findResult.Skip = *skip
	}

	findResult.Limit = fullcount
	if limit != nil {
		findResult.Limit = *limit
	}

	return &findResult, nil
}
