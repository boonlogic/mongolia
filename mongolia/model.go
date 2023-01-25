package mongolia

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FindResult struct {
	Skip       int64
	Limit      int64
	Filtered   int64
	Collection int64
}

// Model interface contains base methods that must be implemented by
// each model. If you're using the `DefaultModel` struct in your model,
// you don't need to implement any of these methods.
type Model interface {
	// PrepareID converts the id value if needed, then
	// returns it (e.g convert string to objectId).
	PrepareID(id interface{}) (interface{}, error)

	IsNew() bool
	GetID() interface{}
	SetID(id interface{})

	GetTagReferences() map[string]string

	Validate() error
	ValidateRead() error
	ValidateCreate() error
	ValidateUpdate() error

	PreCreate() error
	PostCreate() error
	PrePartialUpdate(update any) error
	PreUpdate() error
	PostUpdate(result *mongo.UpdateResult) error
	PreSave() error
	PostSave() error
	PreDelete() error
	PostDelete(result *mongo.DeleteResult) error
}

type DefaultModel struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

func (m *DefaultModel) Validate() error {
	return nil
}

func (m *DefaultModel) ValidateRead() error {
	return m.Validate()
}

func (m *DefaultModel) ValidateCreate() error {
	return m.Validate()
}

func (m *DefaultModel) ValidateUpdate() error {
	return m.Validate()
}

func (m *DefaultModel) PreCreate() error {
	return m.DateFields.PreCreate()
}

func (m *DefaultModel) PostCreate() error {
	return nil
}

func (m *DefaultModel) PrePartialUpdate(update any) error {
	switch v := update.(type) {
	case bson.D:
		BSONUpdateAtHook(update.(bson.D))
		return nil
	default:
		return errors.New(fmt.Sprintf("Unknown Partial Update Type %v \n", v))
	}
}

func (m *DefaultModel) PreUpdate() error {
	return nil
}

func (m *DefaultModel) PostUpdate(result *mongo.UpdateResult) error {
	return nil
}

func (m *DefaultModel) PreSave() error {
	return m.DateFields.PreSave()
}

func (m *DefaultModel) PostSave() error {
	return nil
}

func (m *DefaultModel) PreDelete() error {
	return nil
}

func (m *DefaultModel) PostDelete(result *mongo.DeleteResult) error {
	return nil
}

func (m *DefaultModel) GetTagReferences() map[string]string {
	return GetStructTags(*m, "ref")
}
