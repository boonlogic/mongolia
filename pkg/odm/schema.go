package odm

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// REST      restapi    pkg/odm (validate/hooks) mongo-driver (access layer)
// -----------------------------------------------------------------------------------
// POST   -> create      -> createOne, createMany -> mongo.InsertOne, mongo.InsertMany
// GET    -> read        -> findOne, findMany     -> mongo.FindOne, mongo.FindMany
// PUT    -> update      -> updateOne, updateMany -> mongo.UpdateOne, mongo.UpdateMany
// DELETE -> delete      -> removeOne, removeMany -> mongo.DeleteOne, mongo.DeleteMany

func (s Schema) CreateOne(obj any) (*Document, error) {
	if err := s.Validator(obj); err != nil {
		return nil, err
	}

	insres, err := db.Collection(s.Name).InsertOne(ctx(), obj)
	if err != nil {
		return nil, err
	}
	id := insres.InsertedID.(primitive.ObjectID)

	var doc *Document
	query := bson.M{"_id": id}
	if err := db.Collection(s.Name).FindOne(ctx(), query).Decode(&doc); err != nil {
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

	res, err := db.Collection(s.Name).InsertMany(ctx(), arr)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(res.InsertedIDs))
	for i, v := range res.InsertedIDs {
		ids[i] = v.(primitive.ObjectID)
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
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

func (s Schema) FindOne(any) (*Document, error) {
	filter := bson.D{}
	var doc *Document
	if err := db.Collection(s.Name).FindOne(ctx(), filter).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) FindMany(any) ([]Document, error) {
	filter := bson.D{}
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
	filter := bson.D{}
	update := bson.D{}
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
	var doc *Document
	err := db.Collection(s.Name).FindOneAndUpdate(ctx(), filter, update, opts).Decode(&doc); if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s Schema) UpdateMany(any) ([]Document, error) {
	return nil, nil
}

func (s Schema) RemoveOne(any) (*Document, error) {
	return nil, nil
}

func (s Schema) RemoveMany(any) ([]Document, error) {
	return nil, nil
}

