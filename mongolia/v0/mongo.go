package v0

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func listIndexes(ctx context.Context, coll *mongo.Collection) (bson.M, error) {
	indexview := coll.Indexes()
	var indexes bson.M
	curs, err := indexview.List(ctx)
	if err != nil {
		return nil, err
	}
	if err = curs.All(ctx, &indexes); err != nil {
		return nil, err
	}
	return indexes, nil
}

func idToObjectID(id any) (primitive.ObjectID, error) {
	switch v := id.(type) {
	case primitive.ObjectID:
		return v, nil
	case string:
		oid, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return primitive.ObjectID{}, err
		}
		return oid, nil
	default:
		return primitive.ObjectID{}, errors.New(fmt.Sprintf("expected id of type string (got type %T)", id))
	}
}

func insertOne(ctx context.Context, coll *mongo.Collection, doc map[string]any) (string, error) {
	if id, ok := doc["id"]; ok {
		oid, err := idToObjectID(id)
		if err != nil {
			return "", err
		}
		delete(doc, "id")
		if !oid.IsZero() {
			doc["_id"] = oid
		}
	}
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}
