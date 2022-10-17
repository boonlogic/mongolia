package odm

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Schema struct {
	Name       string
	Definition []byte
	Validator  func(any) error
	Hooks      *Hooks
}

func (s Schema) preValidate(any) *Model {
	return nil
}

func (s Schema) postValidate(any) *Model {
	return nil
}

func (s Schema) preSave(any) *Model {
	return nil
}

func (s Schema) postSave(any) *Model {
	return nil
}

func (s Schema) preCreate(any) *Model {
	return nil
}

func (s Schema) postCreate(any) *Model {
	return nil
}

func (s Schema) preUpdate(any) *Model {
	return nil
}

func (s Schema) postUpdate(any) *Model {
	return nil
}

func (s Schema) preRemove(any) *Model {
	return nil
}

func (s Schema) postRemove(any) *Model {
	return nil
}

// Mongo driver by default creates field names from the struct field name lowercase.
// It ignores the `json` struct tag (uses `bson`).

// REST      pkg/restapi    pkg/odm (validate/hooks) mongo-driver (access layer)
// -----------------------------------------------------------------------------------
// POST   -> create      -> createOne, createMany -> mongo.InsertOne, mongo.InsertMany
// GET    -> read        -> findOne, findMany     -> mongo.FindOne, mongo.FindMany
// PUT    -> update      -> updateOne, updateMany -> mongo.UpdateOne, mongo.UpdateMany
// DELETE -> delete      -> removeOne, removeMany -> mongo.DeleteOne, mongo.DeleteMany

func (s Schema) CreateOne(obj any) (*Document, error) {
	if err := s.Validator(obj); err != nil {
		return nil, err
	}

	var doc *Document

	fn := func(ctx context.Context) error {
		// Insert document.
		res, err := db.Collection(s.Name).InsertOne(ctx, obj)
		if err != nil {
			return err
		}

		id := res.InsertedID.(primitive.ObjectID)

		// Find inserted document.
		query := bson.M{"_id": id}
		if err := db.Collection(s.Name).FindOne(ctx, query).Decode(&doc); err != nil {
			return err
		}
		return nil
	}

	if err := transact(ctx(), fn); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) CreateMany(objs any) ([]Document, error) {
	arr := objs.([]any)
	for _, obj := range arr {
		if err := s.Validator(obj); err != nil {
			return nil, err
		}
	}

	var docs []Document

	fn := func(ctx context.Context) error {
		// Insert documents.
		res, err := db.Collection(s.Name).InsertMany(ctx, arr)
		if err != nil {
			return err
		}

		ids := make([]primitive.ObjectID, len(res.InsertedIDs))
		for i, v := range res.InsertedIDs {
			ids[i] = v.(primitive.ObjectID)
		}

		// Find inserted documents.
		filter := bson.M{"_id": bson.M{"$in": ids}}
		cur, err := db.Collection(s.Name).Find(ctx, filter)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		return nil
	}

	if err := transact(ctx(), fn); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s Schema) FindOne(any) (*Document, error) {
	filter := bson.M{}
	var doc *Document
	if err := db.Collection(s.Name).FindOne(ctx(), filter).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) FindMany(any) ([]Document, error) {
	filter := bson.M{}
	cur, err := db.Collection(s.Name).Find(ctx(), filter)
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
	err := db.Collection(s.Name).FindOneAndUpdate(ctx(), filter, update, opts).Decode(&doc); if err != nil {
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
		res, err := db.Collection(s.Name).UpdateMany(ctx, filter, update, uopts)
		if err != nil {
			return err
		}
		if res.MatchedCount < 1 {
			return nil
		}

		// Find the documents that were updated.
		fopts := options.Find().SetLimit(res.MatchedCount)
		cur, err := db.Collection(s.Name).Find(ctx, filter, fopts)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		return nil
	}

	if err := transact(ctx(), fn); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s Schema) RemoveOne(any) (*Document, error) {
	filter := bson.M{}
	var doc *Document
	if err := db.Collection(s.Name).FindOneAndDelete(ctx(), filter).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) RemoveMany(any) ([]Document, error) {
	var docs []Document

	fn := func(ctx context.Context) error {
		filter := bson.M{}

		// Find the documents that will be matched by the delete query.
		cur, err := db.Collection(s.Name).Find(ctx, filter)
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}
		nmatched := len(docs)

		// Delete them and ensure that the count is the same.
		res, err := db.Collection(s.Name).DeleteMany(ctx, filter)
		if err != nil {
			return err
		}
		if res.DeletedCount != int64(nmatched) {
			return errors.New("deleted count did not equal matched count")
		}
		return nil
	}

	if err := transact(ctx(), fn); err != nil {
		return nil, err
	}
	return docs, nil
}

// transact runs the passed-in function inside a Transaction.
func transact(ctx context.Context, fn func(context.Context) error) error {
	sess, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(ctx)

	txn, err := makeTransaction(sess, fn)
	if err != nil {
		return err
	}
	if err = mongo.WithSession(ctx, sess, txn); err != nil {
		if aerr := sess.AbortTransaction(ctx); aerr != nil {
			return aerr
		}
		return err
	}
	return nil
}

// makeTransaction decorates a func so that it runs as a Transaction under the given mongo.Session.
func makeTransaction(sess mongo.Session, fn func(context.Context) error) (func(mongo.SessionContext) error, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	opts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	txn := func(ctx mongo.SessionContext) error {
		if err := sess.StartTransaction(opts); err != nil {
			return err
		}
		if err := fn(ctx); err != nil {
			return err
		}
		if err := sess.CommitTransaction(ctx); err != nil {
			return err
		}
		return nil
	}
	return txn, nil
}
