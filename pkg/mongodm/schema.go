package mongodm

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Schema implements the Collection interface.
type Schema struct {
	validate   func(any) error   // function which validates Documents against this schema
	collection *mongo.Collection // handle to the underlying mongo collection
	hooks      *Hooks            // function pointers for hooks
}

func (s Schema) CreateOne(obj any) (*Document, error) {
	var doc Document

	// Convert the object into a document.
	buf, err := bson.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err := bson.Unmarshal(buf, &doc); err != nil {
		return nil, err
	}

	// Define a closure that inserts document and populates its id field.
	createOne := func(ctx context.Context, doc *Document) error {
		res, err := s.collection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}
		id := res.InsertedID.(primitive.ObjectID)
		for i := range *doc {
			if (*doc)[i].Key == "id" {
				(*doc)[i].Value = id
			}
		}
		return nil
	}

	if err := s.preValidate(&doc); err != nil {
		return nil, err
	}
	if err := s.validate(&doc); err != nil {
		return nil, err
	}
	if err := s.postValidate(&doc); err != nil {
		return nil, err
	}
	if err := s.preCreate(&doc); err != nil {
		return nil, err
	}
	if err := s.preSave(&doc); err != nil {
		return nil, err
	}
	if err := createOne(ctx(), &doc); err != nil {
		return nil, err
	}
	if err := s.postCreate(&doc); err != nil {
		return nil, err
	}
	if err := s.postSave(&doc); err != nil {
		return nil, err
	}
	out := Document(doc)
	return &out, nil
}

func (s Schema) CreateMany(objs any) ([]Document, error) {
	arr := objs.([]any)
	for _, obj := range arr {
		if err := s.validate(obj); err != nil {
			return nil, err
		}
	}

	var docs []Document

	fn := func(ctx context.Context) error {
		// Insert documents.
		res, err := s.collection.InsertMany(ctx, arr)
		if err != nil {
			return err
		}

		ids := make([]primitive.ObjectID, len(res.InsertedIDs))
		for i, v := range res.InsertedIDs {
			ids[i] = v.(primitive.ObjectID)
		}

		// Find inserted documents.
		filter := bson.M{"_id": bson.M{"$in": ids}}
		cur, err := s.collection.Find(ctx, filter)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		return nil
	}

	if err := fn(ctx()); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s Schema) FindOne(any) (*Document, error) {
	filter := bson.M{}
	var doc *Document
	if err := s.collection.FindOne(ctx(), filter).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) FindMany(any) ([]Document, error) {
	filter := bson.M{}
	cur, err := s.collection.Find(ctx(), filter)
	if err != nil {
		return nil, err
	}
	var docs []Document
	if err := cur.All(ctx(), &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s Schema) UpdateOne(any) (*Document, error) {
	filter := bson.M{}
	update := bson.M{
		"$unset": bson.M{"unset_me": 1},
	}
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)

	var doc *Document
	err := s.collection.FindOneAndUpdate(ctx(), filter, update, opts).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) UpdateMany(any) ([]Document, error) {
	var docs []Document

	fn := func(ctx context.Context) error {
		filter := bson.M{}

		// Update documents.
		update := bson.M{
			"$unset": bson.M{"unset_me": 1},
		}
		uopts := options.Update().SetUpsert(false)
		res, err := s.collection.UpdateMany(ctx, filter, update, uopts)
		if err != nil {
			return err
		}
		if res.MatchedCount < 1 {
			return nil
		}

		// Find the documents that were updated.
		fopts := options.Find().SetLimit(res.MatchedCount)
		cur, err := s.collection.Find(ctx, filter, fopts)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		return nil
	}

	if err := fn(ctx()); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s Schema) RemoveOne(any) (*Document, error) {
	filter := bson.M{}
	var doc *Document
	if err := s.collection.FindOneAndDelete(ctx(), filter).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) RemoveMany(any) ([]Document, error) {
	var docs []Document

	fn := func(ctx context.Context) error {
		filter := bson.M{}

		// Find the documents that will be matched by the delete query.
		cur, err := s.collection.Find(ctx, filter)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		nmatched := len(docs)

		// Delete them and ensure that the count is the same.
		res, err := s.collection.DeleteMany(ctx, filter)
		if err != nil {
			return err
		}
		if res.DeletedCount != int64(nmatched) {
			return errors.New("deleted count did not equal matched count")
		}
		return nil
	}

	if err := fn(ctx()); err != nil {
		return nil, err
	}
	return docs, nil
}

// preValidate is triggered before a document is validated against the schema.
func (s Schema) preValidate(doc *Document) error {
	if s.hooks.PreValidate == nil {
		return nil
	}
	return s.hooks.PreValidate(doc)
}

// postValidate is triggered after a document is validate against the schema.
func (s Schema) postValidate(doc *Document) error {
	if s.hooks.PostValidate == nil {
		return nil
	}
	return s.hooks.PostValidate(doc)
}

// preCreate is triggered after postValidate and before inserting a document.
func (s Schema) preCreate(doc *Document) error {
	if s.hooks.PreCreate == nil {
		return nil
	}
	return s.hooks.PreCreate(doc)
}

// preUpdate is triggered after postValidate and before updating a document.
func (s Schema) preUpdate(doc *Document) error {
	if s.hooks.PreUpdate == nil {
		return nil
	}
	return s.hooks.PreUpdate(doc)
}

// preSave is triggered after preCreate/preUpdate and before inserting or updating a document.
func (s Schema) preSave(doc *Document) error {
	if s.hooks.PreSave == nil {
		return nil
	}
	return s.hooks.PreSave(doc)
}

// preRemove is triggered before removing a document.
func (s Schema) preRemove(doc *Document) error {
	if s.hooks.PreRemove == nil {
		return nil
	}
	return s.hooks.PreRemove(doc)
}

// postCreate is triggered after inserting a document.
func (s Schema) postCreate(doc *Document) error {
	if s.hooks.PostCreate == nil {
		return nil
	}
	return s.hooks.PostCreate(doc)
}

// postUpdate is triggered after updating a document.
func (s Schema) postUpdate(doc *Document) error {
	if s.hooks.PostUpdate == nil {
		return nil
	}
	return s.hooks.PostUpdate(doc)
}

// postSave is triggered after postCreate and postUpdate, after inserting or updating a document.
func (s Schema) postSave(doc *Document) error {
	if s.hooks.PostSave == nil {
		return nil
	}
	return s.hooks.PostSave(doc)
}

// postRemove is triggered after removing a document.
func (s Schema) postRemove(doc *Document) error {
	if s.hooks.PostRemove == nil {
		return nil
	}
	return s.hooks.PostRemove(doc)
}