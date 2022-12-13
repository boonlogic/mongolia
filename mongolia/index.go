package mongolia

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func insertIndexes(ctx context.Context, coll *mongo.Collection, indexes []mongo.IndexModel) error {
	_, err := coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	return nil
}

func listIndexes(ctx context.Context, coll *mongo.Collection) ([]bson.D, error) {
	curs, err := coll.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	var docs []bson.D
	if err = curs.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func dropIndex(coll *mongo.Collection, name string) error {
	if _, err := coll.Indexes().DropOne(context.Background(), name); err != nil {
		return err
	}
	return nil
}

func validateIndexes(have []bson.D, want []mongo.IndexModel) error {
	//Todo
	return nil
}

func PopulateIndexes(ctx context.Context, coll *mongo.Collection, indexes interface{}) error {
	fmt.Printf("PopulateIndexes... \n")
	var indexModel []mongo.IndexModel
	var err error
	switch v := indexes.(type) {
	case []mongo.IndexModel:
		indexModel = indexes.([]mongo.IndexModel)
	default:
		return errors.New(fmt.Sprintf("Unknown index type %v", v))
	}

	fmt.Printf("Requested Indexes: %v \n", indexModel)

	//Get current indexes
	current, err := listIndexes(ctx, coll)
	if err != nil {
		return err
	}
	fmt.Printf("Current Indexes: %v \n", current)

	//Compare with existing
	if err := validateIndexes(current, indexModel); err != nil {
		return err
	}

	//Insert indexes
	if err := insertIndexes(ctx, coll, indexModel); err != nil {
		return err
	}
	return nil
}
