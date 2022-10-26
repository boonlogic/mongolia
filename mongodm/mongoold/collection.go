package mongoold

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Collection wraps a mongoold.Collection.
type Collection struct {
	coll *mongo.Collection
}

func (c *Collection) insertOne(ctx context.Context, doc any) (primitive.ObjectID, error) {
	result, err := c.coll.InsertOne(ctx, doc)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
