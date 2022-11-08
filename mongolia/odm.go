package mongolia

type ODM struct {
	colls map[string]*Collection
}

func NewODM() *ODM {
	return &ODM{
		colls: make(map[string]*Collection),
	}
}
