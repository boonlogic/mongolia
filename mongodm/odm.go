package mongodm

type ODM struct {
	colls map[string]*Collection
}

func (odm *ODM) Collection(*Model) *Collection {
	return nil
}
