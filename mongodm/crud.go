package mongodm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// createOne inserts a document and sets its id field to the inserted id.
func createOne(ctx context.Context, coll string, doc *Document) error {
	id, err := insertOne(ctx, coll, bson.M(*doc.attrs))
	if err != nil {
		return err
	}
	doc.id = id
	return nil
}
