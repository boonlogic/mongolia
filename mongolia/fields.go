package mongolia

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DateFields are the date-related fields that are auto-updated
// by the ODM whenever the Model gets created or updated.
type DateFields struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (f *DateFields) PreCreate() error {
	f.CreatedAt = time.Now().UTC()
	return nil
}

func (f *DateFields) PreSave() error {
	f.UpdatedAt = time.Now().UTC()
	return nil
}

// IDField contains a model's ID field.
type IDField struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

// PrepareID prepares the id value to be used for mongo lookups.
func (f *IDField) PrepareID(id any) (any, error) {
	switch v := id.(type) {
	case string:
		result, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return result, errors.New("Invalid ID")
		}
		return result, nil
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

// LA note:
// I the recommend adding the hook-managed fields "CreatedBy" and "UpdatedBy" once
// ACL for v2 is up and running. There will probably be utility in knowing which entity
// in Amber ecosystem modified the record last. That implementation should go here.

type EntityFields struct {
	CreatedBy string `json:"createdBy" bson:"createdBy"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
}
