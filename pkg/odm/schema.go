package odm

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema struct {
	Name       string
	Definition []byte
	Validator  func(any) error
	Hooks      *Hooks
}

func (s *Schema) preValidate(any) *Model {
	return nil
}

func (s *Schema) postValidate(any) *Model {
	return nil
}

func (s *Schema) preSave(any) *Model {
	return nil
}

func (s *Schema) postSave(any) *Model {
	return nil
}

func (s *Schema) preCreate(any) *Model {
	return nil
}

func (s *Schema) postCreate(any) *Model {
	return nil
}

func (s *Schema) preUpdate(any) *Model {
	return nil
}

func (s *Schema) postUpdate(any) *Model {
	return nil
}

func (s *Schema) preRemove(any) *Model {
	return nil
}

func (s *Schema) postRemove(any) *Model {
	return nil
}

func (s *Schema) CreateOne(obj any) (*Document, error) {
	if err := s.Validator(obj); err != nil {
		return nil, err
	}

	insres, err := db.Collection(s.Name).InsertOne(ctx(), obj)
	if err != nil {
		return nil, err
	}
	id := insres.InsertedID.(primitive.ObjectID)

	var out bson.D
	if err := db.Collection(s.Name).FindOne(ctx(), bson.M{"_id": id}).Decode(&out); err != nil {
		return nil, err
	}
	d := Document(out)
	return &d, nil
}

func (s *Schema) CreateMany(objSlice any) ([]Document, error) {
	objs, ok := objSlice.([]any)
	if !ok {
		return nil, errors.New("CreateMany must take a slice as argument")
	}

	for _, obj := range objs {
		if err := s.Validator(obj); err != nil {
			return nil, err
		}
	}

	res, err := db.Collection(s.Name).InsertMany(ctx(), objs)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(res.InsertedIDs))
	for i, v := range res.InsertedIDs {
		ids[i] = v.(primitive.ObjectID)
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
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

func (s *Schema) FindOne(any) (*Document, error) {
	var doc *Document
	if err := db.Collection(s.Name).FindOne(ctx(), bson.D{}).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Schema) FindMany(any) ([]Document, error) {
	cur, err := db.Collection(s.Name).Find(ctx(), bson.D{})
	if err != nil {
		return nil, err
	}
	var docs []Document
	if err := cur.All(ctx(), &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s *Schema) UpdateOne(any) (*Document, error) {
	return nil, nil
}

func (s *Schema) UpdateMany(any) ([]Document, error) {
	return nil, nil
}

func (s *Schema) RemoveOne(any) (*Document, error) {
	return nil, nil
}

func (s *Schema) RemoveMany(any) ([]Document, error) {
	return nil, nil
}

