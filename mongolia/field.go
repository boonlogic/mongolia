package mongolia

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// IDField contains a model's ID field.
type IDField struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

// DateFields contains the `created_at` and `updated_at`
// fields that autofill when inserting or updating a model.
type DateFields struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// PrepareID prepares the id value to be used for mongo lookups.
func (f *IDField) PrepareID(id any) (any, error) {
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return v, nil
	default:
		return nil, errors.New(fmt.Sprintf("id must be a string or objectId (got type %T)", v))
	}
}

func (f *IDField) GetID() any {
	return f.ID
}

func (f *IDField) SetID(id any) {
	f.ID = id.(primitive.ObjectID)
}

func (f *DateFields) Creating() error {
	f.CreatedAt = time.Now().UTC()
	return nil
}

func (f *DateFields) Saving() error {
	f.UpdatedAt = time.Now().UTC()
	return nil
}