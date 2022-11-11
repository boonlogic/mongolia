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
