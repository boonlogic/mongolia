package mongolia

import "context"

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

	GetID() interface{}
	SetID(id interface{})
}

type DefaultModel struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

// Creating function calls the Creating hooks of DefaultModel's inner fields.
func (m *DefaultModel) PreCreate(ctx context.Context) error {
	return m.DateFields.PreCreate(ctx)
}

// Saving function calls the Saving hooks of DefaultModel's inner fields.
func (m *DefaultModel) PreSave(ctx context.Context) error {
	return m.DateFields.PreSave(ctx)
}
