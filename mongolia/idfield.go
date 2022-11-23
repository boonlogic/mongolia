package mongolia

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDField contains a model's ID field.
type IDField struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

// PrepareID prepares the id value to be used for mongo lookups.
func (f *IDField) PrepareID(id any) (any, error) {
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	default:
		return v, nil // assume it is an objectId
	}
}

func (f *IDField) IsNew() bool {
	return f.GetID() == primitive.ObjectID{}
}

func (f *IDField) GetID() any {
	return f.ID
}

func (f *IDField) SetID(id any) {
	f.ID = id.(primitive.ObjectID)
}

func (f *IDField) Equals(other *IDField) bool {
	if f == nil || other == nil {
		return false
	}
	return *f == *other
}
