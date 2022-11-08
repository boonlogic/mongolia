package mongolia

type CollectionGetter interface {
	Collection() *Collection
}

type CollectionNameGetter interface {
	CollectionName() string
}

type Model interface {
	PrepareID(id any) (any, error)
	IsNew() bool
	GetID() any
	SetID(id any)
}

type DefaultModel struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

func (model *DefaultModel) PreCreate() error {
	return model.DateFields.PreCreate()
}

func (model *DefaultModel) PreSave() error {
	return model.DateFields.PreSave()
}
