package mongolia

// Coll returns the collection associated with a model.
func Coll(m Model) *Collection {
	if v, ok := m.(CollectionGetter); ok {
		return v.Collection()
	}
	return nil // todo: handle
}
