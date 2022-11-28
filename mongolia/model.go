package mongolia

import "github.com/Kamva/mgm"

//
//type CollectionGetter interface {
//	mgm.CollectionGetter
//}
//
//type CollectionNameGetter interface {
//	mgm.CollectionNameGetter
//}

type Model interface {
	mgm.Model
	Equaler
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

// Equals function tells whether this Model equals another.
func (m *DefaultModel) Equals(other Model) bool {
	return m.GetID()
}
