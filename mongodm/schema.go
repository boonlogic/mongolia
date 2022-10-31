package mongodm

import (
	"context"
	"encoding/json"
)

type Schema struct {
	definition map[string]any
	validator  func(map[string]any) error
	hooks      *Hooks
}

// NewSchema creates a new Schema based on the given Spec and Hooks.
func NewSchema(spec *Spec, hooks *Hooks) (*Schema, error) {
	return newSchema(spec, hooks)
}

func newSchema(spec *Spec, hooks *Hooks) (*Schema, error) {
	vfunc, err := spec.GetValidator()
	if err != nil {
		return nil, err
	}

	var def map[string]any
	if err = json.Unmarshal(spec.definition, &def); err != nil {
		return nil, err
	}

	s := &Schema{
		definition: def,
		validator:  vfunc,
		hooks:      hooks,
	}
	return s, nil
}

func (s *Schema) preValidate(doc *Model) error {
	if s.hooks.PreValidate == nil {
		return nil
	}
	return s.hooks.PreValidate(doc)
}

func (s *Schema) postValidate(doc *Model) error {
	if s.hooks.PostValidate == nil {
		return nil
	}
	return s.hooks.PostValidate(doc)
}

func (s *Schema) preCreate(doc *Model) error {
	if s.hooks.PreCreate == nil {
		return nil
	}
	return s.hooks.PreCreate(doc)
}

func (s *Schema) preUpdate(doc *Model) error {
	if s.hooks.PreUpdate == nil {
		return nil
	}
	return s.hooks.PreUpdate(doc)
}

func (s *Schema) preSave(doc *Model) error {
	if s.hooks.PreSave == nil {
		return nil
	}
	return s.hooks.PreSave(doc)
}

func (s *Schema) preRemove(doc *Model) error {
	if s.hooks.PreRemove == nil {
		return nil
	}
	return s.hooks.PreRemove(doc)
}

func (s *Schema) postCreate(doc *Model) error {
	if s.hooks.PostCreate == nil {
		return nil
	}
	return s.hooks.PostCreate(doc)
}

func (s *Schema) postUpdate(doc *Model) error {
	if s.hooks.PostUpdate == nil {
		return nil
	}
	return s.hooks.PostUpdate(doc)
}

func (s *Schema) postSave(doc *Model) error {
	if s.hooks.PostSave == nil {
		return nil
	}
	return s.hooks.PostSave(doc)
}

func (s *Schema) postRemove(doc *Model) error {
	if s.hooks.PostRemove == nil {
		return nil
	}
	return s.hooks.PostRemove(doc)
}

func (s *Schema) runWithHooks(ctx context.Context, operator func(ctx context.Context) error, model *Model) error {
	if err := s.preValidate(model); err != nil {
		return err
	}
	if err := s.validator(map[string]any(model.Document)); err != nil {
		return err
	}
	if err := s.postValidate(model); err != nil {
		return err
	}
	if err := s.preCreate(model); err != nil {
		return err
	}
	if err := s.preSave(model); err != nil {
		return err
	}
	if err := operator(ctx); err != nil {
		return err
	}
	if err := s.postCreate(model); err != nil {
		return err
	}
	if err := s.postSave(model); err != nil {
		return err
	}
	return nil
}

// todo: move these to collection.go?
//func (s *Schema) createMany(attrs []Attributes) ([]Model, error) {
//	for _, v := range attrs {
//		if err := s.validator(v); err != nil {
//			return nil, err
//		}
//	}
//
//	var docs []Model
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
//func FindOne(query Query) (*Model, error) {
//	filter := bson.M{}
//	var doc *Model
//	if err := odm.db.Collection(s.name).FindOne(ctx(), filter).Decode(&doc); err != nil {
//		return nil, err
//	}
//	return doc, nil
//}
//
//func FindMany(query Query) ([]Model, error) {
//	filter := bson.M{}
//	cur, err := odm.db.Collection(s.name).Find(ctx(), filter)
//	if err != nil {
//		return nil, err
//	}
//	var docs []Model
//	if err := cur.All(ctx(), &docs); err != nil {
//		return nil, err
//	}
//	return docs, nil
//}
//
//func UpdateOne(query Query, attr Attributes) (*Model, error) {
//	filter := bson.M{}
//	update := bson.M{
//		"$unset": bson.M{"unset_me": 1},
//	}
//	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
//
//	var doc *Model
//	err := odm.db.Collection(s.name).FindOneAndUpdate(ctx(), filter, update, opts).Decode(&doc)
//	if err != nil {
//		return nil, err
//	}
//	return doc, nil
//}
//
//func UpdateMany(query Query, attr Attributes) ([]Model, error) {
//	var docs []Model
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
//func RemoveOne(query Query) (*Model, error) {
//	filter := bson.M{}
//	var doc *Model
//	if err := odm.db.Collection(s.name).FindOneAndDelete(ctx(), filter).Decode(&doc); err != nil {
//		return nil, err
//	}
//	return doc, nil
//}
//
//func RemoveMany(query Query) ([]Model, error) {
//	var docs []Model
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
