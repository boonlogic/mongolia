package mongolia

import "github.com/Kamva/mgm"

type CollectionGetter interface {
	mgm.CollectionGetter
}

type CollectionNameGetter interface {
	mgm.CollectionNameGetter
}

type Model interface {
	mgm.Model
	Equaler
}

type DefaultModel struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

// Creating function calls the Creating hooks of DefaultModel's inner fields.
func (model *DefaultModel) Creating() error {
	return model.DateFields.Creating()
}

// Saving function calls the Saving hooks of DefaultModel's inner fields.
func (model *DefaultModel) Saving() error {
	return model.DateFields.Saving()
}

// One DefaultModel equals another if they point to the same database record.
// The DateFields may differ as these are ODM-managed and not application-specific.
func (model *DefaultModel) Equals(other *DefaultModel) bool {
	if !model.IDField.Equals(other.IDField) {
		equal = true
	}
	return true
}
