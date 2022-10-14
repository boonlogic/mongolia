package odm

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema struct {
	Name string
	Definition []byte
	Hooks *Hooks
}

func (s *Schema) PreValidate(any) *Model {
	return nil
}

func (s *Schema) PostValidate(any) *Model {
	return nil
}

func (s *Schema) PreSave(any) *Model {
	return nil
}

func (s *Schema) PostSave(any) *Model {
	return nil
}

func (s *Schema) PreCreate(any) *Model {
	return nil
}

func (s *Schema) PostCreate(any) *Model {
	return nil
}

func (s *Schema) PreUpdate(any) *Model {
	return nil
}

func (s *Schema) PostUpdate(any) *Model {
	return nil
}

func (s *Schema) PreRemove(any) *Model {
	return nil
}

func (s *Schema) PostRemove(any) *Model {
	return nil
}

func (s *Schema) CreateOne(obj any) (doc *Document, err error) {
	// todo: JSON schema validation here

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

func (s *Schema) CreateMany(any) ([]Document, error) {
	return nil, nil
}

func (s *Schema) FindOne(any) (*Document, error) {
	return nil, nil
}

func (s *Schema) FindMany(any) ([]Document, error) {
	return nil, nil
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
