package mgm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// DefaultModel implements the Model interface.
type DefaultModel struct {
}

func Creating(context.Context) error {
	return nil
}

func Created(context.Context) error {
	return nil
}

func Updating(context.Context) error {
	return nil
}

func Updated(ctx context.Context, result *mongo.UpdateResult) error {
	return nil
}

func Saving(context.Context) error {
	return nil
}

func Saved(context.Context) error {
	return nil
}

func Deleting(context.Context) error {
	return nil
}

func Deleted(ctx context.Context, result *mongo.DeleteResult) error {
	return nil
}
