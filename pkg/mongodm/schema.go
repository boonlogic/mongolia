package mongodm

import (
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Schema struct {
	name     string          // is globally unique
	validate func(any) error // anything accepted by this function "matches the schema"

	hooks *types.Hooks
}

func convertToDocument(obj any) (types.Document, error) {
	buf, err := bson.MarshalExtJSON(obj, true, false)
	if err != nil {
		return nil, err
	}
	var doc types.Document
	if err := bson.UnmarshalExtJSON(buf, true, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}
