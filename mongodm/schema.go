package mongodm

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Schema struct {
	name     string          // is globally unique
	validate func(any) error // anything accepted by this function "matches the schema"

	hooks *Hooks
}

func convertToDocument(obj any) (Document, error) {
	buf, err := bson.MarshalExtJSON(obj, true, false)
	if err != nil {
		return nil, err
	}
	var doc Document
	if err := bson.UnmarshalExtJSON(buf, true, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}
