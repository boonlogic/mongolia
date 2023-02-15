package mongolia

import (
	validation "github.com/go-ozzo/ozzo-validation"
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
	ValidateUpdate(update any) error
	ValidateUpdateModel() error

	PreCreate() error
	PostCreate() error
	PreUpdate(update any) error
	PreUpdateModel() error
	PostUpdateModel(result *mongo.UpdateResult) error
	PreSave() error
	PostSave() error
	PostRead() error
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
	return validation.ValidateStruct(m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
	)
}

func (m *DefaultModel) ValidateCreate() error {
	return m.Validate()
}

//validation when doing a partial update
func (m *DefaultModel) ValidateUpdate(update any) error {
	return nil
}

//validation when updating entire model
func (m *DefaultModel) ValidateUpdateModel() error {
	return m.Validate()
}

func (m *DefaultModel) PreCreate() error {
	return m.DateFields.PreCreate()
}

func (m *DefaultModel) PostCreate() error {
	return nil
}

//Hook when doing a partial update
func (m *DefaultModel) PreUpdate(update any) error {
	switch update.(type) {
	case bson.D:
		BSONUpdateAtHook(update.(bson.D))
		return nil
	default:
		return nil
	}
}

// Hooks when pre updating entire model
func (m *DefaultModel) PreUpdateModel() error {
	return nil
}

// Hooks when post updating entire model
func (m *DefaultModel) PostUpdateModel(result *mongo.UpdateResult) error {
	return nil
}

func (m *DefaultModel) PreSave() error {
	return m.DateFields.PreSave()
}

func (m *DefaultModel) PostSave() error {
	return nil
}

func (m *DefaultModel) PostRead() error {
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
