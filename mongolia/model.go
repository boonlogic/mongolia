package mongolia

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
func (m *DefaultModel) Creating() error {
	return m.DateFields.Creating()
}

// Saving function calls the Saving hooks of DefaultModel's inner fields.
func (m *DefaultModel) Saving() error {
	return m.DateFields.Saving()
}
