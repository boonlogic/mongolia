package mongodm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema struct {
	name     string
	validate func(Attributes) error
	hooks    *Hooks
}

func AddSchema(name string, spec *Spec, hooks *Hooks) error {
	if _, ok := odm.schemas[name]; ok {
		return errors.New(fmt.Sprintf("a schema named '%s' already exists", name))
	}

	if err := spec.Validate(); err != nil {
		return err
	}
	vfunc, err := spec.ValidatorFunc()
	if err != nil {
		return err
	}

	schema := &Schema{
		name:     name,
		validate: vfunc,
		hooks:    hooks,
	}
	odm.schemas[name] = schema

	coll := odm.db.Collection(name)
	// todo: convert spec into map[string]any
	// todo: parse indexes from spec x-attrs
	// todo: initialize collection indexes
	odm.colls[name] = coll

	return nil
}

func (s *Schema) preValidate(doc *Document) error {
	if s.hooks.PreValidate == nil {
		return nil
	}
	return s.hooks.PreValidate(doc)
}

func (s *Schema) postValidate(doc *Document) error {
	if s.hooks.PostValidate == nil {
		return nil
	}
	return s.hooks.PostValidate(doc)
}

func (s *Schema) preCreate(doc *Document) error {
	if s.hooks.PreCreate == nil {
		return nil
	}
	return s.hooks.PreCreate(doc)
}

func (s *Schema) preUpdate(doc *Document) error {
	if s.hooks.PreUpdate == nil {
		return nil
	}
	return s.hooks.PreUpdate(doc)
}

func (s *Schema) preSave(doc *Document) error {
	if s.hooks.PreSave == nil {
		return nil
	}
	return s.hooks.PreSave(doc)
}

func (s *Schema) preRemove(doc *Document) error {
	if s.hooks.PreRemove == nil {
		return nil
	}
	return s.hooks.PreRemove(doc)
}

func (s *Schema) postCreate(doc *Document) error {
	if s.hooks.PostCreate == nil {
		return nil
	}
	return s.hooks.PostCreate(doc)
}

func (s *Schema) postUpdate(doc *Document) error {
	if s.hooks.PostUpdate == nil {
		return nil
	}
	return s.hooks.PostUpdate(doc)
}

func (s *Schema) postSave(doc *Document) error {
	if s.hooks.PostSave == nil {
		return nil
	}
	return s.hooks.PostSave(doc)
}

func (s *Schema) postRemove(doc *Document) error {
	if s.hooks.PostRemove == nil {
		return nil
	}
	return s.hooks.PostRemove(doc)
}

func (s *Schema) runWithHooks(ctx context.Context, operator func(ctx context.Context) error, doc *Document) error {
	if err := s.preValidate(doc); err != nil {
		return err
	}
	if err := s.validate(*doc.attrs); err != nil {
		return err
	}
	if err := s.postValidate(doc); err != nil {
		return err
	}
	if err := s.preCreate(doc); err != nil {
		return err
	}
	if err := s.preSave(doc); err != nil {
		return err
	}
	if err := operator(ctx); err != nil {
		return err
	}
	if err := s.postCreate(doc); err != nil {
		return err
	}
	if err := s.postSave(doc); err != nil {
		return err
	}
	return nil
}

func (s *Schema) createOne(attrs *Attributes) (*Document, error) {
	doc := &Document{
		id:    primitive.ObjectID{},
		attrs: attrs,
	}
	fn := func(ctx context.Context) error {
		if err := createOne(ctx, s.name, doc); err != nil {
			return err
		}
		return nil
	}
	s.runWithHooks(ctx(), fn, doc)
	return doc, nil
}

//func (s *Schema) createMany(attrs []Attributes) ([]Document, error) {
//	for _, v := range attrs {
//		if err := s.validate(v); err != nil {
//			return nil, err
//		}
//	}
//
//	var docs []Document
//
//	// todo: use for loop and CreateOne
//	fn := func(ctx context.Context) error {
//		arr := make([]any, 0)
//		for _, v := range attrs {
//			arr = append(arr, any(v))
//		}
//
//		res, err := odm.db.Collection(s.name).InsertMany(ctx, arr)
//		if err != nil {
//			return err
//		}
//
//		ids := make([]primitive.ObjectID, len(res.InsertedIDs))
//		for i, v := range res.InsertedIDs {
//			ids[i] = v.(primitive.ObjectID)
//		}
//
//		filter := bson.M{"_id": bson.M{"$in": ids}}
//		cur, err := odm.db.Collection(s.name).Find(ctx, filter)
//		if err != nil {
//			return err
//		}
//		if err := cur.All(ctx, &res); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	if err := fn(ctx()); err != nil {
//		return nil, err
//	}
//	return docs, nil
//}
//
//func FindOne(query Query) (*Document, error) {
//	filter := bson.M{}
//	var doc *Document
//	if err := odm.db.Collection(s.name).FindOne(ctx(), filter).Decode(&doc); err != nil {
//		return nil, err
//	}
//	return doc, nil
//}
//
//func FindMany(query Query) ([]Document, error) {
//	filter := bson.M{}
//	cur, err := odm.db.Collection(s.name).Find(ctx(), filter)
//	if err != nil {
//		return nil, err
//	}
//	var docs []Document
//	if err := cur.All(ctx(), &docs); err != nil {
//		return nil, err
//	}
//	return docs, nil
//}
//
//func UpdateOne(query Query, attr Attributes) (*Document, error) {
//	filter := bson.M{}
//	update := bson.M{
//		"$unset": bson.M{"unset_me": 1},
//	}
//	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
//
//	var doc *Document
//	err := odm.db.Collection(s.name).FindOneAndUpdate(ctx(), filter, update, opts).Decode(&doc)
//	if err != nil {
//		return nil, err
//	}
//	return doc, nil
//}
//
//func UpdateMany(query Query, attr Attributes) ([]Document, error) {
//	var docs []Document
//
//	fn := func(ctx context.Context) error {
//		filter := bson.M{}
//
//		// Update documents.
//		update := bson.M{
//			"$unset": bson.M{"unset_me": 1},
//		}
//		uopts := options.Update().SetUpsert(false)
//		res, err := odm.db.Collection(s.name).UpdateMany(ctx, filter, update, uopts)
//		if err != nil {
//			return err
//		}
//		if res.MatchedCount < 1 {
//			return nil
//		}
//
//		// Find the documents that were updated.
//		fopts := options.Find().SetLimit(res.MatchedCount)
//		cur, err := odm.db.Collection(s.name).Find(ctx, filter, fopts)
//		if err != nil {
//			return err
//		}
//		if err := cur.All(ctx, &docs); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	if err := fn(ctx()); err != nil {
//		return nil, err
//	}
//	return docs, nil
//}
//
//func RemoveOne(query Query) (*Document, error) {
//	filter := bson.M{}
//	var doc *Document
//	if err := odm.db.Collection(s.name).FindOneAndDelete(ctx(), filter).Decode(&doc); err != nil {
//		return nil, err
//	}
//	return doc, nil
//}
//
//func RemoveMany(query Query) ([]Document, error) {
//	var docs []Document
//
//	fn := func(ctx context.Context) error {
//		filter := bson.M{}
//
//		// Find the documents that will be matched by the delete query.
//		cur, err := odm.db.Collection(s.name).Find(ctx, filter)
//		if err != nil {
//			return err
//		}
//		if err := cur.All(ctx, &docs); err != nil {
//			return err
//		}
//		nmatched := len(docs)
//
//		// Delete them and ensure that the count is the same.
//		res, err := odm.db.Collection(s.name).DeleteMany(ctx, filter)
//		if err != nil {
//			return err
//		}
//		if res.DeletedCount != int64(nmatched) {
//			return errors.New("deleted count did not equal matched count")
//		}
//		return nil
//	}
//
//	if err := fn(ctx()); err != nil {
//		return nil, err
//	}
//	return docs, nil
//}
