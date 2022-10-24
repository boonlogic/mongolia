package mongodm

import (
	options2 "gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/options"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Singleton containing all collections.
var collections = make(map[string]*CollectionNew)

type CollectionNew struct {
	name string
	coll *mongo.Collection
}

func GetCollection(name string) (*CollectionNew, bool) {
	c, ok := collections[name]
	if !ok {
		return nil, false
	}
	return &c, true
}

func (c CollectionNew) CreateOne(attr types.Attributes) (*types.Document, error) {
	createOne := func(ctx context.Context, attr *types.Attributes) error {
		res, err := c.collection.InsertOne(ctx, attr.M())
		if err != nil {
			return err
		}
		(*attr)["id"] = res.InsertedID.(primitive.ObjectID)
		return nil
	}

	if err := c.preValidate(&attr); err != nil {
		return nil, err
	}
	if err := c.validate(&attr); err != nil {
		return nil, err
	}
	if err := c.postValidate(&doc); err != nil {
		return nil, err
	}
	if err := c.preCreate(&doc); err != nil {
		return nil, err
	}
	if err := c.preSave(&doc); err != nil {
		return nil, err
	}
	if err := createOne(options2.ctx(), &doc); err != nil {
		return nil, err
	}
	if err := c.postCreate(&doc); err != nil {
		return nil, err
	}
	if err := c.postSave(&doc); err != nil {
		return nil, err
	}
	out := types.Document(doc)
	return &out, nil
}

func (c CollectionNew) CreateMany(attrs []types.Attributes) ([]types.Document, error) {
	for _, v := range attrs {
		if err := c.validate(v); err != nil {
			return nil, err
		}
	}

	var docs []types.Document

	// todo: use for loop and CreateOne
	fn := func(ctx context.Context) error {
		arr := make([]any, 0)
		for _, v := range attrs {
			arr = append(arr, any(v))
		}

		res, err := c.collection.InsertMany(ctx, arr)
		if err != nil {
			return err
		}

		ids := make([]primitive.ObjectID, len(res.InsertedIDs))
		for i, v := range res.InsertedIDs {
			ids[i] = v.(primitive.ObjectID)
		}

		filter := bson.M{"_id": bson.M{"$in": ids}}
		cur, err := c.collection.Find(ctx, filter)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &res); err != nil {
			return err
		}
		return nil
	}

	if err := fn(options2.ctx()); err != nil {
		return nil, err
	}
	return docs, nil
}

func (c CollectionNew) FindOne(query Query) (*types.Document, error) {
	filter := bson.M{}
	var doc *types.Document
	if err := c.collection.FindOne(options2.ctx(), filter).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (c CollectionNew) FindMany(query Query) ([]types.Document, error) {
	filter := bson.M{}
	cur, err := c.collection.Find(options2.ctx(), filter)
	if err != nil {
		return nil, err
	}
	var docs []types.Document
	if err := cur.All(options2.ctx(), &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (c CollectionNew) UpdateOne(query Query, attr types.Attributes) (*types.Document, error) {
	filter := bson.M{}
	update := bson.M{
		"$unset": bson.M{"unset_me": 1},
	}
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)

	var doc *types.Document
	err := c.collection.FindOneAndUpdate(options2.ctx(), filter, update, opts).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (c CollectionNew) UpdateMany(query Query, attr types.Attributes) ([]types.Document, error) {
	var docs []types.Document

	fn := func(ctx context.Context) error {
		filter := bson.M{}

		// Update documents.
		update := bson.M{
			"$unset": bson.M{"unset_me": 1},
		}
		uopts := options.Update().SetUpsert(false)
		res, err := c.collection.UpdateMany(ctx, filter, update, uopts)
		if err != nil {
			return err
		}
		if res.MatchedCount < 1 {
			return nil
		}

		// Find the documents that were updated.
		fopts := options.Find().SetLimit(res.MatchedCount)
		cur, err := c.collection.Find(ctx, filter, fopts)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		return nil
	}

	if err := fn(options2.ctx()); err != nil {
		return nil, err
	}
	return docs, nil
}

func (c CollectionNew) RemoveOne(query Query) (*types.Document, error) {
	filter := bson.M{}
	var doc *types.Document
	if err := c.collection.FindOneAndDelete(options2.ctx(), filter).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (c CollectionNew) RemoveMany(query Query) ([]types.Document, error) {
	var docs []types.Document

	fn := func(ctx context.Context) error {
		filter := bson.M{}

		// Find the documents that will be matched by the delete query.
		cur, err := c.collection.Find(ctx, filter)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		nmatched := len(docs)

		// Delete them and ensure that the count is the same.
		res, err := c.collection.DeleteMany(ctx, filter)
		if err != nil {
			return err
		}
		if res.DeletedCount != int64(nmatched) {
			return errors.New("deleted count did not equal matched count")
		}
		return nil
	}

	if err := fn(options2.ctx()); err != nil {
		return nil, err
	}
	return docs, nil
}

// preValidate is triggered before a document is validated against the schema.
func (c CollectionNew) preValidate(doc *types.Document) error {
	if c.hooks.PreValidate == nil {
		return nil
	}
	return c.hooks.PreValidate(doc)
}

// postValidate is triggered after a document is validate against the schema.
func (c CollectionNew) postValidate(doc *types.Document) error {
	if c.hooks.PostValidate == nil {
		return nil
	}
	return c.hooks.PostValidate(doc)
}

// preCreate is triggered after postValidate and before inserting a document.
func (c CollectionNew) preCreate(doc *types.Document) error {
	if c.hooks.PreCreate == nil {
		return nil
	}
	return c.hooks.PreCreate(doc)
}

// preUpdate is triggered after postValidate and before updating a document.
func (c CollectionNew) preUpdate(doc *types.Document) error {
	if c.hooks.PreUpdate == nil {
		return nil
	}
	return c.hooks.PreUpdate(doc)
}

// preSave is triggered after preCreate/preUpdate and before inserting or updating a document.
func (c CollectionNew) preSave(doc *types.Document) error {
	if c.hooks.PreSave == nil {
		return nil
	}
	return c.hooks.PreSave(doc)
}

// preRemove is triggered before removing a document.
func (c CollectionNew) preRemove(doc *types.Document) error {
	if c.hooks.PreRemove == nil {
		return nil
	}
	return c.hooks.PreRemove(doc)
}

// postCreate is triggered after inserting a document.
func (c CollectionNew) postCreate(doc *types.Document) error {
	if c.hooks.PostCreate == nil {
		return nil
	}
	return c.hooks.PostCreate(doc)
}

// postUpdate is triggered after updating a document.
func (c CollectionNew) postUpdate(doc *types.Document) error {
	if c.hooks.PostUpdate == nil {
		return nil
	}
	return c.hooks.PostUpdate(doc)
}

// postSave is triggered after postCreate and postUpdate, after inserting or updating a document.
func (c CollectionNew) postSave(doc *types.Document) error {
	if c.hooks.PostSave == nil {
		return nil
	}
	return c.hooks.PostSave(doc)
}

// postRemove is triggered after removing a document.
func (c CollectionNew) postRemove(doc *types.Document) error {
	if c.hooks.PostRemove == nil {
		return nil
	}
	return c.hooks.PostRemove(doc)
}
