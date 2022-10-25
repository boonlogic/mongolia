package mongodm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func insertOne(ctx context.Context, coll string, doc bson.M) (primitive.ObjectID, error) {
	result, err := odm.db.Collection(coll).InsertOne(ctx, doc)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
