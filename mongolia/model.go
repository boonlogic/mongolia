package mongolia

type CollectionGetter interface {
	Collection() *Collection
}

type CollectionNameGetter interface {
	CollectionName() string
}

type Model interface {
	// PrepareID converts the id value into a mongo objectId.
	PrepareID(id any) (any, error)

	GetID() any
	SetID(id any)
}
