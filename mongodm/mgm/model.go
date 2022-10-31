package mgm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model interface {
	Creating(context.Context) error
	Created(context.Context) error
	Updating(context.Context) error
	Updated(ctx context.Context, result *mongo.UpdateResult) error
	Saving(context.Context) error
	Saved(context.Context) error
	Deleting(context.Context) error
	Deleted(ctx context.Context, result *mongo.DeleteResult) error
}
